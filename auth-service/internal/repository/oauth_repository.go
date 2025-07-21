package repository

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	"github.com/study-platform/auth-service/internal/model"
	"github.com/study-platform/pkg/database"
)

type OAuthRepository struct {
	db *database.DB
}

func NewOAuthRepository(db *database.DB) *OAuthRepository {
	return &OAuthRepository{db: db}
}

func (r *OAuthRepository) CreateOAuthAccount(account *model.OAuthAccount) error {
	query := `
		INSERT INTO oauth_accounts (id, user_id, provider, provider_id, provider_email, access_token, refresh_token, expires_at, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`
	
	_, err := r.db.Exec(query,
		account.ID, account.UserID, account.Provider, account.ProviderID, account.ProviderEmail,
		account.AccessToken, account.RefreshToken, account.ExpiresAt, account.CreatedAt, account.UpdatedAt,
	)
	if err != nil {
		return fmt.Errorf("failed to create oauth account: %w", err)
	}
	
	return nil
}

func (r *OAuthRepository) GetOAuthAccountByProviderID(provider model.OAuthProvider, providerID string) (*model.OAuthAccount, error) {
	query := `
		SELECT id, user_id, provider, provider_id, provider_email, access_token, refresh_token, expires_at, created_at, updated_at
		FROM oauth_accounts
		WHERE provider = $1 AND provider_id = $2
	`
	
	account := &model.OAuthAccount{}
	err := r.db.QueryRow(query, provider, providerID).Scan(
		&account.ID, &account.UserID, &account.Provider, &account.ProviderID, &account.ProviderEmail,
		&account.AccessToken, &account.RefreshToken, &account.ExpiresAt, &account.CreatedAt, &account.UpdatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("oauth account not found")
		}
		return nil, fmt.Errorf("failed to get oauth account: %w", err)
	}
	
	return account, nil
}

func (r *OAuthRepository) GetOAuthAccountsByUserID(userID uuid.UUID) ([]*model.OAuthAccount, error) {
	query := `
		SELECT id, user_id, provider, provider_id, provider_email, access_token, refresh_token, expires_at, created_at, updated_at
		FROM oauth_accounts
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get oauth accounts: %w", err)
	}
	defer rows.Close()
	
	accounts := []*model.OAuthAccount{}
	for rows.Next() {
		account := &model.OAuthAccount{}
		err := rows.Scan(
			&account.ID, &account.UserID, &account.Provider, &account.ProviderID, &account.ProviderEmail,
			&account.AccessToken, &account.RefreshToken, &account.ExpiresAt, &account.CreatedAt, &account.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan oauth account: %w", err)
		}
		accounts = append(accounts, account)
	}
	
	return accounts, nil
}

func (r *OAuthRepository) UpdateOAuthAccount(account *model.OAuthAccount) error {
	query := `
		UPDATE oauth_accounts
		SET provider_email = $1, access_token = $2, refresh_token = $3, expires_at = $4, updated_at = $5
		WHERE id = $6
	`
	
	result, err := r.db.Exec(query,
		account.ProviderEmail, account.AccessToken, account.RefreshToken, account.ExpiresAt, account.UpdatedAt, account.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update oauth account: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("oauth account not found")
	}
	
	return nil
}

func (r *OAuthRepository) DeleteOAuthAccount(accountID uuid.UUID) error {
	query := `DELETE FROM oauth_accounts WHERE id = $1`
	
	result, err := r.db.Exec(query, accountID)
	if err != nil {
		return fmt.Errorf("failed to delete oauth account: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("oauth account not found")
	}
	
	return nil
}

func (r *OAuthRepository) LinkOAuthAccount(userID uuid.UUID, account *model.OAuthAccount) error {
	account.UserID = userID
	return r.CreateOAuthAccount(account)
}

func (r *OAuthRepository) UnlinkOAuthAccount(userID uuid.UUID, provider model.OAuthProvider) error {
	query := `DELETE FROM oauth_accounts WHERE user_id = $1 AND provider = $2`
	
	result, err := r.db.Exec(query, userID, provider)
	if err != nil {
		return fmt.Errorf("failed to unlink oauth account: %w", err)
	}
	
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}
	
	if rowsAffected == 0 {
		return fmt.Errorf("oauth account not found")
	}
	
	return nil
}