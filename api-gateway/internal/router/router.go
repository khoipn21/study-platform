package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/study-platform/api-gateway/internal/handler"
	"github.com/study-platform/api-gateway/internal/middleware"
)

type Router struct {
	authHandler       *handler.AuthHandler
	courseHandler     *handler.CourseHandler
	progressHandler   *handler.ProgressHandler
	bucketHandler     *handler.BucketHandler
	chatbotHandler    *handler.ChatbotHandler
	forumHandler      *handler.ForumHandler
	docsHandler       *handler.DocsHandler
	authMiddleware    *middleware.AuthMiddleware
	loggingMiddleware *middleware.LoggingMiddleware
	rateLimitMiddleware *middleware.RateLimitMiddleware
	circuitBreakerManager *middleware.CircuitBreakerManager
}

func NewRouter(
	authHandler *handler.AuthHandler,
	courseHandler *handler.CourseHandler,
	progressHandler *handler.ProgressHandler,
	bucketHandler *handler.BucketHandler,
	chatbotHandler *handler.ChatbotHandler,
	forumHandler *handler.ForumHandler,
	docsHandler *handler.DocsHandler,
	authMiddleware *middleware.AuthMiddleware,
	loggingMiddleware *middleware.LoggingMiddleware,
	rateLimitMiddleware *middleware.RateLimitMiddleware,
	circuitBreakerManager *middleware.CircuitBreakerManager,
) *Router {
	return &Router{
		authHandler:       authHandler,
		courseHandler:     courseHandler,
		progressHandler:   progressHandler,
		bucketHandler:     bucketHandler,
		chatbotHandler:    chatbotHandler,
		forumHandler:      forumHandler,
		docsHandler:       docsHandler,
		authMiddleware:    authMiddleware,
		loggingMiddleware: loggingMiddleware,
		rateLimitMiddleware: rateLimitMiddleware,
		circuitBreakerManager: circuitBreakerManager,
	}
}

func (rt *Router) SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	// Apply global middleware
	r.Use(middleware.CORSMiddleware)
	r.Use(middleware.SetJSONContentType)
	r.Use(rt.rateLimitMiddleware.RateLimit)
	r.Use(rt.loggingMiddleware.LogRequest)

	// API version prefix
	api := r.PathPrefix("/api/v1").Subrouter()

	// Health check endpoint
	api.HandleFunc("/health", rt.healthCheck).Methods("GET")
	api.HandleFunc("/health/circuit-breakers", rt.circuitBreakerStatus).Methods("GET")
	
	// Documentation endpoints
	api.HandleFunc("/docs/openapi.json", rt.docsHandler.GetAPISpec).Methods("GET")
	api.HandleFunc("/docs", rt.docsHandler.GetSwaggerUI).Methods("GET")
	
	// Debug endpoint for testing
	api.HandleFunc("/debug/files", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "Debug endpoint works", "path": "` + r.URL.Path + `"}`))
	}).Methods("GET")
	
	// Direct file routes for testing
	api.HandleFunc("/files", rt.bucketHandler.ListFiles).Methods("GET")
	api.HandleFunc("/files/upload", rt.bucketHandler.UploadFile).Methods("POST")

	// Auth routes (no authentication required)
	authRoutes := api.PathPrefix("/auth").Subrouter()
	authRoutes.HandleFunc("/register", rt.authHandler.Register).Methods("POST")
	authRoutes.HandleFunc("/login", rt.authHandler.Login).Methods("POST")
	authRoutes.HandleFunc("/validate", rt.authHandler.ValidateToken).Methods("POST")
	authRoutes.HandleFunc("/oauth/{provider}/url", rt.authHandler.GetOAuthURL).Methods("GET")
	authRoutes.HandleFunc("/oauth/{provider}/callback", rt.authHandler.OAuthCallback).Methods("GET")

	// Protected auth routes (authentication required)
	protectedAuthRoutes := api.PathPrefix("/auth").Subrouter()
	protectedAuthRoutes.Use(rt.authMiddleware.RequireAuth)
	protectedAuthRoutes.HandleFunc("/profile", rt.authHandler.GetProfile).Methods("GET")

	// Course routes (some public, some protected)
	courseRoutes := api.PathPrefix("/courses").Subrouter()
	
	// Public course routes
	courseRoutes.HandleFunc("", rt.courseHandler.ListCourses).Methods("GET")
	courseRoutes.HandleFunc("/search", rt.courseHandler.SearchCourses).Methods("GET")
	courseRoutes.HandleFunc("/{id}", rt.courseHandler.GetCourse).Methods("GET")
	courseRoutes.HandleFunc("/{course_id}/lectures", rt.courseHandler.ListLectures).Methods("GET")
	courseRoutes.HandleFunc("/lectures/{id}", rt.courseHandler.GetLecture).Methods("GET")

	// Protected course routes (authentication required)
	protectedCourseRoutes := api.PathPrefix("/courses").Subrouter()
	protectedCourseRoutes.Use(rt.authMiddleware.RequireAuth)
	protectedCourseRoutes.HandleFunc("", rt.courseHandler.CreateCourse).Methods("POST")
	protectedCourseRoutes.HandleFunc("/{id}", rt.courseHandler.UpdateCourse).Methods("PUT")
	protectedCourseRoutes.HandleFunc("/{id}", rt.courseHandler.DeleteCourse).Methods("DELETE")
	protectedCourseRoutes.HandleFunc("/lectures", rt.courseHandler.CreateLecture).Methods("POST")
	protectedCourseRoutes.HandleFunc("/{course_id}/enroll", rt.courseHandler.EnrollInCourse).Methods("POST")

	// Progress routes (all require authentication)
	progressRoutes := api.PathPrefix("/progress").Subrouter()
	progressRoutes.Use(rt.authMiddleware.RequireAuth)
	
	// Progress tracking
	progressRoutes.HandleFunc("/update", rt.progressHandler.UpdateProgress).Methods("POST")
	progressRoutes.HandleFunc("/courses/{course_id}/lectures/{lecture_id}", rt.progressHandler.GetProgress).Methods("GET")
	progressRoutes.HandleFunc("/courses/{course_id}", rt.progressHandler.GetUserProgress).Methods("GET")
	progressRoutes.HandleFunc("/lectures/{course_id}", rt.progressHandler.GetLectureProgress).Methods("GET")
	progressRoutes.HandleFunc("/courses/{course_id}/completion", rt.progressHandler.GetCourseCompletion).Methods("GET")
	progressRoutes.HandleFunc("/lectures/complete", rt.progressHandler.MarkLectureComplete).Methods("POST")

	// Enrollment routes
	enrollmentRoutes := api.PathPrefix("/enrollments").Subrouter()
	enrollmentRoutes.Use(rt.authMiddleware.RequireAuth)
	enrollmentRoutes.HandleFunc("", rt.progressHandler.CreateEnrollment).Methods("POST")
	enrollmentRoutes.HandleFunc("", rt.progressHandler.ListEnrollments).Methods("GET")
	enrollmentRoutes.HandleFunc("/courses/{course_id}", rt.progressHandler.GetEnrollment).Methods("GET")

	// Analytics routes
	analyticsRoutes := api.PathPrefix("/analytics").Subrouter()
	analyticsRoutes.Use(rt.authMiddleware.RequireAuth)
	analyticsRoutes.HandleFunc("/user", rt.progressHandler.GetUserAnalytics).Methods("GET")

	// File management routes (protected - require authentication)
	protectedFileRoutes := api.PathPrefix("/files").Subrouter()
	protectedFileRoutes.Use(rt.authMiddleware.RequireAuth)
	
	// File operations
	protectedFileRoutes.HandleFunc("/upload", rt.bucketHandler.UploadFile).Methods("POST")
	protectedFileRoutes.HandleFunc("", rt.bucketHandler.ListFiles).Methods("GET")
	protectedFileRoutes.HandleFunc("/{fileId}", rt.bucketHandler.DownloadFile).Methods("GET")
	protectedFileRoutes.HandleFunc("/{fileId}", rt.bucketHandler.DeleteFile).Methods("DELETE")
	protectedFileRoutes.HandleFunc("/{fileId}/metadata", rt.bucketHandler.GetFileMetadata).Methods("GET")
	
	// Multipart upload operations
	protectedFileRoutes.HandleFunc("/upload/start", rt.bucketHandler.StartMultipartUpload).Methods("POST")
	protectedFileRoutes.HandleFunc("/upload/{sessionId}/complete", rt.bucketHandler.CompleteMultipartUpload).Methods("POST")
	protectedFileRoutes.HandleFunc("/upload/{sessionId}", rt.bucketHandler.AbortMultipartUpload).Methods("DELETE")
	protectedFileRoutes.HandleFunc("/upload/{sessionId}/progress", rt.bucketHandler.GetUploadProgress).Methods("GET")

	// Admin routes (require admin role)
	adminRoutes := api.PathPrefix("/admin").Subrouter()
	adminRoutes.Use(rt.authMiddleware.RequireAuth)
	adminRoutes.Use(rt.authMiddleware.RequireAdmin)
	// Add admin-specific routes here as needed

	// Instructor routes (require instructor role)
	instructorRoutes := api.PathPrefix("/instructor").Subrouter()
	instructorRoutes.Use(rt.authMiddleware.RequireAuth)
	instructorRoutes.Use(rt.authMiddleware.RequireInstructor)
	// Add instructor-specific routes here as needed

	// Chatbot routes
	chatRoutes := api.PathPrefix("/chat").Subrouter()
	chatRoutes.Use(rt.authMiddleware.RequireAuth)
	
	// Chat session management
	chatRoutes.HandleFunc("/sessions", rt.chatbotHandler.CreateSession).Methods("POST")
	chatRoutes.HandleFunc("/sessions", rt.chatbotHandler.GetUserSessions).Methods("GET")
	chatRoutes.HandleFunc("/sessions/{sessionId}", rt.chatbotHandler.GetSession).Methods("GET")
	chatRoutes.HandleFunc("/sessions/{sessionId}", rt.chatbotHandler.UpdateSession).Methods("PUT")
	chatRoutes.HandleFunc("/sessions/{sessionId}", rt.chatbotHandler.DeleteSession).Methods("DELETE")
	
	// Chat messaging
	chatRoutes.HandleFunc("/message", rt.chatbotHandler.SendMessage).Methods("POST")
	chatRoutes.HandleFunc("/sessions/{sessionId}/messages", rt.chatbotHandler.GetMessages).Methods("GET")
	
	// WebSocket endpoint for real-time chat
	chatRoutes.HandleFunc("/ws", rt.chatbotHandler.HandleWebSocket).Methods("GET")

	// Forum routes
	forumRoutes := api.PathPrefix("/forum").Subrouter()
	
	// Public forum routes (can work without auth but benefit from it)
	forumRoutes.HandleFunc("/topics", rt.forumHandler.ListTopics).Methods("GET")
	forumRoutes.HandleFunc("/topics/{topicId}", rt.forumHandler.GetTopic).Methods("GET")
	forumRoutes.HandleFunc("/topics/{topicId}/posts", rt.forumHandler.ListPosts).Methods("GET")
	forumRoutes.HandleFunc("/posts/{postId}", rt.forumHandler.GetPost).Methods("GET")
	forumRoutes.HandleFunc("/search", rt.forumHandler.SearchTopics).Methods("GET")
	
	// Course-specific forum routes
	forumRoutes.HandleFunc("/courses/{courseId}/topics", rt.forumHandler.ListCourseTopics).Methods("GET")

	// Protected forum routes (require authentication)
	protectedForumRoutes := api.PathPrefix("/forum").Subrouter()
	protectedForumRoutes.Use(rt.authMiddleware.RequireAuth)
	
	// Topic management
	protectedForumRoutes.HandleFunc("/topics", rt.forumHandler.CreateTopic).Methods("POST")
	protectedForumRoutes.HandleFunc("/topics/{topicId}", rt.forumHandler.UpdateTopic).Methods("PUT")
	protectedForumRoutes.HandleFunc("/topics/{topicId}", rt.forumHandler.DeleteTopic).Methods("DELETE")
	protectedForumRoutes.HandleFunc("/topics/{topicId}/sticky", rt.forumHandler.ToggleTopicSticky).Methods("PUT")
	protectedForumRoutes.HandleFunc("/topics/{topicId}/lock", rt.forumHandler.ToggleTopicLock).Methods("PUT")
	
	// Post management
	protectedForumRoutes.HandleFunc("/posts", rt.forumHandler.CreatePost).Methods("POST")
	protectedForumRoutes.HandleFunc("/posts/{postId}", rt.forumHandler.UpdatePost).Methods("PUT")
	protectedForumRoutes.HandleFunc("/posts/{postId}", rt.forumHandler.DeletePost).Methods("DELETE")
	protectedForumRoutes.HandleFunc("/posts/{postId}/answer", rt.forumHandler.MarkPostAsAnswer).Methods("PUT")
	protectedForumRoutes.HandleFunc("/posts/{postId}/pin", rt.forumHandler.TogglePostPin).Methods("PUT")
	
	// Voting
	protectedForumRoutes.HandleFunc("/votes", rt.forumHandler.VotePost).Methods("POST")
	protectedForumRoutes.HandleFunc("/posts/{postId}/vote", rt.forumHandler.RemoveVote).Methods("DELETE")

	return r
}

func (rt *Router) healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	// Simple JSON encoding
	jsonStr := `{"status":"healthy","service":"api-gateway","version":"1.0.0"}`
	w.Write([]byte(jsonStr))
}

func (rt *Router) circuitBreakerStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	
	status := rt.circuitBreakerManager.GetStatus()
	
	// Simple JSON encoding for circuit breaker status
	jsonStr := `{"status":"healthy","circuit_breakers":{`
	first := true
	for service, state := range status {
		if !first {
			jsonStr += ","
		}
		jsonStr += `"` + service + `":` + `{"state":"` + state.(map[string]interface{})["state"].(string) + `"}`
		first = false
	}
	jsonStr += `}}`
	
	w.Write([]byte(jsonStr))
}