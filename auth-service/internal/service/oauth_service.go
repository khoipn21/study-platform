package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"

	"github.com/study-platform/auth-service/internal/model"
	"github.com/study-platform/auth-service/internal/repository"
	"github.com/study-platform/pkg/logger"
)

type OAuthService struct {
	userRepo        *repository.UserRepository
	oauthRepo       *repository.OAuthRepository
	authService     *AuthService
	logger          *logger.Logger
	oauthConfigs    map[model.OAuthProvider]*oauth2.Config
	defaultRedirect string
}

func NewOAuthService(
	userRepo *repository.UserRepository,
	oauthRepo *repository.OAuthRepository,
	authService *AuthService,
	logger *logger.Logger,
	configs map[model.OAuthProvider]model.OAuthConfig,
	defaultRedirect string,
) *OAuthService {
	oauthConfigs := make(map[model.OAuthProvider]*oauth2.Config)
	
	for provider, config := range configs {
		var endpoint oauth2.Endpoint
		
		switch provider {
		case model.ProviderGoogle:
			endpoint = google.Endpoint
		case model.ProviderGitHub:
			endpoint = github.Endpoint
		case model.ProviderFacebook:
			endpoint = facebook.Endpoint
		}
		
		oauthConfigs[provider] = &oauth2.Config{
			ClientID:     config.ClientID,
			ClientSecret: config.ClientSecret,
			RedirectURL:  config.RedirectURL,
			Scopes:       config.Scopes,
			Endpoint:     endpoint,
		}
	}
	
	return &OAuthService{
		userRepo:        userRepo,
		oauthRepo:       oauthRepo,
		authService:     authService,
		logger:          logger,
		oauthConfigs:    oauthConfigs,
		defaultRedirect: defaultRedirect,
	}
}

func (s *OAuthService) GetAuthURL(provider model.OAuthProvider, state string) (string, error) {
	config, exists := s.oauthConfigs[provider]
	if !exists {
		return "", fmt.Errorf("oauth provider %s not configured", provider)
	}
	
	return config.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

func (s *OAuthService) HandleOAuthCallback(ctx context.Context, provider model.OAuthProvider, code, state string) (*model.OAuthLoginResponse, error) {
	config, exists := s.oauthConfigs[provider]
	if !exists {
		return nil, fmt.Errorf("oauth provider %s not configured", provider)
	}
	
	// Exchange code for token
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}
	
	// Get user info from provider
	userInfo, err := s.getUserInfoFromProvider(ctx, provider, token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	
	// Check if user exists by OAuth provider ID
	existingUser, err := s.userRepo.GetByProviderID(provider, userInfo.ID)
	if err != nil && !strings.Contains(err.Error(), "user not found") {
		return nil, fmt.Errorf("failed to check existing user: %w", err)
	}
	
	var user *model.User
	var isNewUser bool
	
	if existingUser != nil {
		// User exists, update their info
		user = existingUser
		user.AvatarURL = &userInfo.AvatarURL
		user.IsEmailVerified = userInfo.EmailVerified
		isNewUser = false
	} else {
		// Check if user exists by email
		existingUserByEmail, err := s.userRepo.GetByEmail(userInfo.Email)
		if err != nil && !strings.Contains(err.Error(), "user not found") {
			return nil, fmt.Errorf("failed to check existing user by email: %w", err)
		}
		
		if existingUserByEmail != nil {
			// User exists with same email, link OAuth account
			user = existingUserByEmail
			user.AvatarURL = &userInfo.AvatarURL
			user.IsEmailVerified = userInfo.EmailVerified
			isNewUser = false
		} else {
			// Create new user
			user = &model.User{
				ID:              uuid.New(),
				Username:        s.generateUsername(userInfo.Username, userInfo.Name),
				Email:           userInfo.Email,
				PasswordHash:    nil, // OAuth users don't have passwords
				Role:            model.RoleStudent,
				Provider:        provider,
				ProviderID:      &userInfo.ID,
				AvatarURL:       &userInfo.AvatarURL,
				IsEmailVerified: userInfo.EmailVerified,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			}
			
			err = s.userRepo.Create(user)
			if err != nil {
				return nil, fmt.Errorf("failed to create user: %w", err)
			}
			
			isNewUser = true
		}
	}
	
	// Create or update OAuth account
	oauthAccount := &model.OAuthAccount{
		ID:            uuid.New(),
		UserID:        user.ID,
		Provider:      provider,
		ProviderID:    userInfo.ID,
		ProviderEmail: userInfo.Email,
		AccessToken:   token.AccessToken,
		RefreshToken:  token.RefreshToken,
		ExpiresAt:     &token.Expiry,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	
	existingOAuthAccount, err := s.oauthRepo.GetOAuthAccountByProviderID(provider, userInfo.ID)
	if err != nil && !strings.Contains(err.Error(), "oauth account not found") {
		return nil, fmt.Errorf("failed to check existing oauth account: %w", err)
	}
	
	if existingOAuthAccount != nil {
		// Update existing OAuth account
		existingOAuthAccount.AccessToken = token.AccessToken
		existingOAuthAccount.RefreshToken = token.RefreshToken
		existingOAuthAccount.ExpiresAt = &token.Expiry
		existingOAuthAccount.UpdatedAt = time.Now()
		
		err = s.oauthRepo.UpdateOAuthAccount(existingOAuthAccount)
		if err != nil {
			return nil, fmt.Errorf("failed to update oauth account: %w", err)
		}
	} else {
		// Create new OAuth account
		err = s.oauthRepo.CreateOAuthAccount(oauthAccount)
		if err != nil {
			return nil, fmt.Errorf("failed to create oauth account: %w", err)
		}
	}
	
	// Generate JWT token
	jwtToken, err := s.authService.generateToken(user)
	if err != nil {
		return nil, fmt.Errorf("failed to generate JWT token: %w", err)
	}
	
	// Get all linked accounts
	linkedAccounts, err := s.getLinkedProviders(user.ID)
	if err != nil {
		s.logger.Error(fmt.Errorf("failed to get linked accounts: %w", err))
		linkedAccounts = []model.OAuthProvider{provider}
	}
	
	s.logger.Infof("OAuth login successful for user %s via %s", user.Email, provider)
	
	return &model.OAuthLoginResponse{
		User:           *user,
		Token:          jwtToken,
		IsNewUser:      isNewUser,
		LinkedAccounts: linkedAccounts,
	}, nil
}

func (s *OAuthService) LinkAccount(ctx context.Context, userID uuid.UUID, provider model.OAuthProvider, code string) error {
	config, exists := s.oauthConfigs[provider]
	if !exists {
		return fmt.Errorf("oauth provider %s not configured", provider)
	}
	
	// Exchange code for token
	token, err := config.Exchange(ctx, code)
	if err != nil {
		return fmt.Errorf("failed to exchange code for token: %w", err)
	}
	
	// Get user info from provider
	userInfo, err := s.getUserInfoFromProvider(ctx, provider, token.AccessToken)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}
	
	// Check if this OAuth account is already linked to another user
	existingAccount, err := s.oauthRepo.GetOAuthAccountByProviderID(provider, userInfo.ID)
	if err != nil && !strings.Contains(err.Error(), "oauth account not found") {
		return fmt.Errorf("failed to check existing oauth account: %w", err)
	}
	
	if existingAccount != nil {
		return fmt.Errorf("this %s account is already linked to another user", provider)
	}
	
	// Create OAuth account
	oauthAccount := &model.OAuthAccount{
		ID:            uuid.New(),
		UserID:        userID,
		Provider:      provider,
		ProviderID:    userInfo.ID,
		ProviderEmail: userInfo.Email,
		AccessToken:   token.AccessToken,
		RefreshToken:  token.RefreshToken,
		ExpiresAt:     &token.Expiry,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
	
	err = s.oauthRepo.CreateOAuthAccount(oauthAccount)
	if err != nil {
		return fmt.Errorf("failed to link oauth account: %w", err)
	}
	
	s.logger.Infof("Linked %s account to user %s", provider, userID)
	return nil
}

func (s *OAuthService) UnlinkAccount(userID uuid.UUID, provider model.OAuthProvider) error {
	err := s.oauthRepo.UnlinkOAuthAccount(userID, provider)
	if err != nil {
		return fmt.Errorf("failed to unlink oauth account: %w", err)
	}
	
	s.logger.Infof("Unlinked %s account from user %s", provider, userID)
	return nil
}

func (s *OAuthService) GetLinkedAccounts(userID uuid.UUID) ([]model.OAuthProvider, error) {
	return s.getLinkedProviders(userID)
}

func (s *OAuthService) getUserInfoFromProvider(ctx context.Context, provider model.OAuthProvider, accessToken string) (*model.OAuthUserInfo, error) {
	var url string
	
	switch provider {
	case model.ProviderGoogle:
		url = "https://www.googleapis.com/oauth2/v2/userinfo"
	case model.ProviderGitHub:
		url = "https://api.github.com/user"
	case model.ProviderFacebook:
		url = "https://graph.facebook.com/me?fields=id,name,email,picture"
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+accessToken)
	
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user info: status %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	
	return s.parseUserInfo(provider, body)
}

func (s *OAuthService) parseUserInfo(provider model.OAuthProvider, data []byte) (*model.OAuthUserInfo, error) {
	var userInfo model.OAuthUserInfo
	
	switch provider {
	case model.ProviderGoogle:
		var googleUser struct {
			ID            string `json:"id"`
			Email         string `json:"email"`
			Name          string `json:"name"`
			Picture       string `json:"picture"`
			VerifiedEmail bool   `json:"verified_email"`
		}
		
		err := json.Unmarshal(data, &googleUser)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Google user info: %w", err)
		}
		
		userInfo = model.OAuthUserInfo{
			ID:            googleUser.ID,
			Email:         googleUser.Email,
			Name:          googleUser.Name,
			Username:      googleUser.Email,
			AvatarURL:     googleUser.Picture,
			EmailVerified: googleUser.VerifiedEmail,
		}
		
	case model.ProviderGitHub:
		var githubUser struct {
			ID        int    `json:"id"`
			Login     string `json:"login"`
			Name      string `json:"name"`
			Email     string `json:"email"`
			AvatarURL string `json:"avatar_url"`
		}
		
		err := json.Unmarshal(data, &githubUser)
		if err != nil {
			return nil, fmt.Errorf("failed to parse GitHub user info: %w", err)
		}
		
		userInfo = model.OAuthUserInfo{
			ID:            fmt.Sprintf("%d", githubUser.ID),
			Email:         githubUser.Email,
			Name:          githubUser.Name,
			Username:      githubUser.Login,
			AvatarURL:     githubUser.AvatarURL,
			EmailVerified: true, // GitHub emails are considered verified
		}
		
	case model.ProviderFacebook:
		var facebookUser struct {
			ID      string `json:"id"`
			Name    string `json:"name"`
			Email   string `json:"email"`
			Picture struct {
				Data struct {
					URL string `json:"url"`
				} `json:"data"`
			} `json:"picture"`
		}
		
		err := json.Unmarshal(data, &facebookUser)
		if err != nil {
			return nil, fmt.Errorf("failed to parse Facebook user info: %w", err)
		}
		
		userInfo = model.OAuthUserInfo{
			ID:            facebookUser.ID,
			Email:         facebookUser.Email,
			Name:          facebookUser.Name,
			Username:      facebookUser.Email,
			AvatarURL:     facebookUser.Picture.Data.URL,
			EmailVerified: true, // Facebook emails are considered verified
		}
		
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
	
	return &userInfo, nil
}

func (s *OAuthService) generateUsername(username, name string) string {
	if username != "" {
		return username
	}
	
	if name != "" {
		// Convert name to username format
		return strings.ToLower(strings.ReplaceAll(name, " ", ""))
	}
	
	// Generate random username
	return fmt.Sprintf("user_%s", uuid.New().String()[:8])
}

func (s *OAuthService) getLinkedProviders(userID uuid.UUID) ([]model.OAuthProvider, error) {
	accounts, err := s.oauthRepo.GetOAuthAccountsByUserID(userID)
	if err != nil {
		return nil, err
	}
	
	providers := make([]model.OAuthProvider, len(accounts))
	for i, account := range accounts {
		providers[i] = account.Provider
	}
	
	return providers, nil
}