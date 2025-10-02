package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"forum-service/internal/config"
	"forum-service/internal/handler"
	"forum-service/internal/repository"
	"forum-service/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	forumRepo := repository.NewForumRepository(db)

	// Initialize services
	forumService := service.NewForumService(forumRepo)

	// Initialize handlers
	forumHandler := handler.NewForumHandler(forumService)

	// Setup router
	router := setupRouter(forumHandler)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Forum service starting on %s", addr)
	
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

func setupRouter(forumHandler *handler.ForumHandler) *gin.Engine {
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
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "forum"})
	})

	// Authentication middleware (simplified - in production, this should validate JWT tokens)
	_ = func(c *gin.Context) {
		// For now, we'll expect user_id in headers
		userIDStr := c.GetHeader("X-User-ID")
		userRole := c.GetHeader("X-User-Role")
		
		if userIDStr != "" {
			c.Set("user_id", userIDStr)
		}
		if userRole != "" {
			c.Set("user_role", userRole)
		}
		
		c.Next()
	}

	// Optional auth middleware for endpoints that can work with or without auth
	optionalAuthMiddleware := func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-ID")
		userRole := c.GetHeader("X-User-Role")
		
		if userIDStr != "" {
			c.Set("user_id", userIDStr)
		}
		if userRole != "" {
			c.Set("user_role", userRole)
		}
		
		c.Next()
	}

	// Required auth middleware
	requiredAuthMiddleware := func(c *gin.Context) {
		userIDStr := c.GetHeader("X-User-ID")
		if userIDStr == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication required"})
			c.Abort()
			return
		}
		
		userRole := c.GetHeader("X-User-Role")
		c.Set("user_id", userIDStr)
		if userRole != "" {
			c.Set("user_role", userRole)
		}
		
		c.Next()
	}

	// API routes
	api := router.Group("/api/v1")
	{
		// Public endpoints (can work without auth but benefit from it)
		api.GET("/topics", optionalAuthMiddleware, forumHandler.ListTopics)
		api.GET("/topics/:topicId", optionalAuthMiddleware, forumHandler.GetTopic)
		api.GET("/topics/:topicId/posts", optionalAuthMiddleware, forumHandler.ListPosts)
		api.GET("/posts/:postId", optionalAuthMiddleware, forumHandler.GetPost)
		api.GET("/search", optionalAuthMiddleware, forumHandler.SearchTopics)

		// Authenticated endpoints
		authenticated := api.Group("")
		authenticated.Use(requiredAuthMiddleware)
		{
			// Topic management
			authenticated.POST("/topics", forumHandler.CreateTopic)
			authenticated.PUT("/topics/:topicId", forumHandler.UpdateTopic)
			authenticated.DELETE("/topics/:topicId", forumHandler.DeleteTopic)
			authenticated.PUT("/topics/:topicId/sticky", forumHandler.ToggleTopicSticky)
			authenticated.PUT("/topics/:topicId/lock", forumHandler.ToggleTopicLock)

			// Post management
			authenticated.POST("/posts", forumHandler.CreatePost)
			authenticated.PUT("/posts/:postId", forumHandler.UpdatePost)
			authenticated.DELETE("/posts/:postId", forumHandler.DeletePost)
			authenticated.PUT("/posts/:postId/answer", forumHandler.MarkPostAsAnswer)
			authenticated.PUT("/posts/:postId/pin", forumHandler.TogglePostPin)

			// Voting
			authenticated.POST("/votes", forumHandler.VotePost)
			authenticated.DELETE("/posts/:postId/vote", forumHandler.RemoveVote)

			// Approval management (instructors/admins only)
			authenticated.PUT("/topics/:topicId/approve", forumHandler.ApproveTopic)
			authenticated.PUT("/posts/:postId/approve", forumHandler.ApprovePost)
			authenticated.GET("/pending/topics", forumHandler.GetPendingTopics)

			// Pin management (instructors/admins only)
			authenticated.PUT("/topics/:topicId/pin-order", forumHandler.SetTopicPinOrder)
			authenticated.PUT("/posts/:postId/pin-order", forumHandler.SetPostPinOrder)
		}
	}

	return router
}