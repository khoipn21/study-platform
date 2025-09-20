package repository

import (
	"database/sql"
	"fmt"

	"github.com/study-platform/payment-service/internal/model"
)

type PaymentMethodRepository struct {
	db *sql.DB
}

func NewPaymentMethodRepository(db *sql.DB) *PaymentMethodRepository {
	return &PaymentMethodRepository{db: db}
}

func (r *PaymentMethodRepository) Create(pm *model.PaymentMethod) error {
	query := `
		INSERT INTO payment_methods (id, user_id, provider, token, card_last_four, card_expiry, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(query, pm.ID, pm.UserID, pm.Provider, pm.Token, pm.CardLastFour, pm.CardExpiry, pm.IsDefault, pm.CreatedAt, pm.UpdatedAt)
	return err
}

func (r *PaymentMethodRepository) GetByID(id string) (*model.PaymentMethod, error) {
	pm := &model.PaymentMethod{}
	query := `
		SELECT id, user_id, provider, token, card_last_four, card_expiry, is_default, created_at, updated_at
		FROM payment_methods WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&pm.ID, &pm.UserID, &pm.Provider, &pm.Token, &pm.CardLastFour, &pm.CardExpiry, &pm.IsDefault, &pm.CreatedAt, &pm.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return pm, nil
}

func (r *PaymentMethodRepository) GetByUserID(userID string) ([]*model.PaymentMethod, error) {
	query := `
		SELECT id, user_id, provider, token, card_last_four, card_expiry, is_default, created_at, updated_at
		FROM payment_methods WHERE user_id = $1 ORDER BY is_default DESC, created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var paymentMethods []*model.PaymentMethod
	for rows.Next() {
		pm := &model.PaymentMethod{}
		err := rows.Scan(
			&pm.ID, &pm.UserID, &pm.Provider, &pm.Token, &pm.CardLastFour, &pm.CardExpiry, &pm.IsDefault, &pm.CreatedAt, &pm.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		paymentMethods = append(paymentMethods, pm)
	}

	return paymentMethods, nil
}

func (r *PaymentMethodRepository) Update(pm *model.PaymentMethod) error {
	query := `
		UPDATE payment_methods 
		SET provider = $2, token = $3, card_last_four = $4, card_expiry = $5, is_default = $6, updated_at = $7
		WHERE id = $1`

	result, err := r.db.Exec(query, pm.ID, pm.Provider, pm.Token, pm.CardLastFour, pm.CardExpiry, pm.IsDefault, pm.UpdatedAt)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment method not found")
	}

	return nil
}

func (r *PaymentMethodRepository) Delete(id string) error {
	query := `DELETE FROM payment_methods WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("payment method not found")
	}

	return nil
}

func (r *PaymentMethodRepository) SetAllNonDefault(userID string) error {
	query := `UPDATE payment_methods SET is_default = false WHERE user_id = $1`
	_, err := r.db.Exec(query, userID)
	return err
}

func (r *PaymentMethodRepository) CheckOwnership(id, userID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM payment_methods WHERE id = $1 AND user_id = $2`
	err := r.db.QueryRow(query, id, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}