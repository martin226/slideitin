package main

import (
	"log"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/slideitin/backend/controllers"
	"github.com/slideitin/backend/services/gemini"
	"github.com/slideitin/backend/services/queue"
)

func main() {
	// Initialize the router
	router := gin.Default()

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Cache-Control", "Connection", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Cache-Control", "Content-Encoding", "Transfer-Encoding"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Get Gemini API key from environment variable or use a default for development
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("Warning: GEMINI_API_KEY not set, using placeholder")
		apiKey = "sk-placeholder"
	}

	// Initialize services
	geminiService := gemini.NewService(apiKey)
	
	// Initialize queue service with Gemini service
	queueService := queue.NewService(geminiService)

	// Initialize controllers
	slideController := controllers.NewSlideController(geminiService, queueService)

	// API routes
	v1 := router.Group("/v1")
	{
		// Slide generation endpoint - adds job to queue and returns immediately
		v1.POST("/generate", slideController.GenerateSlides)
		
		// Streaming status endpoint - combines status checking and streaming
		v1.GET("/slides/:id", slideController.StreamSlideStatus)
	}

	// Start the server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("Starting server on port %s\n", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
} 