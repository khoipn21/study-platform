package repository

import (
	"database/sql"
	"fmt"

	"payment-service/internal/model"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(tx *model.Transaction) error {
	query := `
		INSERT INTO transactions (id, user_id, course_id, payment_method_id, amount, currency, status, transaction_reference, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.Exec(query, tx.ID, tx.UserID, tx.CourseID, tx.PaymentMethodID, tx.Amount, tx.Currency, tx.Status, tx.TransactionReference, tx.CreatedAt, tx.UpdatedAt)
	return err
}

func (r *TransactionRepository) GetByID(id string) (*model.Transaction, error) {
	tx := &model.Transaction{}
	query := `
		SELECT id, user_id, course_id, payment_method_id, amount, currency, status, transaction_reference, created_at, updated_at
		FROM transactions WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&tx.ID, &tx.UserID, &tx.CourseID, &tx.PaymentMethodID, &tx.Amount, &tx.Currency, &tx.Status, &tx.TransactionReference, &tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (r *TransactionRepository) GetByUserID(userID string, limit, offset int) ([]*model.Transaction, error) {
	query := `
		SELECT id, user_id, course_id, payment_method_id, amount, currency, status, transaction_reference, created_at, updated_at
		FROM transactions WHERE user_id = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		tx := &model.Transaction{}
		err := rows.Scan(
			&tx.ID, &tx.UserID, &tx.CourseID, &tx.PaymentMethodID, &tx.Amount, &tx.Currency, &tx.Status, &tx.TransactionReference, &tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (r *TransactionRepository) GetByReference(reference string) (*model.Transaction, error) {
	tx := &model.Transaction{}
	query := `
		SELECT id, user_id, course_id, payment_method_id, amount, currency, status, transaction_reference, created_at, updated_at
		FROM transactions WHERE transaction_reference = $1`

	err := r.db.QueryRow(query, reference).Scan(
		&tx.ID, &tx.UserID, &tx.CourseID, &tx.PaymentMethodID, &tx.Amount, &tx.Currency, &tx.Status, &tx.TransactionReference, &tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}

func (r *TransactionRepository) Update(tx *model.Transaction) error {
	query := `
		UPDATE transactions 
		SET course_id = $2, payment_method_id = $3, amount = $4, currency = $5, status = $6, transaction_reference = $7, updated_at = $8
		WHERE id = $1`

	result, err := r.db.Exec(query, tx.ID, tx.CourseID, tx.PaymentMethodID, tx.Amount, tx.Currency, tx.Status, tx.TransactionReference, tx.UpdatedAt)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("transaction not found")
	}

	return nil
}

func (r *TransactionRepository) GetByStatus(status string, limit, offset int) ([]*model.Transaction, error) {
	query := `
		SELECT id, user_id, course_id, payment_method_id, amount, currency, status, transaction_reference, created_at, updated_at
		FROM transactions WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		tx := &model.Transaction{}
		err := rows.Scan(
			&tx.ID, &tx.UserID, &tx.CourseID, &tx.PaymentMethodID, &tx.Amount, &tx.Currency, &tx.Status, &tx.TransactionReference, &tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, tx)
	}

	return transactions, nil
}

func (r *TransactionRepository) CheckOwnership(id, userID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM transactions WHERE id = $1 AND user_id = $2`
	err := r.db.QueryRow(query, id, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *TransactionRepository) GetByCourseAndUser(courseID, userID string) (*model.Transaction, error) {
	tx := &model.Transaction{}
	query := `
		SELECT id, user_id, course_id, payment_method_id, amount, currency, status, transaction_reference, created_at, updated_at
		FROM transactions WHERE course_id = $1 AND user_id = $2 AND status = $3`

	err := r.db.QueryRow(query, courseID, userID, model.TransactionStatusCompleted).Scan(
		&tx.ID, &tx.UserID, &tx.CourseID, &tx.PaymentMethodID, &tx.Amount, &tx.Currency, &tx.Status, &tx.TransactionReference, &tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return tx, nil
}