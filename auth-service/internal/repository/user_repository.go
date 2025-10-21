package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/study-platform/auth-service/internal/model"
	"github.com/study-platform/pkg/database"
)

type UserRepository struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *model.User) error {
	query := `
		INSERT INTO users (id, username, email, password_hash, role, provider, provider_id, avatar_url, is_email_verified, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`
	
	_, err := r.db.Exec(query, user.ID, user.Username, user.Email, user.PasswordHash, user.Role, user.Provider, user.ProviderID, user.AvatarURL, user.IsEmailVerified, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}
	
	return nil
}

func (r *UserRepository) GetByID(id uuid.UUID) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, provider, provider_id, avatar_url, is_email_verified, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	
	user := &model.User{}
	err := r.db.QueryRow(query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.Provider, &user.ProviderID, &user.AvatarURL, &user.IsEmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}

func (r *UserRepository) GetByEmail(email string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, provider, provider_id, avatar_url, is_email_verified, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	
	user := &model.User{}
	err := r.db.QueryRow(query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.Provider, &user.ProviderID, &user.AvatarURL, &user.IsEmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}

func (r *UserRepository) GetByUsername(username string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at
		FROM users
		WHERE username = $1
	`
	
	user := &model.User{}
	err := r.db.QueryRow(query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}

func (r *UserRepository) UpdateRole(userID uuid.UUID, role model.UserRole) error {
	query := `
		UPDATE users
		SET role = $1, updated_at = NOW()
		WHERE id = $2
	`
	
	result, err := r.db.Exec(query, role, userID)
	if err != nil {
		return fmt.Errorf("failed to update user role: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

func (r *UserRepository) Delete(userID uuid.UUID) error {
	query := `DELETE FROM users WHERE id = $1`
	
	result, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	
	return nil
}

func (r *UserRepository) List(limit, offset int) ([]*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, created_at, updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}
	defer rows.Close()
	
	users := []*model.User{}
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(
			&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.CreatedAt, &user.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		users = append(users, user)
	}
	
	return users, nil
}

func (r *UserRepository) GetByProviderID(provider model.OAuthProvider, providerID string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, provider, provider_id, avatar_url, is_email_verified, created_at, updated_at
		FROM users
		WHERE provider = $1 AND provider_id = $2
	`
	
	user := &model.User{}
	err := r.db.QueryRow(query, provider, providerID).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.Role, &user.Provider, &user.ProviderID, &user.AvatarURL, &user.IsEmailVerified, &user.CreatedAt, &user.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}
	
	return user, nil
}
func (r *UserRepository) SearchUsers(query string, limit int) ([]*model.User, error) {
	sqlQuery := `
		SELECT id, username, email, role, avatar_url
		FROM users
		WHERE username ILIKE $1
		ORDER BY username ASC
		LIMIT $2
	`
	
	rows, err := r.db.Query(sqlQuery, "%"+query+"%", limit)
	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}
	defer rows.Close()
	
	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		var avatarURL sql.NullString
		
		err := rows.Scan(&user.ID, &user.Username, &user.Email, &user.Role, &avatarURL)
		if err != nil {
			return nil, fmt.Errorf("failed to scan user: %w", err)
		}
		
		if avatarURL.Valid {
			avatarStr := avatarURL.String
			user.AvatarURL = &avatarStr
		}
		
		users = append(users, user)
	}
	
	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}
	
	return users, nil
}
