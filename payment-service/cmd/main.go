package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"payment-service/internal/config"
	"payment-service/internal/handler"
	"payment-service/internal/repository"
	"payment-service/internal/service"

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
	paymentMethodRepo := repository.NewPaymentMethodRepository(db)
	transactionRepo := repository.NewTransactionRepository(db)
	subscriptionRepo := repository.NewSubscriptionRepository(db)

	// Initialize Progress service client
	progressClient, err := service.NewProgressClient(cfg.Services.ProgressServiceURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to progress service: %v", err)
		// Don't fail the service, just log the warning
	}

	// Initialize services
	paymentService := service.NewPaymentService(paymentMethodRepo, transactionRepo, subscriptionRepo, cfg.Payment, progressClient)

	// Initialize handlers
	paymentHandler := handler.NewPaymentHandler(paymentService)

	// Setup router
	router := setupRouter(paymentHandler)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Payment service starting on %s", addr)
	
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

func setupRouter(paymentHandler *handler.PaymentHandler) *gin.Engine {
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
		c.JSON(http.StatusOK, gin.H{"status": "healthy", "service": "payment"})
	})

	// Authentication middleware
	authMiddleware := func(c *gin.Context) {
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
	api.Use(authMiddleware)
	{
		// Payment Methods
		api.POST("/payment-methods", paymentHandler.CreatePaymentMethod)
		api.GET("/payment-methods", paymentHandler.GetPaymentMethods)
		api.PUT("/payment-methods/:methodId", paymentHandler.UpdatePaymentMethod)
		api.DELETE("/payment-methods/:methodId", paymentHandler.DeletePaymentMethod)
		api.PUT("/payment-methods/:methodId/default", paymentHandler.SetDefaultPaymentMethod)

		// Course Purchase
		api.POST("/purchase/course/:courseId", paymentHandler.PurchaseCourse)
		api.POST("/purchase/validate", paymentHandler.ValidatePayment)

		// Transactions
		api.GET("/transactions", paymentHandler.GetTransactions)
		api.GET("/transactions/:transactionId", paymentHandler.GetTransaction)
		api.POST("/transactions/:transactionId/refund", paymentHandler.RefundTransaction)

		// Subscriptions
		api.POST("/subscriptions", paymentHandler.CreateSubscription)
		api.GET("/subscriptions", paymentHandler.GetSubscriptions)
		api.PUT("/subscriptions/:subscriptionId", paymentHandler.UpdateSubscription)
		api.DELETE("/subscriptions/:subscriptionId", paymentHandler.CancelSubscription)

		// Enrollments
		api.GET("/enrollments", paymentHandler.GetUserEnrollments)
		api.GET("/enrollments/check/:courseId", paymentHandler.CheckEnrollment)

		// Webhooks (no auth required)
		router.POST("/webhooks/stripe", paymentHandler.HandleStripeWebhook)
		router.POST("/webhooks/paypal", paymentHandler.HandlePayPalWebhook)
	}

	return router
}