package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	"github.com/study-platform/auth-service/internal/handler"
	"github.com/study-platform/auth-service/internal/model"
	"github.com/study-platform/auth-service/internal/repository"
	"github.com/study-platform/auth-service/internal/service"
	"github.com/study-platform/pkg/database"
	"github.com/study-platform/pkg/logger"
	pb "github.com/study-platform/auth-service/proto"
)

func main() {
	log := logger.New()
	log.Info("Starting Auth Service...")

	jwtSecret := getEnv("JWT_SECRET", "your-secret-key-here")
	grpcPort := getEnv("GRPC_PORT", "8081")

	dbConfig := database.Config{
		Host:     getEnv("DB_HOST", "postgres"),
		Port:     getEnvAsInt("DB_PORT", 5432),
		User:     getEnv("DB_USER", "admin"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "studyplatform"),
		SSLMode:  getEnv("DB_SSL_MODE", "disable"),
	}

	db, err := database.New(dbConfig)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to connect to database: %w", err))
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	oauthRepo := repository.NewOAuthRepository(db)
	authService := service.NewAuthService(userRepo, jwtSecret, log)
	
	// OAuth configurations
	oauthConfigs := map[model.OAuthProvider]model.OAuthConfig{
		model.ProviderGoogle: {
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
			Scopes:       []string{"openid", "email", "profile"},
		},
		model.ProviderGitHub: {
			ClientID:     getEnv("GITHUB_CLIENT_ID", ""),
			ClientSecret: getEnv("GITHUB_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GITHUB_REDIRECT_URL", "http://localhost:8080/auth/github/callback"),
			Scopes:       []string{"user:email"},
		},
		model.ProviderFacebook: {
			ClientID:     getEnv("FACEBOOK_CLIENT_ID", ""),
			ClientSecret: getEnv("FACEBOOK_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("FACEBOOK_REDIRECT_URL", "http://localhost:8080/auth/facebook/callback"),
			Scopes:       []string{"email", "public_profile"},
		},
	}
	
	oauthService := service.NewOAuthService(userRepo, oauthRepo, authService, log, oauthConfigs, "http://localhost:3000")
	authHandler := handler.NewAuthHandler(authService, oauthService, log)

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to listen: %w", err))
	}

	s := grpc.NewServer()
	pb.RegisterAuthServiceServer(s, authHandler)

	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthServer)
	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(s)

	log.Infof("Auth Service listening on port %s", grpcPort)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(fmt.Errorf("failed to serve: %w", err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Auth Service...")
	s.GracefulStop()
	log.Info("Auth Service stopped")
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