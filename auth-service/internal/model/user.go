package model

import (
	"time"

	"github.com/google/uuid"
)

type UserRole string

const (
	RoleAdmin      UserRole = "admin"
	RoleInstructor UserRole = "instructor"
	RoleStudent    UserRole = "student"
)

type User struct {
	ID              uuid.UUID     `json:"id" db:"id"`
	Username        string        `json:"username" db:"username"`
	Email           string        `json:"email" db:"email"`
	PasswordHash    *string       `json:"-" db:"password_hash"`
	Role            UserRole      `json:"role" db:"role"`
	Provider        OAuthProvider `json:"provider" db:"provider"`
	ProviderID      *string       `json:"provider_id" db:"provider_id"`
	AvatarURL       *string       `json:"avatar_url" db:"avatar_url"`
	IsEmailVerified bool          `json:"is_email_verified" db:"is_email_verified"`
	CreatedAt       time.Time     `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time     `json:"updated_at" db:"updated_at"`
}

type CreateUserRequest struct {
	Username string   `json:"username" validate:"required,min=3,max=50"`
	Email    string   `json:"email" validate:"required,email"`
	Password string   `json:"password" validate:"required,min=8"`
	Role     UserRole `json:"role" validate:"required"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	User  User   `json:"user"`
	Token string `json:"token"`
}

func (r UserRole) String() string {
	return string(r)
}

func (r UserRole) IsValid() bool {
	switch r {
	case RoleAdmin, RoleInstructor, RoleStudent:
		return true
	default:
		return false
	}
}