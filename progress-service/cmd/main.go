package main

import (
	"context"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/study-platform/progress-service/internal/handler"
	"github.com/study-platform/progress-service/internal/repository"
	"github.com/study-platform/progress-service/internal/service"
	pb "github.com/study-platform/progress-service/proto"
	"github.com/study-platform/pkg/database"
	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize logger
	log := logger.New()
	log.Info("Starting Progress Service...")

	// Database configuration
	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "postgres"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "admin"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "studyplatform"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	// Connect to database
	db, err := database.New(dbConfig)
	if err != nil {
		log.Errorf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer db.Close()

	log.Info("Database connected successfully")

	// Initialize repositories
	progressRepo := repository.NewProgressRepository(db)

	// Initialize services
	progressService := service.NewProgressService(progressRepo, log)

	// Initialize handlers
	progressHandler := handler.NewProgressHandler(progressService, log)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterProgressServiceServer(grpcServer, progressHandler)

	// Enable reflection for development
	reflection.Register(grpcServer)

	// Start server
	port := getEnv("GRPC_PORT", "8080")
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Errorf("Failed to listen on port %s: %v", port, err)
		os.Exit(1)
	}

	log.Infof("Progress Service listening on port %s", port)

	// Start server in a goroutine
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Errorf("Failed to serve gRPC server: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Progress Service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	done := make(chan bool)
	go func() {
		grpcServer.GracefulStop()
		done <- true
	}()

	select {
	case <-ctx.Done():
		log.Info("Shutdown timeout, forcing exit")
		grpcServer.Stop()
	case <-done:
		log.Info("Progress Service shutdown complete")
	}
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