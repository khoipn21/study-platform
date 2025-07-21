package middleware

import (
	"context"
	"net/http"
	"strings"
	"time"

	authpb "github.com/study-platform/auth-service/proto"
	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc"
)

type AuthMiddleware struct {
	authClient authpb.AuthServiceClient
	logger     logger.Logger
}

func NewAuthMiddleware(authConn *grpc.ClientConn, logger logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authpb.NewAuthServiceClient(authConn),
		logger:     logger,
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
			am.logger.Errorf("Token validation failed: %v", err)
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

		// Extract token from "Bearer <token>" format
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			// Invalid format, continue without authentication
			next.ServeHTTP(w, r)
			return
		}

		token := parts[1]
		if token == "" {
			// No token, continue without authentication
			next.ServeHTTP(w, r)
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
			am.logger.Errorf("Token validation failed: %v", err)
			// Continue without authentication on error
			next.ServeHTTP(w, r)
			return
		}

		if resp.Valid && resp.User != nil {
			// Add user information to request context
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

// sendError sends an error response
func (am *AuthMiddleware) sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	// Simple JSON encoding without using json.NewEncoder to avoid import
	jsonStr := `{"success":false,"message":"` + message + `","error":"` + message + `"}`
	w.Write([]byte(jsonStr))
}