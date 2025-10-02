package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/study-platform/auth-service/internal/model"
	"github.com/study-platform/auth-service/internal/repository"
	"github.com/study-platform/pkg/logger"
)

type AuthService struct {
	userRepo  *repository.UserRepository
	jwtSecret string
	logger    logger.Logger
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, logger logger.Logger) *AuthService {
	return &AuthService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		logger:    logger,
	}
}

func (s *AuthService) Register(req *model.CreateUserRequest) (*model.User, string, error) {
	if !req.Role.IsValid() {
		return nil, "", fmt.Errorf("invalid role: %s", req.Role)
	}

	existingUser, _ := s.userRepo.GetByEmail(req.Email)
	if existingUser != nil {
		return nil, "", fmt.Errorf("user with email %s already exists", req.Email)
	}

	existingUser, _ = s.userRepo.GetByUsername(req.Username)
	if existingUser != nil {
		return nil, "", fmt.Errorf("user with username %s already exists", req.Username)
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", fmt.Errorf("failed to hash password: %w", err)
	}

	hashedPasswordStr := string(hashedPassword)
	user := &model.User{
		ID:              uuid.New(),
		Username:        req.Username,
		Email:           req.Email,
		PasswordHash:    &hashedPasswordStr,
		Role:            req.Role,
		Provider:        model.ProviderLocal,
		ProviderID:      nil,
		AvatarURL:       nil,
		IsEmailVerified: false,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err = s.userRepo.Create(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Infof("User registered successfully: %s", user.Email)
	return user, token, nil
}

func (s *AuthService) Login(req *model.LoginRequest) (*model.User, string, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	if user.PasswordHash == nil {
		return nil, "", fmt.Errorf("user registered via OAuth, use OAuth login")
	}

	err = bcrypt.CompareHashAndPassword([]byte(*user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, "", fmt.Errorf("invalid credentials")
	}

	token, err := s.generateToken(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	s.logger.Infof("User logged in successfully: %s", user.Email)
	return user, token, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*model.User, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	userID, err := uuid.Parse(claims.UserID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in token")
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("user not found")
	}

	return user, nil
}

func (s *AuthService) AssignRole(userID uuid.UUID, role model.UserRole) error {
	if !role.IsValid() {
		return fmt.Errorf("invalid role: %s", role)
	}

	err := s.userRepo.UpdateRole(userID, role)
	if err != nil {
		return fmt.Errorf("failed to assign role: %w", err)
	}

	s.logger.Infof("Role assigned successfully: %s to user %s", role, userID)
	return nil
}

func (s *AuthService) GetUserByID(userID uuid.UUID) (*model.User, error) {
	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (s *AuthService) SearchUsers(query string, limit, offset int) ([]*model.User, int, error) {
	if limit <= 0 {
		limit = 10 // Default limit
	}
	if limit > 100 {
		limit = 100 // Max limit
	}
	if offset < 0 {
		offset = 0
	}

	users, total, err := s.userRepo.SearchUsers(query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search users: %w", err)
	}

	return users, total, nil
}

func (s *AuthService) generateToken(user *model.User) (string, error) {
	claims := &Claims{
		UserID: user.ID.String(),
		Role:   user.Role.String(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Subject:   user.ID.String(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}