package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/study-platform/course-service/internal/handler"
	"github.com/study-platform/course-service/internal/repository"
	"github.com/study-platform/course-service/internal/service"
	pb "github.com/study-platform/course-service/proto"
	"github.com/study-platform/pkg/database"
	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Initialize logger
	log := logger.New()
	log.Info("Starting Course Service...")

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
	courseRepo := repository.NewCourseRepository(db)
	lectureRepo := repository.NewLectureRepository(db)
	enrollmentRepo := repository.NewEnrollmentRepository(db)
	courseResourceRepo := repository.NewCourseResourceRepository(db.GetSqlxDB(), log)
	notesRepo := repository.NewNotesRepository(db.GetSqlxDB())

	// Initialize services
	courseService := service.NewCourseService(courseRepo, lectureRepo, enrollmentRepo, courseResourceRepo, log)
	notesService := service.NewNotesService(notesRepo, log)

	// Initialize handlers
	courseHandler := handler.NewCourseHandler(courseService, log)
	accessValidationHandler := handler.NewAccessValidationHandler(courseService, log)
	notesHandler := handler.NewNotesHandler(notesService, log)

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterCourseServiceServer(grpcServer, courseHandler)
	pb.RegisterAccessValidationServiceServer(grpcServer, accessValidationHandler)

	// Enable reflection for development
	reflection.Register(grpcServer)

	// Setup HTTP server for notes API
	gin.SetMode(gin.ReleaseMode)
	httpRouter := gin.New()
	httpRouter.Use(gin.Recovery())
	httpRouter.Use(func(c *gin.Context) {
		// Add user authentication middleware
		userID := c.GetHeader("X-User-ID")
		if userID != "" {
			c.Set("user_id", userID)
		}
		c.Next()
	})

	// Notes routes
	httpRouter.POST("/courses/:course_id/lectures/:lecture_id/notes", notesHandler.CreateNote)
	httpRouter.GET("/courses/:course_id/lectures/:lecture_id/notes", notesHandler.GetNotesByLecture)
	httpRouter.GET("/courses/:course_id/notes", notesHandler.GetNotesByCourse)
	httpRouter.GET("/notes/:note_id", notesHandler.GetNote)
	httpRouter.PUT("/notes/:note_id", notesHandler.UpdateNote)
	httpRouter.DELETE("/notes/:note_id", notesHandler.DeleteNote)

	// Start gRPC server
	grpcPort := getEnv("GRPC_PORT", "8082")
	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Errorf("Failed to listen on gRPC port %s: %v", grpcPort, err)
		os.Exit(1)
	}

	log.Infof("Course Service gRPC listening on port %s", grpcPort)

	// Start gRPC server in a goroutine
	go func() {
		if err := grpcServer.Serve(listener); err != nil {
			log.Errorf("Failed to serve gRPC server: %v", err)
			os.Exit(1)
		}
	}()

	// Start HTTP server for notes API
	httpPort := getEnv("HTTP_PORT", "8092")
	httpServer := &http.Server{
		Addr:    ":" + httpPort,
		Handler: httpRouter,
	}

	log.Infof("Course Service HTTP listening on port %s", httpPort)

	// Start HTTP server in a goroutine
	go func() {
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Errorf("Failed to serve HTTP server: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info("Shutting down Course Service...")

	// Graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	done := make(chan bool)
	go func() {
		// Shutdown HTTP server
		httpServer.Shutdown(ctx)
		// Shutdown gRPC server
		grpcServer.GracefulStop()
		done <- true
	}()

	select {
	case <-ctx.Done():
		log.Info("Shutdown timeout, forcing exit")
		grpcServer.Stop()
	case <-done:
		log.Info("Course Service shutdown complete")
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