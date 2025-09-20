package middleware

import (
	"context"
	"fmt"
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

// sendError sends an error response
func (am *AuthMiddleware) sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	// Simple JSON encoding without using json.NewEncoder to avoid import
	jsonStr := `{"success":false,"message":"` + message + `","error":"` + message + `"}`
	w.Write([]byte(jsonStr))
}