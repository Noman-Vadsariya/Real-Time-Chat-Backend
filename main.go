package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/noman.nooruddin/chat-backend/docs"
	gsfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

// @title Go Chat Backend API
// @version 1.0
// @description Real-time chat backend with JWT, WebSocket, and PostgreSQL
// @host localhost:8080
// @BasePath /

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Initialize database
	var err error
	if os.Getenv("DB_TYPE") == "sqlite" {
		db, err = gorm.Open(sqlite.Open("chat.db"), &gorm.Config{})
	} else {
		dsn := "host=" + os.Getenv("DB_HOST") +
			" user=" + os.Getenv("DB_USER") +
			" password=" + os.Getenv("DB_PASSWORD") +
			" dbname=" + os.Getenv("DB_NAME") +
			" port=" + os.Getenv("DB_PORT") +
			" sslmode=disable"
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	}

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Auto-migrate the schema
	db.AutoMigrate(&User{}, &Message{})

	// Start WebSocket hub
	go startWebSocketHub()

	// Initialize Gin router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Serve static files
	r.Static("/client", "./client")

	// Public routes
	r.POST("/register", register)
	r.POST("/login", login)

	// Protected routes
	auth := r.Group("/")
	auth.Use(authMiddleware())
	{
		auth.GET("/ws", wsHandler)
		auth.GET("/messages/:userId", getMessages)
		auth.GET("/users", getUsers)
	}

	// Serve Swagger docs
	r.GET("/swagger/*any", ginSwagger.WrapHandler(gsfiles.Handler))

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
