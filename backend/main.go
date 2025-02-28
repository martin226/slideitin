package main

import (
	"context"
	"log"
	"os"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/slideitin/backend/controllers"
	"github.com/slideitin/backend/services/slides"
	"github.com/slideitin/backend/services/queue"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using system environment variables")
	}

	// Initialize the router
	router := gin.Default()

	// Get frontend URL from environment variable
	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		frontendURL = "http://localhost:3000" // Fallback for local development
		log.Println("Warning: FRONTEND_URL not set, using default:", frontendURL)
	}

	// Configure CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{frontendURL}, // Use environment variable
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Cache-Control", "Connection", "Access-Control-Allow-Origin"},
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "Cache-Control", "Content-Encoding", "Transfer-Encoding"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Get API keys and project configuration from environment variables
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		log.Println("Warning: GEMINI_API_KEY not set, using placeholder")
		apiKey = "sk-placeholder"
	}

	// Initialize Firestore client
	ctx := context.Background()
	projectID := os.Getenv("GOOGLE_CLOUD_PROJECT")
	if projectID == "" {
		log.Println("Warning: GOOGLE_CLOUD_PROJECT not set, using default")
		projectID = "slideitin"
	}

	firestoreClient, err := firestore.NewClient(ctx, projectID)

	if err != nil {
		log.Fatalf("Failed to initialize Firestore: %v", err)
	}
	defer firestoreClient.Close()

	// Initialize services
	slideService := slides.NewSlideService(apiKey)
	
	// Initialize queue service with Firestore and Gemini service
	queueService := queue.NewService(firestoreClient, slideService)

	// Initialize controllers
	slideController := controllers.NewSlideController(slideService, queueService)

	// API routes
	v1 := router.Group("/v1")
	{
		// Slide generation endpoint - adds job to queue and returns immediately
		v1.POST("/generate", slideController.GenerateSlides)
		
		// Streaming status endpoint - combines status checking and streaming
		v1.GET("/slides/:id", slideController.StreamSlideStatus)
        
		// Result retrieval endpoint - serves the generated presentation
		v1.GET("/results/:id", slideController.GetSlideResult)
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