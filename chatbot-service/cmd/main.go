package main

import (
	"context"
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
	"github.com/redis/go-redis/v9"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// Connect to Redis
	redisClient, err := connectRedis(cfg.Redis)
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Close()

	// Initialize repositories
	chatRepo := repository.NewSimpleChatRepository(db)
	redisRepo := repository.NewRedisRepository(redisClient)
	simpleAnalyticsRepo := repository.NewSimpleAnalyticsRepository(db)

	// Initialize services
	aiService, err := service.NewAIService(&cfg.Gemini)
	if err != nil {
		log.Fatalf("Failed to initialize AI service: %v", err)
	}
	rateLimiter := service.NewRateLimiter(redisClient)
	analyticsService := service.NewSimpleAnalyticsService(simpleAnalyticsRepo, chatRepo)

	// Initialize handlers
	chatHandler := handler.NewSimpleChatHandler(chatRepo, redisRepo, aiService)
	wsHandler := handler.NewSimpleWebSocketHandler(chatRepo, aiService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)

	// Setup router
	router := setupRouter(chatHandler, wsHandler, analyticsHandler, rateLimiter, redisRepo)

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

func connectRedis(cfg config.RedisConfig) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return client, nil
}

func setupRouter(chatHandler *handler.SimpleChatHandler, wsHandler *handler.SimpleWebSocketHandler, analyticsHandler *handler.AnalyticsHandler, rateLimiter *service.RateLimiter, redisRepo *repository.RedisRepository) *gin.Engine {
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

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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

	// Rate limiter middleware
	rateLimiterMiddleware := func(c *gin.Context) {
		userIDInterface, exists := c.Get("user_id")
		if !exists {
			c.Next()
			return
		}

		userID, ok := userIDInterface.(uuid.UUID)
		if !ok {
			c.Next()
			return
		}

		allowed, remaining, err := rateLimiter.CheckLimit(c.Request.Context(), userID)
		if err != nil {
			log.Printf("Rate limiter error: %v", err)
			c.Next()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", "10")
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))

		if !allowed {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Rate limit exceeded. You can send up to 10 prompts per day.",
				"remaining": 0,
			})
			c.Abort()
			return
		}

		c.Next()
	}

	// API routes
	api := router.Group("/api/v1")
	api.Use(authMiddleware)
	{
		// Rate limit info endpoint (no rate limiting on this)
		api.GET("/rate-limit", func(c *gin.Context) {
			userIDInterface, exists := c.Get("user_id")
			if !exists {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
				return
			}

			userID, ok := userIDInterface.(uuid.UUID)
			if !ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
				return
			}

			usage, _ := rateLimiter.GetUsage(c.Request.Context(), userID)
			remaining, _ := rateLimiter.GetRemaining(c.Request.Context(), userID)

			c.JSON(http.StatusOK, gin.H{
				"limit":     10,
				"usage":     usage,
				"remaining": remaining,
			})
		})
		// Chat endpoints - with rate limiting
		chatGroup := api.Group("")
		chatGroup.Use(rateLimiterMiddleware)
		{
			// Messages
			chatGroup.POST("/chat", func(c *gin.Context) {
				// Call the handler first
				chatHandler.SendMessage(c)
				
				// Increment usage only if the request was successful
				if c.Writer.Status() < 400 {
					userIDInterface, exists := c.Get("user_id")
					if exists {
						if userID, ok := userIDInterface.(uuid.UUID); ok {
							rateLimiter.IncrementUsage(c.Request.Context(), userID)
						}
					}
				}
			})

			// WebSocket endpoint
			chatGroup.GET("/ws", wsHandler.HandleWebSocket)
		}

		// Chat sessions - no rate limiting
		api.POST("/sessions", chatHandler.CreateSession)
		api.GET("/sessions", chatHandler.GetUserSessions)
		api.GET("/sessions/:sessionId", chatHandler.GetSession)
		api.PUT("/sessions/:sessionId", chatHandler.UpdateSession)
		api.DELETE("/sessions/:sessionId", chatHandler.DeleteSession)
		api.GET("/sessions/:sessionId/messages", chatHandler.GetMessages)

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