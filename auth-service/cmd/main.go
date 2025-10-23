package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/redis/go-redis/v9"
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
		Port:     5432,
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

	// Initialize Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     getEnv("REDIS_URL", "redis:6379"),
		Password: getEnv("REDIS_PASSWORD", ""),
		DB:       0,
	})
	
	// Test Redis connection
	ctx := context.Background()
	if err := redisClient.Ping(ctx).Err(); err != nil {
		log.Fatal(fmt.Errorf("failed to connect to Redis: %w", err))
	}
	defer redisClient.Close()
	log.Info("Connected to Redis successfully")

	userRepo := repository.NewUserRepository(db)
	oauthRepo := repository.NewOAuthRepository(db)
	verificationRedisRepo := repository.NewVerificationRedisRepository(redisClient)
	
	// Initialize email service
	emailService := service.NewEmailService(
		getEnv("RESEND_API_KEY", ""),
		getEnv("VERIFICATION_EMAIL_FROM", "noreply@study.khoipn.id.vn"),
		getEnv("BASE_URL", "https://study.khoipn.id.vn"),
	)
	
	// Initialize verification service
	verificationService := service.NewVerificationService(
		verificationRedisRepo,
		userRepo,
		emailService,
		log,
	)
	
	authService := service.NewAuthService(userRepo, jwtSecret, log)
	authService.SetVerificationService(verificationService)
	
	// OAuth configurations - Google only
	oauthConfigs := map[model.OAuthProvider]model.OAuthConfig{
		model.ProviderGoogle: {
			ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
			ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
			RedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "http://localhost:8080/auth/google/callback"),
			Scopes:       []string{"openid", "email", "profile"},
		},
	}
	
	oauthService := service.NewOAuthService(userRepo, oauthRepo, authService, log, oauthConfigs, "http://localhost:3000")
	authHandler := handler.NewAuthHandler(authService, oauthService, verificationService, log)

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