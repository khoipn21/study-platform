package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/study-platform/pkg/logger"
)

// Permission represents a specific permission
type Permission struct {
	Resource string
	Action   string
}

// AuthorizationMiddleware handles fine-grained authorization
type AuthorizationMiddleware struct {
	logger      logger.Logger
	permissions map[string][]Permission // endpoint pattern -> required permissions
}

// NewAuthorizationMiddleware creates a new authorization middleware
func NewAuthorizationMiddleware(logger logger.Logger) *AuthorizationMiddleware {
	am := &AuthorizationMiddleware{
		logger:      logger,
		permissions: make(map[string][]Permission),
	}

	// Define default permissions for endpoints
	am.setupDefaultPermissions()
	return am
}

// setupDefaultPermissions defines default permission requirements for endpoints
func (am *AuthorizationMiddleware) setupDefaultPermissions() {
	// Course management permissions
	am.permissions["POST:/api/v1/courses"] = []Permission{
		{Resource: "courses", Action: "create"},
	}
	am.permissions["PUT:/api/v1/courses/*"] = []Permission{
		{Resource: "courses", Action: "update"},
	}
	am.permissions["DELETE:/api/v1/courses/*"] = []Permission{
		{Resource: "courses", Action: "delete"},
	}

	// Lecture management permissions
	am.permissions["POST:/api/v1/lectures"] = []Permission{
		{Resource: "lectures", Action: "create"},
	}
	am.permissions["PUT:/api/v1/lectures/*"] = []Permission{
		{Resource: "lectures", Action: "update"},
	}
	am.permissions["DELETE:/api/v1/lectures/*"] = []Permission{
		{Resource: "lectures", Action: "delete"},
	}

	// User management permissions
	am.permissions["GET:/api/v1/users"] = []Permission{
		{Resource: "users", Action: "view"},
	}
	am.permissions["PUT:/api/v1/users/*"] = []Permission{
		{Resource: "users", Action: "update"},
	}
	am.permissions["DELETE:/api/v1/users/*"] = []Permission{
		{Resource: "users", Action: "delete"},
	}

	// Analytics permissions
	am.permissions["GET:/api/v1/analytics"] = []Permission{
		{Resource: "analytics", Action: "view"},
	}

	// Forum moderation permissions
	am.permissions["DELETE:/api/v1/forums/posts/*"] = []Permission{
		{Resource: "forums", Action: "moderate"},
	}
	am.permissions["PUT:/api/v1/forums/posts/*/moderate"] = []Permission{
		{Resource: "forums", Action: "moderate"},
	}
}

// AddPermissionRequirement adds a permission requirement for an endpoint
func (am *AuthorizationMiddleware) AddPermissionRequirement(endpoint string, permissions []Permission) {
	am.permissions[endpoint] = permissions
}

// Authorize middleware checks if the user has required permissions
func (am *AuthorizationMiddleware) Authorize(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user information from context (set by authentication middleware)
		userInfo := getUserInfoFromContext(r.Context())
		if userInfo == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check permissions for the requested endpoint
		endpointKey := fmt.Sprintf("%s:%s", r.Method, r.URL.Path)
		requiredPerms := am.getRequiredPermissions(endpointKey)

		if len(requiredPerms) > 0 && !am.hasPermissions(userInfo, requiredPerms) {
			am.logger.Warnf("User %s lacks required permissions for %s", userInfo.UserID, endpointKey)
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// getRequiredPermissions gets required permissions for an endpoint
func (am *AuthorizationMiddleware) getRequiredPermissions(endpoint string) []Permission {
	// First try exact match
	if perms, exists := am.permissions[endpoint]; exists {
		return perms
	}

	// Try pattern matching (for endpoints with IDs)
	for pattern, perms := range am.permissions {
		if am.matchesPattern(endpoint, pattern) {
			return perms
		}
	}

	return nil
}

// matchesPattern checks if an endpoint matches a permission pattern
func (am *AuthorizationMiddleware) matchesPattern(endpoint, pattern string) bool {
	// Simple wildcard matching for now
	// This can be enhanced with more sophisticated pattern matching
	if strings.Contains(pattern, "*") {
		prefix := strings.Split(pattern, "*")[0]
		return strings.HasPrefix(endpoint, prefix)
	}
	return endpoint == pattern
}

// hasPermissions checks if user has all required permissions
func (am *AuthorizationMiddleware) hasPermissions(userInfo *UserInfo, requiredPerms []Permission) bool {
	// Admin role has all permissions
	if am.hasRole(userInfo, "admin") {
		return true
	}

	// Check each required permission
	for _, requiredPerm := range requiredPerms {
		if !am.hasPermission(userInfo, requiredPerm) {
			return false
		}
	}

	return true
}

// hasPermission checks if user has a specific permission
func (am *AuthorizationMiddleware) hasPermission(userInfo *UserInfo, permission Permission) bool {
	// Check role-based permissions
	switch userInfo.Role {
	case "admin":
		return true
	case "instructor":
		return am.hasInstructorPermission(permission)
	case "student":
		return am.hasStudentPermission(permission)
	}

	return false
}

// hasInstructorPermission checks instructor-specific permissions
func (am *AuthorizationMiddleware) hasInstructorPermission(permission Permission) bool {
	instructorPermissions := map[string][]string{
		"courses":  {"create", "update", "delete", "view"},
		"lectures": {"create", "update", "delete", "view"},
		"students": {"view"},
		"forums":   {"moderate"},
	}

	if actions, exists := instructorPermissions[permission.Resource]; exists {
		for _, action := range actions {
			if action == permission.Action {
				return true
			}
		}
	}

	return false
}

// hasStudentPermission checks student-specific permissions
func (am *AuthorizationMiddleware) hasStudentPermission(permission Permission) bool {
	studentPermissions := map[string][]string{
		"courses":     {"view"},
		"lectures":    {"view"},
		"enrollment":  {"create", "view"},
		"progress":    {"create", "update", "view"},
		"forums":      {"create", "update", "view"},
	}

	if actions, exists := studentPermissions[permission.Resource]; exists {
		for _, action := range actions {
			if action == permission.Action {
				return true
			}
		}
	}

	return false
}

// hasRole checks if user has a specific role
func (am *AuthorizationMiddleware) hasRole(userInfo *UserInfo, role string) bool {
	return userInfo.Role == role
}

// ResourceOwnershipMiddleware checks if user owns the requested resource
type ResourceOwnershipMiddleware struct {
	logger logger.Logger
}

// NewResourceOwnershipMiddleware creates a new resource ownership middleware
func NewResourceOwnershipMiddleware(logger logger.Logger) *ResourceOwnershipMiddleware {
	return &ResourceOwnershipMiddleware{
		logger: logger,
	}
}

// CheckOwnership middleware for resource ownership
func (rom *ResourceOwnershipMiddleware) CheckOwnership(resourceType string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userInfo := getUserInfoFromContext(r.Context())
			if userInfo == nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Admin can access all resources
			if userInfo.Role == "admin" {
				next.ServeHTTP(w, r)
				return
			}

			// Extract resource ID from URL
			resourceID := extractResourceID(r.URL.Path)
			if resourceID == "" {
				rom.logger.Warnf("Could not extract resource ID from path: %s", r.URL.Path)
				http.Error(w, "Bad Request", http.StatusBadRequest)
				return
			}

			// Check ownership based on resource type
			owns, err := rom.checkResourceOwnership(userInfo.UserID, resourceType, resourceID)
			if err != nil {
				rom.logger.Errorf("Error checking resource ownership: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if !owns {
				rom.logger.Warnf("User %s does not own %s %s", userInfo.UserID, resourceType, resourceID)
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// checkResourceOwnership checks if user owns a specific resource
func (rom *ResourceOwnershipMiddleware) checkResourceOwnership(userID, resourceType, resourceID string) (bool, error) {
	// This would typically query the database to check ownership
	// For now, we'll implement a basic check
	// In a real implementation, you'd inject a repository or service here

	switch resourceType {
	case "course":
		// Check if user is the creator of the course
		return rom.checkCourseOwnership(userID, resourceID)
	case "lecture":
		// Check if user owns the course that contains this lecture
		return rom.checkLectureOwnership(userID, resourceID)
	default:
		return false, fmt.Errorf("unknown resource type: %s", resourceType)
	}
}

// checkCourseOwnership checks if user owns a course
func (rom *ResourceOwnershipMiddleware) checkCourseOwnership(userID, courseID string) (bool, error) {
	// TODO: Implement database query to check course ownership
	// This is a placeholder implementation
	return true, nil
}

// checkLectureOwnership checks if user owns a lecture (through course ownership)
func (rom *ResourceOwnershipMiddleware) checkLectureOwnership(userID, lectureID string) (bool, error) {
	// TODO: Implement database query to check lecture ownership through course
	// This is a placeholder implementation
	return true, nil
}

// extractResourceID extracts resource ID from URL path
func extractResourceID(path string) string {
	// Simple implementation - extract last segment after known patterns
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// UserInfo represents user information from JWT token
type UserInfo struct {
	UserID string
	Role   string
	Email  string
}

// getUserInfoFromContext extracts user info from request context
func getUserInfoFromContext(ctx context.Context) *UserInfo {
	if userInfo, ok := ctx.Value("user").(*UserInfo); ok {
		return userInfo
	}
	return nil
}