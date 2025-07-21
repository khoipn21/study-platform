package model

import (
	"time"

	"github.com/google/uuid"
)

type OAuthProvider string

const (
	ProviderLocal    OAuthProvider = "local"
	ProviderGoogle   OAuthProvider = "google"
	ProviderGitHub   OAuthProvider = "github"
	ProviderFacebook OAuthProvider = "facebook"
)

type OAuthAccount struct {
	ID            uuid.UUID     `json:"id" db:"id"`
	UserID        uuid.UUID     `json:"user_id" db:"user_id"`
	Provider      OAuthProvider `json:"provider" db:"provider"`
	ProviderID    string        `json:"provider_id" db:"provider_id"`
	ProviderEmail string        `json:"provider_email" db:"provider_email"`
	AccessToken   string        `json:"-" db:"access_token"`
	RefreshToken  string        `json:"-" db:"refresh_token"`
	ExpiresAt     *time.Time    `json:"expires_at" db:"expires_at"`
	CreatedAt     time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at" db:"updated_at"`
}

type OAuthUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	Name          string `json:"name"`
	Username      string `json:"username"`
	AvatarURL     string `json:"avatar_url"`
	EmailVerified bool   `json:"email_verified"`
}

type OAuthConfig struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
	Scopes       []string
}

type OAuthLoginRequest struct {
	Provider OAuthProvider `json:"provider" validate:"required"`
	Code     string        `json:"code" validate:"required"`
	State    string        `json:"state" validate:"required"`
}

type OAuthLoginResponse struct {
	User         User   `json:"user"`
	Token        string `json:"token"`
	IsNewUser    bool   `json:"is_new_user"`
	LinkedAccounts []OAuthProvider `json:"linked_accounts"`
}

func (p OAuthProvider) String() string {
	return string(p)
}

func (p OAuthProvider) IsValid() bool {
	switch p {
	case ProviderLocal, ProviderGoogle, ProviderGitHub, ProviderFacebook:
		return true
	default:
		return false
	}
}

func (p OAuthProvider) GetAuthURL() string {
	switch p {
	case ProviderGoogle:
		return "https://accounts.google.com/o/oauth2/auth"
	case ProviderGitHub:
		return "https://github.com/login/oauth/authorize"
	case ProviderFacebook:
		return "https://www.facebook.com/v18.0/dialog/oauth"
	default:
		return ""
	}
}

func (p OAuthProvider) GetTokenURL() string {
	switch p {
	case ProviderGoogle:
		return "https://oauth2.googleapis.com/token"
	case ProviderGitHub:
		return "https://github.com/login/oauth/access_token"
	case ProviderFacebook:
		return "https://graph.facebook.com/v18.0/oauth/access_token"
	default:
		return ""
	}
}

func (p OAuthProvider) GetUserInfoURL() string {
	switch p {
	case ProviderGoogle:
		return "https://www.googleapis.com/oauth2/v2/userinfo"
	case ProviderGitHub:
		return "https://api.github.com/user"
	case ProviderFacebook:
		return "https://graph.facebook.com/me?fields=id,name,email,picture"
	default:
		return ""
	}
}