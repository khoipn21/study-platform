package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"chatbot-service/internal/config"
	"chatbot-service/internal/handler"
	"chatbot-service/internal/repository"
	"chatbot-service/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	db, err := connectDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repositories
	chatRepo := repository.NewSimpleChatRepository(db)
	simpleAnalyticsRepo := repository.NewSimpleAnalyticsRepository(db)

	// Initialize services
	aiService := service.NewAIService(&cfg.OpenAI)
	analyticsService := service.NewSimpleAnalyticsService(simpleAnalyticsRepo, chatRepo)

	// Initialize handlers
	chatHandler := handler.NewSimpleChatHandler(chatRepo, aiService)
	wsHandler := handler.NewSimpleWebSocketHandler(chatRepo, aiService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	// Setup router
	router := setupRouter(chatHandler, wsHandler, analyticsHandler)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Chatbot service starting on %s", addr)
	
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func connectDB(cfg config.DatabaseConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func setupRouter(chatHandler *handler.SimpleChatHandler, wsHandler *handler.SimpleWebSocketHandler, analyticsHandler *handler.AnalyticsHandler) *gin.Engine {
	router := gin.Default()

	// CORS middleware
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "chatbot"})
	})

	// Authentication middleware (simplified - in production, this should validate JWT tokens)
	authMiddleware := func(c *gin.Context) {
		// For now, we'll expect user_id in headers
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID required"})
			c.Abort()
			return
		}

		// Parse UUID
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID format"})
			c.Abort()
			return
		}

		// In production, validate JWT token and extract user info
		c.Set("user_id", userID)
		c.Next()
	}

	// API routes
	api := router.Group("/api/v1")
	api.Use(authMiddleware)
	{
		// Chat sessions
		api.POST("/sessions", chatHandler.CreateSession)
		api.GET("/sessions", chatHandler.GetUserSessions)
		api.GET("/sessions/:sessionId", chatHandler.GetSession)
		api.PUT("/sessions/:sessionId", chatHandler.UpdateSession)
		api.DELETE("/sessions/:sessionId", chatHandler.DeleteSession)

		// Messages
		api.POST("/chat", chatHandler.SendMessage)
		api.GET("/sessions/:sessionId/messages", chatHandler.GetMessages)

		// WebSocket endpoint
		api.GET("/ws", wsHandler.HandleWebSocket)

		// Analytics endpoints
		analytics := api.Group("/analytics")
		{
			analytics.GET("/overall", analyticsHandler.GetOverallAnalytics)
			analytics.GET("/me", analyticsHandler.GetMyAnalytics)
			analytics.GET("/user/:userID", analyticsHandler.GetUserAnalytics)
			analytics.GET("/course/:courseID", analyticsHandler.GetCourseAnalytics)
			analytics.GET("/time-based", analyticsHandler.GetTimeBasedAnalytics)
			analytics.GET("/real-time", analyticsHandler.GetRealTimeMetrics)
			analytics.GET("/sessions", analyticsHandler.GetSessionMetrics)
			analytics.GET("/usage", analyticsHandler.GetUsageStats)
			analytics.GET("/quality", analyticsHandler.GetResponseQuality)
			analytics.GET("/dashboard", analyticsHandler.GetAnalyticsDashboard)
			analytics.POST("/report", analyticsHandler.GenerateReport)
		}
	}

	return router
}