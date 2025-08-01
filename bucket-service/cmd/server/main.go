package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"bucket-service/internal/config"
	"bucket-service/internal/handler"
	"bucket-service/internal/middleware"
	"bucket-service/internal/repository"
	"bucket-service/internal/service"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/golang-migrate/migrate/v4"
	migratepostgres "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Set up logging
	if cfg.Server.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Connect to database
	db, err := connectDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	// Run migrations
	if err := runMigrations(cfg); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize services
	s3Service, err := service.NewS3Service(cfg)
	if err != nil {
		log.Fatal("Failed to initialize S3 service:", err)
	}

	imageService := service.NewImageService()
	fileRepo := repository.NewFileRepository(db)

	// Initialize handlers
	fileHandler := handler.NewFileHandler(fileRepo, s3Service, imageService, cfg)
	healthHandler := handler.NewHealthHandler(db)

	// Initialize auth middleware
	authMiddleware := middleware.NewAuthMiddleware(cfg.Security.JWTSecret)

	// Set up router
	router := setupRouter(cfg, fileHandler, healthHandler, authMiddleware)

	// Create server
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler: router,
	}

	// Start server
	go func() {
		log.Printf("Starting bucket service on %s:%s", cfg.Server.Host, cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("Failed to start server:", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
}

func connectDB(cfg *config.Config) (*gorm.DB, error) {
	var gormLogger logger.Interface
	if cfg.Server.Env == "production" {
		gormLogger = logger.Default.LogMode(logger.Silent)
	} else {
		gormLogger = logger.Default.LogMode(logger.Info)
	}

	db, err := gorm.Open(postgres.Open(cfg.GetDSN()), &gorm.Config{
		Logger: gormLogger,
	})
	if err != nil {
		return nil, err
	}

	// Configure connection pool
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	return db, nil
}

func runMigrations(cfg *config.Config) error {
	db, err := sql.Open("postgres", cfg.GetDSN())
	if err != nil {
		return err
	}
	defer db.Close()

	driver, err := migratepostgres.WithInstance(db, &migratepostgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Database migrations completed")
	return nil
}

func setupRouter(cfg *config.Config, fileHandler *handler.FileHandler, healthHandler *handler.HealthHandler, authMiddleware *middleware.AuthMiddleware) *gin.Engine {
	router := gin.New()

	// Global middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(middleware.RequestID())
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.Security.CORSOrigins,
		AllowMethods:     cfg.Security.CORSMethods,
		AllowHeaders:     cfg.Security.CORSHeaders,
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "ETag", "x-amz-version-id"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Health endpoints (no auth required)
	router.GET("/health", healthHandler.HealthCheck)
	router.GET("/ready", healthHandler.ReadinessCheck)

	// API routes
	api := router.Group("/api")
	{
		// File routes (require authentication)
		files := api.Group("/files")
		files.Use(authMiddleware.RequireAuth())
		{
			files.POST("/upload", fileHandler.UploadFile)
			files.POST("/upload/start", fileHandler.StartMultipartUpload)
			files.POST("/upload/sessions/:sessionId/parts", fileHandler.GetPartURLs)
			files.POST("/upload/complete/:sessionId", fileHandler.CompleteMultipartUpload)
			files.DELETE("/upload/:sessionId", fileHandler.AbortMultipartUpload)
			
			files.GET("/:fileId", fileHandler.GetFile)
			files.GET("/:fileId/metadata", fileHandler.GetFileMetadata)
			files.DELETE("/:fileId", fileHandler.DeleteFile)
			files.GET("", fileHandler.ListFiles)
		}

		// Public file routes (no auth required)
		api.GET("/files/public/:fileId", fileHandler.GetPublicFile)
	}

	return router
}