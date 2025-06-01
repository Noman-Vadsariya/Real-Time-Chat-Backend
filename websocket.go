package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for development
	},
}

// Client represents a connected WebSocket client
type WSClient struct {
	conn     *websocket.Conn
	userID   uint
	username string
	send     chan []byte
}

// Message represents the structure of our WebSocket messages
type WSMessage struct {
	Type    string      `json:"type"`
	Content string      `json:"content,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

var (
	wsClients    = make(map[*WSClient]bool)
	wsBroadcast  = make(chan []byte)
	wsRegister   = make(chan *WSClient)
	wsUnregister = make(chan *WSClient)
	wsMutex      = &sync.Mutex{}
)

func wsHandler(c *gin.Context) {
	// Get token from query parameter
	token := c.Query("token")
	log.Printf("Received WebSocket connection request with token: %s", token)

	if token == "" {
		log.Printf("No token provided in WebSocket request")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Token is required"})
		return
	}

	// Validate token
	claims, err := validateToken(token)
	if err != nil {
		log.Printf("Token validation failed: %v", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
		return
	}

	log.Printf("Token validated successfully, upgrading connection for user %d (%s)", claims.UserID, claims.Username)

	// Set headers for WebSocket upgrade
	c.Header("Upgrade", "websocket")
	c.Header("Connection", "Upgrade")

	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to upgrade connection: %v", err)
		return
	}

	log.Printf("WebSocket connection upgraded successfully for user %d (%s)", claims.UserID, claims.Username)

	client := &WSClient{
		conn:     conn,
		userID:   claims.UserID,
		username: claims.Username,
		send:     make(chan []byte, 256),
	}

	wsRegister <- client
	log.Printf("WebSocket client registered for user %d (%s)", client.userID, client.username)

	go wsWritePump(client)
	go wsReadPump(client)
}

func wsWritePump(c *WSClient) {
	defer func() {
		log.Printf("Write pump closing for user %d", c.userID)
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				log.Printf("Send channel closed for user %d", c.userID)
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Printf("Error getting next writer for user %d: %v", c.userID, err)
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				log.Printf("Error closing writer for user %d: %v", c.userID, err)
				return
			}
		}
	}
}

func wsReadPump(c *WSClient) {
	defer func() {
		log.Printf("Read pump closing for user %d", c.userID)
		wsUnregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error for user %d: %v", c.userID, err)
			}
			break
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("Error unmarshaling message from user %d: %v", c.userID, err)
			continue
		}

		switch wsMsg.Type {
		case "message":
			log.Printf("Received message from user %d: %s", c.userID, wsMsg.Content)
			if data, ok := wsMsg.Data.(map[string]interface{}); ok {
				if receiverID, ok := data["receiver_id"].(float64); ok {
					// Create message in database
					message := Message{
						SenderID:   c.userID,
						ReceiverID: uint(receiverID),
						Content:    wsMsg.Content,
					}
					if err := db.Create(&message).Error; err != nil {
						log.Printf("Error saving message: %v", err)
						continue
					}
					log.Printf("Message saved to database: %+v", message)

					// Send to receiver if online
					wsMutex.Lock()
					for client := range wsClients {
						if client.userID == uint(receiverID) {
							msgBytes, _ := json.Marshal(WSMessage{
								Type:    "message",
								Content: wsMsg.Content,
								Data: map[string]interface{}{
									"sender_id":   c.userID,
									"sender_name": c.username,
								},
							})
							client.send <- msgBytes
							log.Printf("Message forwarded to receiver %d", receiverID)
						}
					}
					wsMutex.Unlock()
				}
			}
		}
	}
}

func startWebSocketHub() {
	log.Println("Starting WebSocket hub...")
	for {
		select {
		case client := <-wsRegister:
			log.Printf("Registering new client for user %d (%s)", client.userID, client.username)
			wsMutex.Lock()
			// Remove any existing connection for this user
			for existingClient := range wsClients {
				if existingClient.userID == client.userID {
					log.Printf("Removing existing connection for user %d", client.userID)
					existingClient.conn.Close()
					delete(wsClients, existingClient)
				}
			}
			// Add new client
			wsClients[client] = true
			log.Printf("User %d (%s) connected. Total clients: %d", client.userID, client.username, len(wsClients))

			// Broadcast user's online status
			statusMsg := WSMessage{
				Type: "user_status",
				Data: map[string]interface{}{
					"user_id":  client.userID,
					"username": client.username,
					"online":   true,
				},
			}
			statusBytes, _ := json.Marshal(statusMsg)
			log.Printf("Broadcasting online status for user %d", client.userID)
			wsBroadcast <- statusBytes
			wsMutex.Unlock()

		case client := <-wsUnregister:
			log.Printf("Unregistering client for user %d (%s)", client.userID, client.username)
			wsMutex.Lock()
			if _, ok := wsClients[client]; ok {
				delete(wsClients, client)
				close(client.send)
				log.Printf("User %d (%s) disconnected. Total clients: %d", client.userID, client.username, len(wsClients))

				// Notify all clients about the user going offline
				if client.userID != 0 {
					statusMsg := WSMessage{
						Type: "user_status",
						Data: map[string]interface{}{
							"user_id":  client.userID,
							"username": client.username,
							"online":   false,
						},
					}
					statusBytes, _ := json.Marshal(statusMsg)
					log.Printf("Broadcasting offline status for user %d", client.userID)
					wsBroadcast <- statusBytes
				}
			}
			wsMutex.Unlock()

		case message := <-wsBroadcast:
			wsMutex.Lock()
			log.Printf("Broadcasting message to %d clients", len(wsClients))
			for client := range wsClients {
				select {
				case client.send <- message:
					log.Printf("Message sent to user %d", client.userID)
				default:
					log.Printf("Failed to send message to user %d, closing connection", client.userID)
					close(client.send)
					delete(wsClients, client)
				}
			}
			wsMutex.Unlock()
		}
	}
}
