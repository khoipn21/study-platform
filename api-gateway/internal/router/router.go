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