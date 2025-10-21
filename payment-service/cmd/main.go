package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/study-platform/payment-service/internal/config"
	"github.com/study-platform/payment-service/internal/handler"
	"github.com/study-platform/payment-service/internal/lemonsqueezy"
	"github.com/study-platform/payment-service/internal/repository"
	"github.com/study-platform/payment-service/internal/service"
	"github.com/study-platform/payment-service/internal/stripe"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
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
	paymentMethodRepo := repository.NewPaymentMethodRepository(db.DB)
	transactionRepo := repository.NewTransactionRepository(db.DB)
	subscriptionRepo := repository.NewSubscriptionRepository(db.DB)
	webhookEventRepo := repository.NewWebhookEventRepository(db)

	// Initialize Stripe repositories
	stripeCustomerRepo := repository.NewStripeCustomerRepository(db)
	stripeProductRepo := repository.NewStripeProductRepository(db)
	stripeWebhookRepo := repository.NewStripeWebhookRepository(db)

	// Initialize other repositories needed for Stripe
	courseRepo := repository.NewCourseRepository(db.DB)
	enrollmentRepo := repository.NewEnrollmentRepository(db.DB)

	// Initialize Progress service client
	progressClient, err := service.NewProgressClient(cfg.Services.ProgressServiceURL)
	if err != nil {
		log.Printf("Warning: Failed to connect to progress service: %v", err)
		// Don't fail the service, just log the warning
	}

	// Initialize Lemon Squeezy client
	lemonSqueezyClient := lemonsqueezy.NewClient(lemonsqueezy.Config{
		APIKey:        cfg.Payment.LemonSqueezyAPIKey,
		StoreID:       cfg.Payment.LemonSqueezyStoreID,
		Environment:   "production", // Use "test" for testing
		WebhookSecret: cfg.Payment.LemonSqueezyWebhookSecret,
	})

	// Initialize Stripe client
	var stripeClient *stripe.Client
	var stripeService *service.StripeService
	if cfg.Payment.StripeSecretKey != "" {
		stripeClient, err = stripe.NewClient(cfg.Payment.StripeSecretKey)
		if err != nil {
			log.Printf("Warning: Failed to initialize Stripe client: %v", err)
		} else {
			stripeService = service.NewStripeService(
				stripeClient,
				cfg,
				stripeCustomerRepo,
				stripeProductRepo,
				stripeWebhookRepo,
				transactionRepo,
				courseRepo,
				enrollmentRepo,
			)
			log.Println("Stripe client initialized successfully")
		}
	} else {
		log.Println("Stripe secret key not configured, Stripe payments disabled")
	}

	// Initialize services
	paymentService := service.NewPaymentService(paymentMethodRepo, transactionRepo, subscriptionRepo, webhookEventRepo, cfg.Payment, progressClient)

	// Initialize handlers
	paymentHandler := handler.NewPaymentHandler(paymentService)
	lemonSqueezyHandler := handler.NewLemonSqueezyHandler(
		lemonSqueezyClient,
		transactionRepo,
		webhookEventRepo,
		paymentService,
		progressClient,
		cfg.Payment.LemonSqueezyVariantID,
	)

	// Initialize Stripe handler if Stripe is configured
	var stripeHandler *handler.StripeHandler
	if stripeService != nil {
		stripeHandler = handler.NewStripeHandler(
			stripeService,
			stripeWebhookRepo,
			cfg,
		)
		log.Println("Stripe handler initialized successfully")
	}

	// Setup router
	router := setupRouter(paymentHandler, lemonSqueezyHandler, stripeHandler, cfg)

	// Start server
	addr := fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port)
	log.Printf("Payment service starting on %s", addr)
	
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func connectDB(cfg config.DatabaseConfig) (*sqlx.DB, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode)

	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func setupRouter(paymentHandler *handler.PaymentHandler, lemonSqueezyHandler *handler.LemonSqueezyHandler, stripeHandler *handler.StripeHandler, cfg *config.Config) *gin.Engine {
	router := gin.Default()

	// CORS middleware - configured via environment variable
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORS.AllowedOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization", "X-Requested-With", "Accept", "Accept-Language", "X-User-ID", "X-User-Role"},
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

		// Lemon Squeezy endpoints
		api.POST("/lemonsqueezy/products", lemonSqueezyHandler.CreateProduct)
		api.POST("/lemonsqueezy/checkout/course/:course_id", lemonSqueezyHandler.CreateCheckout)
		api.POST("/lemonsqueezy/verify/:order_id", lemonSqueezyHandler.VerifyPayment)
		api.GET("/lemonsqueezy/products", lemonSqueezyHandler.GetProducts)
		api.GET("/lemonsqueezy/variants", lemonSqueezyHandler.GetVariants)

		// Stripe endpoints (if available)
		if stripeHandler != nil {
			stripe := api.Group("/stripe")
			{
				stripe.POST("/payment-intents", stripeHandler.CreatePaymentIntent)
				stripe.GET("/payment-intents/:payment_intent_id", stripeHandler.GetPaymentIntent)
				stripe.POST("/payment-intents/confirm", stripeHandler.ConfirmPaymentIntent)
				stripe.GET("/transactions", stripeHandler.ListTransactions)
				stripe.POST("/purchase/course/:courseId", stripeHandler.PurchaseCourse)
			}
		}
	}

	// Public endpoints (no auth required)
	// Stripe config endpoint (public)
	if stripeHandler != nil {
		router.GET("/api/v1/stripe/config", stripeHandler.GetConfig)
	}

	// Webhooks (no auth required)
	router.POST("/api/v1/payments/lemonsqueezy/webhook", lemonSqueezyHandler.HandleWebhook)

	// Stripe webhook (no auth required)
	if stripeHandler != nil {
		router.POST("/api/v1/payments/stripe/webhook", stripeHandler.HandleWebhook)
	}

	return router
}