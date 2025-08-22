package main

import (
	"database/sql"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"

	"video-service/internal/config"
	"video-service/internal/handler"
	"video-service/internal/middleware"
	"video-service/internal/queue"
	"video-service/internal/repository"
	"video-service/internal/service"
	"video-service/internal/websocket"
)

func main() {
	// Load configuration
	cfg := config.Load()

	// Initialize database connection
	db, err := initDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize Redis client
	redisClient, err := queue.NewRedisClient(cfg)
	if err != nil {
		log.Fatalf("Failed to initialize Redis client: %v", err)
	}
	defer redisClient.Close()

	// Initialize repositories
	videoRepo := repository.NewVideoRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// Initialize services
	cloudflareService := service.NewCloudflareService(cfg)
	networkService := service.NewNetworkIntelligenceService(cfg, redisClient)
	videoService := service.NewVideoService(cfg, videoRepo, sessionRepo, cloudflareService, redisClient, networkService)

	// Initialize WebSocket hub
	hub := websocket.NewHub(redisClient, networkService)
	go hub.Run()

	// Initialize handlers
	healthHandler := handler.NewHealthHandler()
	videoHandler := handler.NewVideoHandler(videoService)
	wsHandler := handler.NewWebSocketHandler(hub)

	// Initialize Gin router
	router := gin.New()
	router.Use(gin.Recovery())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.CORSMiddleware())

	// Health endpoint
	router.GET("/health", healthHandler.Health)

	// API routes
	api := router.Group("/api")
	{
		videos := api.Group("/videos")
		{
			// Public routes (no auth required)
			videos.GET("/search", videoHandler.SearchVideos)
			videos.GET("/:video_id", middleware.OptionalAuthMiddleware(cfg.JWTSecret), videoHandler.GetVideo)
			videos.GET("/course/:course_id", videoHandler.ListCourseVideos)

			// Protected routes (auth required)
			protected := videos.Group("")
			protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
			{
				protected.POST("/upload", videoHandler.UploadVideo)
				protected.PUT("/:video_id", videoHandler.UpdateVideo)
				protected.DELETE("/:video_id", videoHandler.DeleteVideo)
				protected.GET("/user/:user_id", videoHandler.ListUserVideos)
				protected.POST("/:video_id/sessions", videoHandler.CreateSession)
				protected.PUT("/sessions/:session_id/progress", videoHandler.UpdateSessionProgress)
				protected.POST("/sessions/:session_id/network", videoHandler.UpdateNetworkStatus)
				protected.GET("/:video_id/analytics", videoHandler.GetVideoAnalytics)
			}

			// Webhook endpoints (no auth for external services)
			videos.POST("/webhooks/cloudflare", videoHandler.CloudflareWebhook)

			// WebSocket endpoints
			videos.GET("/ws/:session_id", wsHandler.HandleWebSocket)
			videos.GET("/ws/stats", wsHandler.GetWebSocketStats)
			videos.POST("/ws/broadcast", middleware.AuthMiddleware(cfg.JWTSecret), wsHandler.BroadcastMessage)
			videos.POST("/ws/session/:session_id/send", middleware.AuthMiddleware(cfg.JWTSecret), wsHandler.SendToSession)
			videos.GET("/ws/session/:session_id", middleware.AuthMiddleware(cfg.JWTSecret), wsHandler.GetSessionInfo)
		}
	}

	// Start server
	serverAddr := cfg.Host + ":" + cfg.Port
	log.Printf("Video service starting on %s", serverAddr)

	server := &http.Server{
		Addr:    serverAddr,
		Handler: router,
	}

	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initDatabase(cfg *config.Config) (*sql.DB, error) {
	dbURL := "host=" + cfg.DBHost + " port=" + cfg.DBPort + " user=" + cfg.DBUser + 
		" password=" + cfg.DBPassword + " dbname=" + cfg.DBName + " sslmode=disable"

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Database connection established")
	return db, nil
}