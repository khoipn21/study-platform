package handler

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/study-platform/auth-service/internal/model"
	"github.com/study-platform/auth-service/internal/service"
	"github.com/study-platform/pkg/logger"
	pb "github.com/study-platform/auth-service/proto"
)

type AuthHandler struct {
	pb.UnimplementedAuthServiceServer
	authService  *service.AuthService
	oauthService *service.OAuthService
	logger       logger.Logger
}

func NewAuthHandler(authService *service.AuthService, oauthService *service.OAuthService, logger logger.Logger) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		oauthService: oauthService,
		logger:       logger,
	}
}

func (h *AuthHandler) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	if req.Username == "" || req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "username, email, and password are required")
	}

	role := model.UserRole(req.Role)
	if req.Role == "" {
		role = model.RoleStudent
	}

	createReq := &model.CreateUserRequest{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
		Role:     role,
	}

	user, token, err := h.authService.Register(createReq)
	if err != nil {
		h.logger.Error(fmt.Errorf("register error: %w", err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RegisterResponse{
		User:    h.userToProto(user),
		Token:   token,
		Message: "User registered successfully",
	}, nil
}

func (h *AuthHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "email and password are required")
	}

	loginReq := &model.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}

	user, token, err := h.authService.Login(loginReq)
	if err != nil {
		h.logger.Error(fmt.Errorf("login error: %w", err))
		return nil, status.Error(codes.Unauthenticated, "invalid credentials")
	}

	return &pb.LoginResponse{
		User:    h.userToProto(user),
		Token:   token,
		Message: "Login successful",
	}, nil
}

func (h *AuthHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.Token == "" {
		return nil, status.Error(codes.InvalidArgument, "token is required")
	}

	user, err := h.authService.ValidateToken(req.Token)
	if err != nil {
		h.logger.Error(fmt.Errorf("validate token error: %w", err))
		return &pb.ValidateTokenResponse{
			Valid:   false,
			Message: "Invalid token",
		}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:   true,
		User:    h.userToProto(user),
		Message: "Token is valid",
	}, nil
}

func (h *AuthHandler) GetUserRoles(ctx context.Context, req *pb.GetUserRolesRequest) (*pb.GetUserRolesResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		h.logger.Error(fmt.Errorf("get user roles error: %w", err))
		return nil, status.Error(codes.NotFound, "user not found")
	}

	return &pb.GetUserRolesResponse{
		Roles: []string{user.Role.String()},
	}, nil
}

func (h *AuthHandler) AssignRole(ctx context.Context, req *pb.AssignRoleRequest) (*pb.AssignRoleResponse, error) {
	if req.UserId == "" || req.Role == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and role are required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	role := model.UserRole(req.Role)
	if !role.IsValid() {
		return nil, status.Error(codes.InvalidArgument, "invalid role")
	}

	err = h.authService.AssignRole(userID, role)
	if err != nil {
		h.logger.Error(fmt.Errorf("assign role error: %w", err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AssignRoleResponse{
		Success: true,
		Message: "Role assigned successfully",
	}, nil
}

func (h *AuthHandler) RemoveRole(ctx context.Context, req *pb.RemoveRoleRequest) (*pb.RemoveRoleResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	err = h.authService.AssignRole(userID, model.RoleStudent)
	if err != nil {
		h.logger.Error(fmt.Errorf("remove role error: %w", err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.RemoveRoleResponse{
		Success: true,
		Message: "Role removed successfully",
	}, nil
}

func (h *AuthHandler) GetOAuthURL(ctx context.Context, req *pb.GetOAuthURLRequest) (*pb.GetOAuthURLResponse, error) {
	if req.Provider == "" {
		return nil, status.Error(codes.InvalidArgument, "provider is required")
	}

	provider := model.OAuthProvider(req.Provider)
	if !provider.IsValid() {
		return nil, status.Error(codes.InvalidArgument, "invalid provider")
	}

	state := req.State
	if state == "" {
		state = "default-state"
	}

	url, err := h.oauthService.GetAuthURL(provider, state)
	if err != nil {
		h.logger.Error(fmt.Errorf("get oauth url error: %w", err))
		return nil, status.Error(codes.Internal, "failed to get oauth url")
	}

	return &pb.GetOAuthURLResponse{
		Url: url,
	}, nil
}

func (h *AuthHandler) OAuthCallback(ctx context.Context, req *pb.OAuthCallbackRequest) (*pb.OAuthCallbackResponse, error) {
	if req.Provider == "" || req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "provider and code are required")
	}

	provider := model.OAuthProvider(req.Provider)
	if !provider.IsValid() {
		return nil, status.Error(codes.InvalidArgument, "invalid provider")
	}

	result, err := h.oauthService.HandleOAuthCallback(ctx, provider, req.Code, req.State)
	if err != nil {
		h.logger.Error(fmt.Errorf("oauth callback error: %w", err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	linkedAccounts := make([]string, len(result.LinkedAccounts))
	for i, account := range result.LinkedAccounts {
		linkedAccounts[i] = account.String()
	}

	return &pb.OAuthCallbackResponse{
		User:           h.userToProto(&result.User),
		Token:          result.Token,
		IsNewUser:      result.IsNewUser,
		LinkedAccounts: linkedAccounts,
		Message:        "OAuth login successful",
	}, nil
}

func (h *AuthHandler) LinkOAuthAccount(ctx context.Context, req *pb.LinkOAuthAccountRequest) (*pb.LinkOAuthAccountResponse, error) {
	if req.UserId == "" || req.Provider == "" || req.Code == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id, provider, and code are required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	provider := model.OAuthProvider(req.Provider)
	if !provider.IsValid() {
		return nil, status.Error(codes.InvalidArgument, "invalid provider")
	}

	err = h.oauthService.LinkAccount(ctx, userID, provider, req.Code)
	if err != nil {
		h.logger.Error(fmt.Errorf("link oauth account error: %w", err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.LinkOAuthAccountResponse{
		Success: true,
		Message: "OAuth account linked successfully",
	}, nil
}

func (h *AuthHandler) UnlinkOAuthAccount(ctx context.Context, req *pb.UnlinkOAuthAccountRequest) (*pb.UnlinkOAuthAccountResponse, error) {
	if req.UserId == "" || req.Provider == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id and provider are required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	provider := model.OAuthProvider(req.Provider)
	if !provider.IsValid() {
		return nil, status.Error(codes.InvalidArgument, "invalid provider")
	}

	err = h.oauthService.UnlinkAccount(userID, provider)
	if err != nil {
		h.logger.Error(fmt.Errorf("unlink oauth account error: %w", err))
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UnlinkOAuthAccountResponse{
		Success: true,
		Message: "OAuth account unlinked successfully",
	}, nil
}

func (h *AuthHandler) GetLinkedAccounts(ctx context.Context, req *pb.GetLinkedAccountsRequest) (*pb.GetLinkedAccountsResponse, error) {
	if req.UserId == "" {
		return nil, status.Error(codes.InvalidArgument, "user_id is required")
	}

	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid user_id format")
	}

	providers, err := h.oauthService.GetLinkedAccounts(userID)
	if err != nil {
		h.logger.Error(fmt.Errorf("get linked accounts error: %w", err))
		return nil, status.Error(codes.Internal, "failed to get linked accounts")
	}

	providerStrings := make([]string, len(providers))
	for i, provider := range providers {
		providerStrings[i] = provider.String()
	}

	return &pb.GetLinkedAccountsResponse{
		Providers: providerStrings,
	}, nil
}

func (h *AuthHandler) userToProto(user *model.User) *pb.User {
	return &pb.User{
		Id:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role.String(),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: timestamppb.New(user.UpdatedAt),
	}
}
func (h *AuthHandler) ListUsers(ctx context.Context, req *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	// Default limit if not specified
	limit := req.Limit
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	users, err := h.authService.ListUsers(req.Query, limit)
	if err != nil {
		h.logger.Error(fmt.Errorf("list users error: %w", err))
		return nil, status.Error(codes.Internal, "Failed to fetch users")
	}

	// Convert to protobuf UserInfo
	pbUsers := make([]*pb.UserInfo, 0, len(users))
	for _, user := range users {
		avatar := ""
		if user.AvatarURL != nil {
			avatar = *user.AvatarURL
		}

		pbUsers = append(pbUsers, &pb.UserInfo{
			Id:       user.ID.String(),
			Username: user.Username,
			Role:     string(user.Role),
			Avatar:   avatar,
		})
	}

	return &pb.ListUsersResponse{
		Users: pbUsers,
	}, nil
}
