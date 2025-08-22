package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strconv"
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

	// Service connection configurations
	authServiceURL := getEnv("AUTH_SERVICE_URL", "localhost:8081")
	courseServiceURL := getEnv("COURSE_SERVICE_URL", "localhost:8082")
	progressServiceURL := getEnv("PROGRESS_SERVICE_URL", "localhost:8083")
	bucketServiceURL := getEnv("BUCKET_SERVICE_URL", "http://localhost:8085")
	chatbotServiceURL := getEnv("CHATBOT_SERVICE_URL", "http://localhost:8086")
	forumServiceURL := getEnv("FORUM_SERVICE_URL", "http://localhost:8087")
	paymentServiceURL := getEnv("PAYMENT_SERVICE_URL", "http://localhost:8088")

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
	courseHandler := handler.NewCourseHandler(courseConn, log)
	progressHandler := handler.NewProgressHandler(progressConn, log)
	videoHandler := handler.NewVideoHandler()
	bucketHandler := handler.NewBucketHandler(bucketServiceURL, log)
	chatbotHandler := handler.NewChatbotHandler(chatbotServiceURL)
	forumHandler := handler.NewForumHandler(forumServiceURL)
	paymentHandler := handler.NewPaymentHandler(paymentServiceURL, log)
	docsHandler := handler.NewDocsHandler()

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(authConn, log)
	loggingMiddleware := middleware.NewLoggingMiddleware(log)
	rateLimitMiddleware := middleware.NewRateLimitMiddleware(log, 100, 200) // 100 req/sec, 200 burst
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
		docsHandler,
		authMiddleware,
		loggingMiddleware,
		rateLimitMiddleware,
		circuitBreakerManager,
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