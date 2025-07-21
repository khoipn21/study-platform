package handler

import (
	"context"
	"encoding/json"
	"net/http"
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