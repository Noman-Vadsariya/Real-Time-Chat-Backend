package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins in development
	},
}

func handleWebSocket(c *gin.Context) {
	userID := c.GetUint("user_id")
	username := c.GetString("username")

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256),
		userID:   userID,
		username: username,
	}

	client.hub.register <- client

	// Start goroutines for reading and writing
	go client.writePump()
	go client.readPump()
}

func getMessages(c *gin.Context) {
	userID := c.GetUint("user_id")
	otherUserID, err := strconv.ParseUint(c.Param("userId"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var messages []Message
	if err := db.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		userID, otherUserID, otherUserID, userID).
		Preload("Sender").
		Preload("Receiver").
		Order("created_at desc").
		Limit(50).
		Find(&messages).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching messages"})
		return
	}

	c.JSON(http.StatusOK, messages)
}

func getUsers(c *gin.Context) {
	var users []User
	if err := db.Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching users"})
		return
	}

	// Create a map of online users
	onlineUsers := make(map[uint]bool)
	hub.mu.RLock()
	for client := range hub.clients {
		onlineUsers[client.userID] = true
	}
	hub.mu.RUnlock()

	// Add online status to user data
	var userData []gin.H
	for _, user := range users {
		userData = append(userData, gin.H{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"online":   onlineUsers[user.ID],
		})
	}

	c.JSON(http.StatusOK, userData)
}
