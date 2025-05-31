package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/noman.nooruddin/chat-backend/docs"
	gsfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB
var hub *Hub

// @title Go Chat Backend API
// @version 1.0
// @description Real-time chat backend with JWT, WebSocket, and PostgreSQL
// @host localhost:8080
// @BasePath /

func initDB() {
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=disable"

	var err error
	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Auto migrate the schema
	db.AutoMigrate(&User{}, &Message{})
}

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	initDB()

	// Initialize WebSocket hub
	hub = newHub()
	go hub.run()

	// Initialize router
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})

	// Auth routes
	r.POST("/register", register)
	r.POST("/login", login)

	// Protected routes
	auth := r.Group("/")
	auth.Use(authMiddleware())
	{
		auth.GET("/ws", handleWebSocket)
		auth.GET("/messages/:userId", getMessages)
		auth.GET("/users", getUsers)
	}

	// Serve Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(gsfiles.Handler))
	// Serve static web client
	r.Static("/client", "./client")

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
