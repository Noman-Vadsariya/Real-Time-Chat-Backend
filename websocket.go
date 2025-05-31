package main

import (
	"encoding/json"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// Client represents a connected WebSocket client
type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan []byte
	userID   uint
	username string
}

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

// Message represents the structure of our WebSocket messages
type WSMessage struct {
	Type      string      `json:"type"`
	SenderID  uint        `json:"sender_id"`
	Content   string      `json:"content"`
	Timestamp string      `json:"timestamp"`
	Data      interface{} `json:"data,omitempty"`
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			h.mu.Unlock()
			// Broadcast user online status
			h.broadcastUserStatus(client.userID, true)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
			h.mu.Unlock()
			// Broadcast user offline status
			h.broadcastUserStatus(client.userID, false)

		case message := <-h.broadcast:
			h.mu.RLock()
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
			h.mu.RUnlock()
		}
	}
}

func (h *Hub) broadcastUserStatus(userID uint, online bool) {
	statusMsg := WSMessage{
		Type: "user_status",
		Data: map[string]interface{}{"user_id": userID, "online": online},
	}
	msg, err := json.Marshal(statusMsg)
	if err != nil {
		log.Printf("Error marshaling status message: %v", err)
		return
	}
	h.broadcast <- msg
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var wsMsg WSMessage
		if err := json.Unmarshal(message, &wsMsg); err != nil {
			log.Printf("Error unmarshaling message: %v", err)
			continue
		}

		// Handle different message types
		switch wsMsg.Type {
		case "message":
			// Store message in database
			msg := Message{
				SenderID:   c.userID,
				ReceiverID: uint(wsMsg.Data.(map[string]interface{})["receiver_id"].(float64)),
				Content:    wsMsg.Content,
			}
			if err := db.Create(&msg).Error; err != nil {
				log.Printf("Error storing message: %v", err)
				continue
			}

			// Broadcast to specific recipient
			c.hub.mu.RLock()
			for client := range c.hub.clients {
				if client.userID == msg.ReceiverID {
					client.send <- message
					break
				}
			}
			c.hub.mu.RUnlock()
		}
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for {
		message, ok := <-c.send
		if !ok {
			c.conn.WriteMessage(websocket.CloseMessage, []byte{})
			return
		}

		w, err := c.conn.NextWriter(websocket.TextMessage)
		if err != nil {
			return
		}
		w.Write(message)

		if err := w.Close(); err != nil {
			return
		}
	}
}
