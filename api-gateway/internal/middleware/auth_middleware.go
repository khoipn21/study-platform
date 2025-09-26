package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	authpb "github.com/study-platform/auth-service/proto"
	coursepb "github.com/study-platform/course-service/proto"
	"github.com/study-platform/pkg/logger"
	"github.com/gorilla/mux"
	"google.golang.org/grpc"
)

type AuthMiddleware struct {
	authClient   authpb.AuthServiceClient
	courseClient coursepb.CourseServiceClient
	logger       logger.Logger
}

func NewAuthMiddleware(authConn *grpc.ClientConn, logger logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authpb.NewAuthServiceClient(authConn),
		logger:     logger,
	}
}

func NewAuthMiddlewareWithCourse(authConn *grpc.ClientConn, courseConn *grpc.ClientConn, logger logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authClient:   authpb.NewAuthServiceClient(authConn),
		courseClient: coursepb.NewCourseServiceClient(courseConn),
		logger:       logger,
	}
}

// RequireAuth is a middleware that requires authentication
func (am *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			am.sendError(w, http.StatusUnauthorized, "Authorization header is required")
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			am.sendError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		token := parts[1]
		if token == "" {
			am.sendError(w, http.StatusUnauthorized, "Token is required")
			return
		}

		// Validate token with auth service
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &authpb.ValidateTokenRequest{
			Token: token,
		}

		resp, err := am.authClient.ValidateToken(ctx, req)
		if err != nil {
			am.logger.Error(fmt.Errorf("token validation failed: %w", err))
			am.sendError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		if !resp.Valid {
			am.sendError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Add user information to request context
		if resp.User != nil {
			ctx = context.WithValue(r.Context(), "user_id", resp.User.Id)
			ctx = context.WithValue(ctx, "user_role", resp.User.Role)
			ctx = context.WithValue(ctx, "user_email", resp.User.Email)
			ctx = context.WithValue(ctx, "user_username", resp.User.Username)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// OptionalAuth is a middleware that optionally validates authentication
func (am *AuthMiddleware) OptionalAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			// No auth header, continue without authentication
			next.ServeHTTP(w, r)
			return
		}

		// CRITICAL FIX for BUG-006: Strict token validation for security
		// If Authorization header is present, the token MUST be valid

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			am.logger.Errorf("SECURITY: Invalid authorization header format from %s", r.RemoteAddr)
			am.sendError(w, http.StatusUnauthorized, "Invalid authorization header format")
			return
		}

		token := parts[1]
		if token == "" {
			am.logger.Errorf("SECURITY: Empty token provided from %s", r.RemoteAddr)
			am.sendError(w, http.StatusUnauthorized, "Empty token provided")
			return
		}

		// Basic token format validation (JWT should have 3 parts separated by dots)
		tokenParts := strings.Split(token, ".")
		if len(tokenParts) != 3 {
			am.logger.Errorf("SECURITY: Malformed JWT token from %s", r.RemoteAddr)
			am.sendError(w, http.StatusUnauthorized, "Malformed token")
			return
		}

		// Validate token with auth service
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		req := &authpb.ValidateTokenRequest{
			Token: token,
		}

		resp, err := am.authClient.ValidateToken(ctx, req)
		if err != nil {
			am.logger.Errorf("SECURITY: Token validation service error from %s: %v", r.RemoteAddr, err)
			am.sendError(w, http.StatusUnauthorized, "Token validation failed")
			return
		}

		if !resp.Valid {
			am.logger.Errorf("SECURITY: Invalid token provided from %s", r.RemoteAddr)
			am.sendError(w, http.StatusUnauthorized, "Invalid token")
			return
		}

		// Token is valid - add user information to request context
		if resp.User != nil {
			ctx = context.WithValue(r.Context(), "user_id", resp.User.Id)
			ctx = context.WithValue(ctx, "user_role", resp.User.Role)
			ctx = context.WithValue(ctx, "user_email", resp.User.Email)
			ctx = context.WithValue(ctx, "user_username", resp.User.Username)
			r = r.WithContext(ctx)
		}

		next.ServeHTTP(w, r)
	})
}

// RequireRole is a middleware that requires a specific role
func (am *AuthMiddleware) RequireRole(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userRole, ok := r.Context().Value("user_role").(string)
			if !ok {
				am.sendError(w, http.StatusUnauthorized, "User not authenticated")
				return
			}

			// Check if user has required role
			hasRole := false
			for _, role := range roles {
				if userRole == role {
					hasRole = true
					break
				}
			}

			if !hasRole {
				am.sendError(w, http.StatusForbidden, "Insufficient permissions")
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RequireInstructor is a middleware that requires instructor role
func (am *AuthMiddleware) RequireInstructor(next http.Handler) http.Handler {
	return am.RequireRole("instructor", "admin")(next)
}

// RequireAdmin is a middleware that requires admin role
func (am *AuthMiddleware) RequireAdmin(next http.Handler) http.Handler {
	return am.RequireRole("admin")(next)
}

// RequireCourseAccess is a middleware that checks if user has access to a specific course
func (am *AuthMiddleware) RequireCourseAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (must be authenticated first)
		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			am.sendError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		// Get course ID from URL path variables
		vars := mux.Vars(r)
		courseID := vars["courseId"]
		if courseID == "" {
			courseID = vars["course_id"] // Try alternative param name
		}
		if courseID == "" {
			am.sendError(w, http.StatusBadRequest, "Course ID is required")
			return
		}

		// Check if course client is available
		if am.courseClient == nil {
			am.logger.Errorf("Course client not initialized")
			am.sendError(w, http.StatusInternalServerError, "Course service unavailable")
			return
		}

		// Check enrollment
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		enrollmentReq := &coursepb.GetEnrollmentRequest{
			UserId:   userID,
			CourseId: courseID,
		}

		enrollment, err := am.courseClient.GetEnrollment(ctx, enrollmentReq)
		if err != nil {
			am.logger.Errorf("Failed to check enrollment: %v", err)
			am.sendError(w, http.StatusForbidden, "Access denied: Not enrolled in this course")
			return
		}

		// Check if enrollment is active
		if enrollment.Enrollment.Status != "active" {
			am.sendError(w, http.StatusForbidden, "Access denied: Enrollment is not active")
			return
		}

		// Add course access info to context
		ctx = context.WithValue(r.Context(), "course_id", courseID)
		ctx = context.WithValue(ctx, "enrollment_status", enrollment.Enrollment.Status)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// RequireLectureAccess is a middleware that checks if user has access to a specific lecture
func (am *AuthMiddleware) RequireLectureAccess(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get user ID from context (must be authenticated first)
		userID, ok := r.Context().Value("user_id").(string)
		if !ok {
			am.sendError(w, http.StatusUnauthorized, "User not authenticated")
			return
		}

		// Get lecture ID from URL path variables
		vars := mux.Vars(r)
		lectureID := vars["lectureId"]
		if lectureID == "" {
			lectureID = vars["lecture_id"] // Try alternative param name
		}
		if lectureID == "" {
			am.sendError(w, http.StatusBadRequest, "Lecture ID is required")
			return
		}

		// Check if course client is available
		if am.courseClient == nil {
			am.logger.Errorf("Course client not initialized")
			am.sendError(w, http.StatusInternalServerError, "Course service unavailable")
			return
		}

		// Get lecture details to find the course ID
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		lectureReq := &coursepb.GetLectureRequest{
			Id: lectureID,
		}

		lecture, err := am.courseClient.GetLecture(ctx, lectureReq)
		if err != nil {
			am.logger.Errorf("Failed to get lecture: %v", err)
			am.sendError(w, http.StatusNotFound, "Lecture not found")
			return
		}

		// Check if lecture is free (public access)
		if lecture.Lecture.IsFree {
			// Add lecture access info to context
			ctx = context.WithValue(r.Context(), "lecture_id", lectureID)
			ctx = context.WithValue(ctx, "course_id", lecture.Lecture.CourseId)
			ctx = context.WithValue(ctx, "is_free_lecture", true)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
			return
		}

		// Check enrollment for paid lectures
		enrollmentReq := &coursepb.GetEnrollmentRequest{
			UserId:   userID,
			CourseId: lecture.Lecture.CourseId,
		}

		enrollment, err := am.courseClient.GetEnrollment(ctx, enrollmentReq)
		if err != nil {
			am.logger.Errorf("Failed to check enrollment: %v", err)
			am.sendError(w, http.StatusForbidden, "Access denied: Not enrolled in this course")
			return
		}

		// Check if enrollment is active
		if enrollment.Enrollment.Status != "active" {
			am.sendError(w, http.StatusForbidden, "Access denied: Enrollment is not active")
			return
		}

		// Add lecture access info to context
		ctx = context.WithValue(r.Context(), "lecture_id", lectureID)
		ctx = context.WithValue(ctx, "course_id", lecture.Lecture.CourseId)
		ctx = context.WithValue(ctx, "is_free_lecture", false)
		ctx = context.WithValue(ctx, "enrollment_status", enrollment.Enrollment.Status)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// sendError sends an error response
func (am *AuthMiddleware) sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	// Simple JSON encoding without using json.NewEncoder to avoid import
	jsonStr := `{"success":false,"message":"` + message + `","error":"` + message + `"}`
	w.Write([]byte(jsonStr))
}