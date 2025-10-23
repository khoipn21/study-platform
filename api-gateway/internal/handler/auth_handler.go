package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	authpb "github.com/study-platform/auth-service/proto"
	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthHandler struct {
	authClient authpb.AuthServiceClient
	logger     logger.Logger
}

func NewAuthHandler(authConn *grpc.ClientConn, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		authClient: authpb.NewAuthServiceClient(authConn),
		logger:     logger,
	}
}

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type ValidateTokenRequest struct {
	Token string `json:"token"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account with username, email, and password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body RegisterRequest true "Registration details"
// @Success      201 {object} APIResponse "User registered successfully"
// @Failure      400 {object} APIResponse "Invalid request"
// @Failure      409 {object} APIResponse "User already exists"
// @Router       /auth/register [post]
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Username == "" || req.Email == "" || req.Password == "" {
		h.sendError(w, http.StatusBadRequest, "Username, email, and password are required")
		return
	}

	// Set default role if not provided
	if req.Role == "" {
		req.Role = "student"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call auth service
	grpcReq := &authpb.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     req.Role,
	}

	resp, err := h.authClient.Register(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Registration failed")
		return
	}

	// Format response
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"id":       resp.User.Id,
			"username": resp.User.Username,
			"email":    resp.User.Email,
			"role":     resp.User.Role,
		},
		"token": resp.Token,
	}

	h.sendSuccess(w, resp.Message, data)
}

// Login godoc
// @Summary      User login
// @Description  Authenticate user with email and password, returns JWT token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body LoginRequest true "Login credentials"
// @Success      200 {object} APIResponse "Login successful"
// @Failure      400 {object} APIResponse "Invalid request"
// @Failure      401 {object} APIResponse "Invalid credentials"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		h.sendError(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call auth service
	grpcReq := &authpb.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	resp, err := h.authClient.Login(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Login failed")
		return
	}

	// Format response
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"id":       resp.User.Id,
			"username": resp.User.Username,
			"email":    resp.User.Email,
			"role":     resp.User.Role,
		},
		"token": resp.Token,
	}

	h.sendSuccess(w, resp.Message, data)
}

// ValidateToken godoc
// @Summary      Validate JWT token
// @Description  Validate a JWT token and return user information
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body ValidateTokenRequest true "Token to validate"
// @Success      200 {object} APIResponse "Token is valid"
// @Failure      400 {object} APIResponse "Invalid request"
// @Failure      401 {object} APIResponse "Invalid or expired token"
// @Router       /auth/validate [post]
func (h *AuthHandler) ValidateToken(w http.ResponseWriter, r *http.Request) {
	var req ValidateTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Token == "" {
		h.sendError(w, http.StatusBadRequest, "Token is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call auth service
	grpcReq := &authpb.ValidateTokenRequest{
		Token: req.Token,
	}

	resp, err := h.authClient.ValidateToken(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Token validation failed")
		return
	}

	// Format response
	data := map[string]interface{}{
		"valid": resp.Valid,
		"user":  nil,
	}

	if resp.User != nil {
		data["user"] = map[string]interface{}{
			"id":       resp.User.Id,
			"username": resp.User.Username,
			"email":    resp.User.Email,
			"role":     resp.User.Role,
		}
	}

	h.sendSuccess(w, resp.Message, data)
}

// GetProfile godoc
// @Summary      Get user profile
// @Description  Get the authenticated user's profile information
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Success      200 {object} APIResponse "User profile"
// @Failure      401 {object} APIResponse "Unauthorized"
// @Failure      404 {object} APIResponse "User not found"
// @Security     BearerAuth
// @Router       /auth/profile [get]
func (h *AuthHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call auth service to get user roles
	grpcReq := &authpb.GetUserRolesRequest{
		UserId: userID,
	}

	resp, err := h.authClient.GetUserRoles(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to get user profile")
		return
	}

	// Format response
	data := map[string]interface{}{
		"id":    userID,
		"roles": resp.Roles,
	}

	h.sendSuccess(w, "Profile retrieved successfully", data)
}

func (h *AuthHandler) GetOAuthURL(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]

	if provider == "" {
		h.sendError(w, http.StatusBadRequest, "Provider is required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call auth service
	grpcReq := &authpb.GetOAuthURLRequest{
		Provider: provider,
	}

	resp, err := h.authClient.GetOAuthURL(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to get OAuth URL")
		return
	}

	// Format response
	data := map[string]interface{}{
		"url": resp.Url,
	}

	h.sendSuccess(w, "OAuth URL generated successfully", data)
}

func (h *AuthHandler) OAuthCallback(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	provider := vars["provider"]

	code := r.URL.Query().Get("code")
	state := r.URL.Query().Get("state")

	if provider == "" || code == "" || state == "" {
		h.sendError(w, http.StatusBadRequest, "Provider, code, and state are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call auth service
	grpcReq := &authpb.OAuthCallbackRequest{
		Provider: provider,
		Code:     code,
		State:    state,
	}

	resp, err := h.authClient.OAuthCallback(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "OAuth callback failed")
		return
	}

	// Redirect to frontend OAuth callback handler
	// The frontend will handle role-based routing
	frontendURL := getEnv("FRONTEND_URL", "http://localhost:3000")
	redirectURL := frontendURL + "/auth/oauth/callback?token=" + resp.Token + "&user=" + resp.User.Id

	// Redirect to frontend with authentication data
	http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
}

// Helper methods
func (h *AuthHandler) sendSuccess(w http.ResponseWriter, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) sendError(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := APIResponse{
		Success: false,
		Message: message,
		Error:   message,
	}

	json.NewEncoder(w).Encode(response)
}

func (h *AuthHandler) handleGRPCError(w http.ResponseWriter, err error, defaultMessage string) {
	if grpcErr, ok := status.FromError(err); ok {
		switch grpcErr.Code() {
		case codes.NotFound:
			h.sendError(w, http.StatusNotFound, grpcErr.Message())
		case codes.InvalidArgument:
			h.sendError(w, http.StatusBadRequest, grpcErr.Message())
		case codes.Unauthenticated:
			h.sendError(w, http.StatusUnauthorized, grpcErr.Message())
		case codes.PermissionDenied:
			h.sendError(w, http.StatusForbidden, grpcErr.Message())
		case codes.AlreadyExists:
			h.sendError(w, http.StatusConflict, grpcErr.Message())
		default:
			h.sendError(w, http.StatusInternalServerError, defaultMessage)
		}
	} else {
		h.sendError(w, http.StatusInternalServerError, defaultMessage)
	}

	h.logger.Errorf("gRPC error: %v", err)
}

// getEnv retrieves an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ListUsers godoc
// @Summary      List users for mentions
// @Description  Get list of users for autocomplete/mentions
// @Tags         Users
// @Accept       json
// @Produce      json
// @Param        q query string false "Search query (username)"
// @Param        limit query int false "Limit results (default 20, max 100)"
// @Success      200 {object} APIResponse "Users retrieved successfully"
// @Failure      500 {object} APIResponse "Failed to fetch users"
// @Router       /users [get]
func (h *AuthHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	limitStr := r.URL.Query().Get("limit")

	limit := int32(20)
	if limitStr != "" {
		if parsedLimit, err := strconv.Atoi(limitStr); err == nil && parsedLimit > 0 && parsedLimit <= 100 {
			limit = int32(parsedLimit)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call auth service
	grpcReq := &authpb.ListUsersRequest{
		Query: query,
		Limit: limit,
	}

	resp, err := h.authClient.ListUsers(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to fetch users")
		return
	}

	// Format response
	users := make([]map[string]interface{}, 0, len(resp.Users))
	for _, user := range resp.Users {
		users = append(users, map[string]interface{}{
			"id":       user.Id,
			"username": user.Username,
			"role":     user.Role,
			"avatar":   user.Avatar,
		})
	}

	data := map[string]interface{}{
		"users": users,
	}

	h.sendSuccess(w, "Users retrieved successfully", data)
}

// GetCurrentUser godoc
// @Summary      Get current user details
// @Description  Get the authenticated user's full profile information
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Success      200 {object} APIResponse "User retrieved successfully"
// @Failure      401 {object} APIResponse "Unauthorized"
// @Failure      404 {object} APIResponse "User not found"
// @Security     BearerAuth
// @Router       /auth/me [get]
func (h *AuthHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call auth service
	grpcReq := &authpb.GetCurrentUserRequest{
		UserId: userID,
	}

	resp, err := h.authClient.GetCurrentUser(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to get user")
		return
	}

	// Format response
	data := map[string]interface{}{
		"id":         resp.User.Id,
		"username":   resp.User.Username,
		"email":      resp.User.Email,
		"role":       resp.User.Role,
		"created_at": resp.User.CreatedAt,
		"updated_at": resp.User.UpdatedAt,
		"avatar_url": resp.User.AvatarUrl,
	}

	h.sendSuccess(w, resp.Message, data)
}

type UpdateProfileRequest struct {
	Username  string `json:"username,omitempty"`
	Email     string `json:"email,omitempty"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// UpdateProfile godoc
// @Summary      Update user profile
// @Description  Update the authenticated user's profile information
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body UpdateProfileRequest true "Profile update details"
// @Success      200 {object} APIResponse "Profile updated successfully"
// @Failure      400 {object} APIResponse "Invalid request"
// @Failure      401 {object} APIResponse "Unauthorized"
// @Failure      409 {object} APIResponse "Username or email already exists"
// @Security     BearerAuth
// @Router       /auth/profile [put]
func (h *AuthHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// At least one field must be provided
	if req.Username == "" && req.Email == "" && req.AvatarURL == "" {
		h.sendError(w, http.StatusBadRequest, "At least one field must be provided")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call auth service
	grpcReq := &authpb.UpdateProfileRequest{
		UserId:    userID,
		Username:  req.Username,
		Email:     req.Email,
		AvatarUrl: req.AvatarURL,
	}

	resp, err := h.authClient.UpdateProfile(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to update profile")
		return
	}

	// Format response
	data := map[string]interface{}{
		"id":         resp.User.Id,
		"username":   resp.User.Username,
		"email":      resp.User.Email,
		"role":       resp.User.Role,
		"created_at": resp.User.CreatedAt,
		"updated_at": resp.User.UpdatedAt,
	}

	h.sendSuccess(w, resp.Message, data)
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ChangePassword godoc
// @Summary      Change user password
// @Description  Change the authenticated user's password
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body ChangePasswordRequest true "Password change details"
// @Success      200 {object} APIResponse "Password changed successfully"
// @Failure      400 {object} APIResponse "Invalid request"
// @Failure      401 {object} APIResponse "Unauthorized"
// @Failure      403 {object} APIResponse "Current password is incorrect"
// @Security     BearerAuth
// @Router       /auth/password [put]
func (h *AuthHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	var req ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.CurrentPassword == "" || req.NewPassword == "" {
		h.sendError(w, http.StatusBadRequest, "Current password and new password are required")
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Call auth service
	grpcReq := &authpb.ChangePasswordRequest{
		UserId:          userID,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	}

	resp, err := h.authClient.ChangePassword(ctx, grpcReq)
	if err != nil {
		h.handleGRPCError(w, err, "Failed to change password")
		return
	}

	// Format response
	data := map[string]interface{}{
		"success": resp.Success,
	}

	h.sendSuccess(w, resp.Message, data)
}

// UploadAvatar godoc
// @Summary      Upload user avatar
// @Description  Upload a new avatar image for the authenticated user. Replaces existing avatar at {userId}/avatar path.
// @Tags         Authentication
// @Accept       multipart/form-data
// @Produce      json
// @Param        file formData file true "Avatar image file (max 5MB, jpg/png/webp)"
// @Success      200 {object} APIResponse "Avatar uploaded successfully"
// @Failure      400 {object} APIResponse "Invalid file"
// @Failure      401 {object} APIResponse "Unauthorized"
// @Security     BearerAuth
// @Router       /auth/avatar [post]
func (h *AuthHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context (set by auth middleware)
	userID, ok := r.Context().Value("user_id").(string)
	if !ok {
		h.sendError(w, http.StatusUnauthorized, "User not authenticated")
		return
	}

	// Parse multipart form (max 10MB)
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		h.sendError(w, http.StatusBadRequest, "Failed to parse form data")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.sendError(w, http.StatusBadRequest, "File is required")
		return
	}
	defer file.Close()

	// Validate file size (5MB max)
	if header.Size > 5<<20 {
		h.sendError(w, http.StatusBadRequest, "File size must be less than 5MB")
		return
	}

	// Validate file type
	contentType := header.Header.Get("Content-Type")
	if contentType == "" {
		// Detect from filename extension
		ext := header.Filename[len(header.Filename)-4:]
		switch ext {
		case ".jpg", "jpeg":
			contentType = "image/jpeg"
		case ".png":
			contentType = "image/png"
		case "webp":
			contentType = "image/webp"
		default:
			h.sendError(w, http.StatusBadRequest, "Invalid file type. Only jpg, png, webp allowed")
			return
		}
	}

	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		h.sendError(w, http.StatusBadRequest, "Invalid file type. Only jpg, png, webp allowed")
		return
	}

	// Create multipart form for bucket service
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Add file
	part, err := writer.CreateFormFile("file", header.Filename)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to create form file")
		return
	}

	if _, err := io.Copy(part, file); err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to copy file")
		return
	}

	// Add metadata
	writer.WriteField("bucket", "avatars")
	writer.WriteField("is_public", "true")
	writer.WriteField("key", userID+"/avatar") // Custom key for avatar path
	writer.Close()

	// Forward to bucket service
	bucketServiceURL := getEnv("BUCKET_SERVICE_URL", "bucket-service:8088")
	uploadURL := fmt.Sprintf("http://%s/api/files/upload", bucketServiceURL)

	req, err := http.NewRequest("POST", uploadURL, body)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to create request")
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	req = req.WithContext(ctx)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		h.logger.Errorf("Failed to upload to bucket service: %v", err)
		h.sendError(w, http.StatusInternalServerError, "Failed to upload avatar")
		return
	}
	defer resp.Body.Close()

	// Read response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		h.sendError(w, http.StatusInternalServerError, "Failed to read response")
		return
	}

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		h.logger.Errorf("Bucket service error: %s", string(respBody))
		h.sendError(w, resp.StatusCode, "Failed to upload avatar")
		return
	}

	// Parse bucket service response
	var bucketResp struct {
		FileID      string `json:"file_id"`
		Filename    string `json:"filename"`
		URL         string `json:"url"`
		ContentType string `json:"content_type"`
	}

	if err := json.Unmarshal(respBody, &bucketResp); err != nil {
		h.logger.Errorf("Failed to parse bucket response: %v", err)
		h.sendError(w, http.StatusInternalServerError, "Failed to parse response")
		return
	}

	// Update user profile with new avatar URL
	ctx2, cancel2 := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel2()

	grpcReq := &authpb.UpdateProfileRequest{
		UserId:    userID,
		AvatarUrl: bucketResp.URL,
	}

	_, err = h.authClient.UpdateProfile(ctx2, grpcReq)
	if err != nil {
		h.logger.Errorf("Failed to update profile with avatar: %v", err)
		// Avatar uploaded but profile update failed - not critical
	}

	// Return success with avatar URL
	data := map[string]interface{}{
		"avatar_url": bucketResp.URL,
		"message":    "Avatar uploaded successfully",
	}

	h.sendSuccess(w, "Avatar uploaded successfully", data)
}
