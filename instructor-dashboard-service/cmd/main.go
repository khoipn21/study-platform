package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"instructor-dashboard-service/internal/handler"
	"instructor-dashboard-service/internal/middleware"
	"instructor-dashboard-service/internal/repository"
	"instructor-dashboard-service/internal/service"

	"database/sql"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	port := getEnv("PORT", "8089")
	dbURL := getEnv("DATABASE_URL", "postgres://admin:admin123@localhost:2345/studyplatform?sslmode=disable")

	// Connect to database
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	// Initialize repositories
	dashboardRepo := repository.NewDashboardRepository(db)
	analyticsRepo := repository.NewAnalyticsRepository(db)
	communicationRepo := repository.NewCommunicationRepository(db)

	// Initialize services
	dashboardService := service.NewDashboardService(dashboardRepo, analyticsRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo)
	communicationService := service.NewCommunicationService(communicationRepo)

	// Initialize handlers
	dashboardHandler := handler.NewDashboardHandler(dashboardService)
	analyticsHandler := handler.NewAnalyticsHandler(analyticsService)
	communicationHandler := handler.NewCommunicationHandler(communicationService)
	healthHandler := handler.NewHealthHandler(db)

	// Setup Gin router
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	// Add CORS middleware
	router.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Health check endpoints
	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/health/ready", healthHandler.ReadinessCheck)
	router.GET("/health/live", healthHandler.LivenessCheck)
	router.GET("/metrics", healthHandler.MetricsEndpoint)

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Instructor dashboard routes - apply auth middleware
		instructor := v1.Group("/instructor")
		instructor.Use(middleware.AuthMiddleware())
		{
			// Dashboard overview
			instructor.GET("/dashboard/overview", dashboardHandler.GetDashboardOverview)
			instructor.PUT("/dashboard/settings", dashboardHandler.UpdateDashboardSettings)

			// Course management
			instructor.GET("/courses", dashboardHandler.GetInstructorCourses)
			instructor.GET("/courses/:id/analytics", analyticsHandler.GetCourseAnalytics)
			instructor.POST("/courses/:id/bulk-operations", dashboardHandler.BulkCourseOperations)

			// Analytics routes
			instructor.GET("/analytics/revenue", analyticsHandler.GetRevenueAnalytics)
			instructor.GET("/analytics/engagement", analyticsHandler.GetEngagementAnalytics)
			instructor.GET("/analytics/students", analyticsHandler.GetStudentAnalytics)

			// Student management
			instructor.GET("/students", dashboardHandler.GetStudents)
			instructor.GET("/students/:id", dashboardHandler.GetStudentDetails)

			// Video analytics
			instructor.GET("/videos/analytics", analyticsHandler.GetVideoAnalytics)
			instructor.GET("/videos/:id/engagement", analyticsHandler.GetVideoEngagement)

			// Communication
			instructor.POST("/communication/broadcast", communicationHandler.SendBroadcast)
			instructor.GET("/communication/history", communicationHandler.GetCommunicationHistory)
			instructor.POST("/communication/automated", communicationHandler.SetupAutomatedMessages)

			// AI suggestions
			instructor.GET("/suggestions", dashboardHandler.GetAISuggestions)
			instructor.POST("/suggestions/:id/implement", dashboardHandler.ImplementSuggestion)

			// Team management
			instructor.GET("/team", dashboardHandler.GetTeamMembers)
			instructor.POST("/team/invite", dashboardHandler.InviteTeamMember)
			instructor.PUT("/team/:id", dashboardHandler.UpdateTeamMember)
			instructor.DELETE("/team/:id", dashboardHandler.RemoveTeamMember)

			// Notifications
			instructor.GET("/notifications", communicationHandler.GetNotifications)
			instructor.GET("/notifications/settings", dashboardHandler.GetNotificationSettings)
			instructor.PUT("/notifications/settings", dashboardHandler.UpdateNotificationSettings)
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Starting instructor dashboard service on port %s", port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests 30 seconds to complete
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}