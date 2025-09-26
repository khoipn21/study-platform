package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/study-platform/api-gateway/internal/handler"
	"github.com/study-platform/api-gateway/internal/middleware"
	"github.com/study-platform/api-gateway/internal/router"
	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Initialize logger
	log := logger.New()
	log.Info("Starting API Gateway...")

	// CRITICAL: Validate environment variables first
	if err := middleware.ValidateEnvVars(); err != nil {
		log.Fatalf("Environment validation failed: %v", err)
	}
	log.Info("Environment variables validated successfully")

	// Service connection configurations
	authServiceURL := getEnv("AUTH_SERVICE_URL", "localhost:8081")
	courseServiceURL := getEnv("COURSE_SERVICE_URL", "localhost:8082")
	progressServiceURL := getEnv("PROGRESS_SERVICE_URL", "localhost:8083")
	bucketServiceURL := getEnv("BUCKET_SERVICE_URL", "http://bucket-service:8085")
	chatbotServiceURL := getEnv("CHATBOT_SERVICE_URL", "http://localhost:8086")
	forumServiceURL := getEnv("FORUM_SERVICE_URL", "http://localhost:8087")
	paymentServiceURL := getEnv("PAYMENT_SERVICE_URL", "http://localhost:8088")
	instructorDashboardServiceURL := getEnv("INSTRUCTOR_DASHBOARD_SERVICE_URL", "http://localhost:8089")

	// Initialize circuit breaker and retry managers
	circuitBreakerConfig := middleware.CircuitBreakerConfig{
		MaxFailures:  5,
		Timeout:      10 * time.Second,
		ResetTimeout: 60 * time.Second,
	}
	circuitBreakerManager := middleware.NewCircuitBreakerManager(circuitBreakerConfig, log)
	
	retryManager := middleware.NewRetryManager(log)
	retryManager.SetConfig("auth-service", middleware.DefaultRetryConfig())
	retryManager.SetConfig("course-service", middleware.DefaultRetryConfig())
	retryManager.SetConfig("progress-service", middleware.DefaultRetryConfig())
	retryManager.SetConfig("bucket-service", middleware.DefaultRetryConfig())
	retryManager.SetConfig("instructor-dashboard", middleware.DefaultRetryConfig())

	// Connect to auth service with circuit breaker and retry
	authCB := circuitBreakerManager.GetCircuitBreaker("auth-service")
	authConn, err := grpc.NewClient(authServiceURL, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.CircuitBreakerInterceptor(authCB)),
		grpc.WithUnaryInterceptor(retryManager.GetInterceptor("auth-service")))
	if err != nil {
		log.Errorf("Failed to connect to auth service: %v", err)
		os.Exit(1)
	}
	defer authConn.Close()

	// Connect to course service with circuit breaker and retry
	courseCB := circuitBreakerManager.GetCircuitBreaker("course-service")
	courseConn, err := grpc.NewClient(courseServiceURL, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.CircuitBreakerInterceptor(courseCB)),
		grpc.WithUnaryInterceptor(retryManager.GetInterceptor("course-service")))
	if err != nil {
		log.Errorf("Failed to connect to course service: %v", err)
		os.Exit(1)
	}
	defer courseConn.Close()

	// Connect to progress service with circuit breaker and retry
	progressCB := circuitBreakerManager.GetCircuitBreaker("progress-service")
	progressConn, err := grpc.NewClient(progressServiceURL, 
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(middleware.CircuitBreakerInterceptor(progressCB)),
		grpc.WithUnaryInterceptor(retryManager.GetInterceptor("progress-service")))
	if err != nil {
		log.Errorf("Failed to connect to progress service: %v", err)
		os.Exit(1)
	}
	defer progressConn.Close()

	log.Info("Connected to all gRPC services")

	// Initialize handlers
	authHandler := handler.NewAuthHandler(authConn, log)
	courseHandler := handler.NewCourseHandler(courseConn, bucketServiceURL, log)
	progressHandler := handler.NewProgressHandler(progressConn, log)
	videoHandler := handler.NewVideoHandlerWithCourse(courseConn, log)
	bucketHandler := handler.NewBucketHandler(bucketServiceURL, log)
	chatbotHandler := handler.NewChatbotHandler(chatbotServiceURL)
	forumHandler := handler.NewForumHandler(forumServiceURL)
	paymentHandler := handler.NewPaymentHandler(paymentServiceURL, log)
	lemonSqueezyHandler := handler.NewLemonSqueezyHandler(log)
	instructorDashboardHandler := handler.NewInstructorDashboardHandler(instructorDashboardServiceURL, courseHandler)
	studentDashboardHandler := handler.NewStudentDashboardHandler()
	docsHandler := handler.NewDocsHandler()

	// Initialize new course access handlers
	courseAccessHandler := handler.NewCourseAccessHandler(courseConn, log)
	progressTrackingHandler := handler.NewProgressTrackingHandler(progressConn, log)
	notesHandler := handler.NewNotesHandler("http://course-service:8092", log)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddlewareWithCourse(authConn, courseConn, log)
	loggingMiddleware := middleware.NewLoggingMiddleware(log)
	
	// Security Configuration
	securityConfig := middleware.SecurityConfig{
		JWTSecret:             os.Getenv("JWT_SECRET"),
		ContentSecurityPolicy: getEnv("CONTENT_SECURITY_POLICY", "default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'"),
		EnableSecurityHeaders: getEnvAsBool("SECURITY_HEADERS_ENABLED", true),
		AllowedOrigins:       strings.Split(getEnv("CORS_ORIGINS", "http://localhost:3000"), ","),
	}
	
	securityMiddleware := middleware.NewSecurityMiddleware(securityConfig, log)
	
	// Rate limiting configuration
	rateLimit := getEnvAsInt("RATE_LIMIT_REQUESTS", 100)
	rateBurst := getEnvAsInt("RATE_LIMIT_BURST", 200)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(log, rateLimit, rateBurst)
	rateLimitMiddleware.StartCleanupRoutine()

	// Initialize router
	rt := router.NewRouter(
		authHandler,
		courseHandler,
		progressHandler,
		videoHandler,
		bucketHandler,
		chatbotHandler,
		forumHandler,
		paymentHandler,
		lemonSqueezyHandler,
		instructorDashboardHandler,
		studentDashboardHandler,
		docsHandler,
		courseAccessHandler,
		progressTrackingHandler,
		notesHandler,
		authMiddleware,
		loggingMiddleware,
		rateLimitMiddleware,
		circuitBreakerManager,
		securityMiddleware,
	)

	// Setup routes
	routes := rt.SetupRoutes()

	// HTTP server configuration
	port := getEnv("HTTP_PORT", "8080")
	server := &http.Server{
		Addr:         ":" + port,
		Handler:      routes,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Infof("API Gateway listening on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("Failed to start server: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down API Gateway...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Errorf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	log.Info("API Gateway shutdown complete")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		return strings.ToLower(value) == "true"
	}
	return defaultValue
}