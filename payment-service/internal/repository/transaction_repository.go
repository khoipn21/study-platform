package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/study-platform/payment-service/internal/model"
)

type TransactionRepository struct {
	db *sql.DB
}

func NewTransactionRepository(db *sql.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) Create(ctx context.Context, tx *model.Transaction) error {
	customDataJSON, err := json.Marshal(tx.CustomData)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO transactions (
			id, user_id, course_id, payment_method_id, amount, currency, status,
			transaction_reference, lemon_squeezy_order_id, lemon_squeezy_checkout_id,
			webhook_event_id, payment_verified_at, payment_provider, custom_data, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	_, err = r.db.ExecContext(ctx, query,
		tx.ID, tx.UserID, tx.CourseID, tx.PaymentMethodID, tx.Amount, tx.Currency,
		tx.Status, tx.TransactionReference, tx.LemonSqueezyOrderID,
		tx.LemonSqueezyCheckoutID, tx.WebhookEventID, tx.PaymentVerifiedAt,
		tx.PaymentProvider, customDataJSON, tx.CreatedAt, tx.UpdatedAt)
	return err
}

func (r *TransactionRepository) GetByID(ctx context.Context, id string) (*model.Transaction, error) {
	tx := &model.Transaction{}
	var customDataJSON []byte

	query := `
		SELECT id, user_id, course_id, payment_method_id, amount, currency, status,
		       transaction_reference, lemon_squeezy_order_id, lemon_squeezy_checkout_id,
		       webhook_event_id, payment_verified_at, payment_provider, custom_data, created_at, updated_at
		FROM transactions WHERE id = $1`

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&tx.ID, &tx.UserID, &tx.CourseID, &tx.PaymentMethodID, &tx.Amount, &tx.Currency,
		&tx.Status, &tx.TransactionReference, &tx.LemonSqueezyOrderID,
		&tx.LemonSqueezyCheckoutID, &tx.WebhookEventID, &tx.PaymentVerifiedAt,
		&tx.PaymentProvider, &customDataJSON, &tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if customDataJSON != nil {
		if err := json.Unmarshal(customDataJSON, &tx.CustomData); err != nil {
			return nil, err
		}
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
	customDataJSON, err := json.Marshal(tx.CustomData)
	if err != nil {
		return err
	}

	query := `
		UPDATE transactions
		SET course_id = $2, payment_method_id = $3, amount = $4, currency = $5, status = $6,
		    transaction_reference = $7, lemon_squeezy_order_id = $8, lemon_squeezy_checkout_id = $9,
		    webhook_event_id = $10, custom_data = $11, updated_at = $12
		WHERE id = $1`

	result, err := r.db.Exec(query,
		tx.ID, tx.CourseID, tx.PaymentMethodID, tx.Amount, tx.Currency, tx.Status,
		tx.TransactionReference, tx.LemonSqueezyOrderID, tx.LemonSqueezyCheckoutID,
		tx.WebhookEventID, customDataJSON, tx.UpdatedAt)
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
	var customDataJSON []byte

	query := `
		SELECT id, user_id, course_id, payment_method_id, amount, currency, status,
		       transaction_reference, lemon_squeezy_order_id, lemon_squeezy_checkout_id,
		       webhook_event_id, custom_data, created_at, updated_at
		FROM transactions WHERE course_id = $1 AND user_id = $2 AND status = $3`

	err := r.db.QueryRow(query, courseID, userID, model.TransactionStatusCompleted).Scan(
		&tx.ID, &tx.UserID, &tx.CourseID, &tx.PaymentMethodID, &tx.Amount, &tx.Currency,
		&tx.Status, &tx.TransactionReference, &tx.LemonSqueezyOrderID,
		&tx.LemonSqueezyCheckoutID, &tx.WebhookEventID, &customDataJSON,
		&tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if customDataJSON != nil {
		if err := json.Unmarshal(customDataJSON, &tx.CustomData); err != nil {
			return nil, err
		}
	}

	return tx, nil
}

// GetByLemonSqueezyOrderID finds a transaction by Lemon Squeezy order ID
func (r *TransactionRepository) GetByLemonSqueezyOrderID(orderID string) (*model.Transaction, error) {
	tx := &model.Transaction{}
	var customDataJSON []byte

	query := `
		SELECT id, user_id, course_id, payment_method_id, amount, currency, status,
		       transaction_reference, lemon_squeezy_order_id, lemon_squeezy_checkout_id,
		       webhook_event_id, custom_data, created_at, updated_at
		FROM transactions WHERE lemon_squeezy_order_id = $1`

	err := r.db.QueryRow(query, orderID).Scan(
		&tx.ID, &tx.UserID, &tx.CourseID, &tx.PaymentMethodID, &tx.Amount, &tx.Currency,
		&tx.Status, &tx.TransactionReference, &tx.LemonSqueezyOrderID,
		&tx.LemonSqueezyCheckoutID, &tx.WebhookEventID, &customDataJSON,
		&tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if customDataJSON != nil {
		if err := json.Unmarshal(customDataJSON, &tx.CustomData); err != nil {
			return nil, err
		}
	}

	return tx, nil
}

// GetByLemonSqueezyCheckoutID finds a transaction by Lemon Squeezy checkout ID
func (r *TransactionRepository) GetByLemonSqueezyCheckoutID(checkoutID string) (*model.Transaction, error) {
	tx := &model.Transaction{}
	var customDataJSON []byte

	query := `
		SELECT id, user_id, course_id, payment_method_id, amount, currency, status,
		       transaction_reference, lemon_squeezy_order_id, lemon_squeezy_checkout_id,
		       webhook_event_id, custom_data, created_at, updated_at
		FROM transactions WHERE lemon_squeezy_checkout_id = $1`

	err := r.db.QueryRow(query, checkoutID).Scan(
		&tx.ID, &tx.UserID, &tx.CourseID, &tx.PaymentMethodID, &tx.Amount, &tx.Currency,
		&tx.Status, &tx.TransactionReference, &tx.LemonSqueezyOrderID,
		&tx.LemonSqueezyCheckoutID, &tx.WebhookEventID, &customDataJSON,
		&tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if customDataJSON != nil {
		if err := json.Unmarshal(customDataJSON, &tx.CustomData); err != nil {
			return nil, err
		}
	}

	return tx, nil
}

// GetPendingByCourseAndUser finds pending transactions for a specific course and user
func (r *TransactionRepository) GetPendingByCourseAndUser(ctx context.Context, courseID, userID string) (*model.Transaction, error) {
	tx := &model.Transaction{}
	var customDataJSON []byte

	query := `
		SELECT id, user_id, course_id, payment_method_id, amount, currency, status,
		       transaction_reference, lemon_squeezy_order_id, lemon_squeezy_checkout_id,
		       webhook_event_id, payment_verified_at, payment_provider, custom_data, created_at, updated_at
		FROM transactions WHERE course_id = $1 AND user_id = $2 AND status = $3
		ORDER BY created_at DESC LIMIT 1`

	err := r.db.QueryRowContext(ctx, query, courseID, userID, model.TransactionStatusPending).Scan(
		&tx.ID, &tx.UserID, &tx.CourseID, &tx.PaymentMethodID, &tx.Amount, &tx.Currency,
		&tx.Status, &tx.TransactionReference, &tx.LemonSqueezyOrderID,
		&tx.LemonSqueezyCheckoutID, &tx.WebhookEventID, &tx.PaymentVerifiedAt,
		&tx.PaymentProvider, &customDataJSON, &tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if customDataJSON != nil {
		if err := json.Unmarshal(customDataJSON, &tx.CustomData); err != nil {
			return nil, err
		}
	}

	return tx, nil
}

// GetByStripePaymentIntentID finds a transaction by Stripe payment intent ID
func (r *TransactionRepository) GetByStripePaymentIntentID(paymentIntentID string) (*model.Transaction, error) {
	tx := &model.Transaction{}
	var customDataJSON []byte

	query := `
		SELECT id, user_id, course_id, payment_method_id, amount, currency, status,
		       transaction_reference, lemon_squeezy_order_id, lemon_squeezy_checkout_id,
		       webhook_event_id, stripe_payment_intent_id, stripe_customer_id,
		       stripe_charge_id, stripe_session_id, stripe_invoice_id,
		       stripe_subscription_id, payment_verified_at, payment_provider,
		       custom_data, created_at, updated_at
		FROM transactions WHERE stripe_payment_intent_id = $1`

	err := r.db.QueryRow(query, paymentIntentID).Scan(
		&tx.ID, &tx.UserID, &tx.CourseID, &tx.PaymentMethodID, &tx.Amount, &tx.Currency,
		&tx.Status, &tx.TransactionReference, &tx.LemonSqueezyOrderID,
		&tx.LemonSqueezyCheckoutID, &tx.WebhookEventID, &tx.StripePaymentIntentID,
		&tx.StripeCustomerID, &tx.StripeChargeID, &tx.StripeSessionID,
		&tx.StripeInvoiceID, &tx.StripeSubscriptionID, &tx.PaymentVerifiedAt,
		&tx.PaymentProvider, &customDataJSON, &tx.CreatedAt, &tx.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if customDataJSON != nil {
		if err := json.Unmarshal(customDataJSON, &tx.CustomData); err != nil {
			return nil, err
		}
	}

	return tx, nil
}

// GetByStripeCustomerID finds transactions by Stripe customer ID
func (r *TransactionRepository) GetByStripeCustomerID(customerID string) ([]*model.Transaction, error) {
	query := `
		SELECT id, user_id, course_id, payment_method_id, amount, currency, status,
		       transaction_reference, lemon_squeezy_order_id, lemon_squeezy_checkout_id,
		       webhook_event_id, stripe_payment_intent_id, stripe_customer_id,
		       stripe_charge_id, stripe_session_id, stripe_invoice_id,
		       stripe_subscription_id, payment_verified_at, payment_provider,
		       custom_data, created_at, updated_at
		FROM transactions WHERE stripe_customer_id = $1
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		tx := &model.Transaction{}
		var customDataJSON []byte

		err := rows.Scan(
			&tx.ID, &tx.UserID, &tx.CourseID, &tx.PaymentMethodID, &tx.Amount, &tx.Currency,
			&tx.Status, &tx.TransactionReference, &tx.LemonSqueezyOrderID,
			&tx.LemonSqueezyCheckoutID, &tx.WebhookEventID, &tx.StripePaymentIntentID,
			&tx.StripeCustomerID, &tx.StripeChargeID, &tx.StripeSessionID,
			&tx.StripeInvoiceID, &tx.StripeSubscriptionID, &tx.PaymentVerifiedAt,
			&tx.PaymentProvider, &customDataJSON, &tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		if customDataJSON != nil {
			if err := json.Unmarshal(customDataJSON, &tx.CustomData); err != nil {
				return nil, err
			}
		}

		transactions = append(transactions, tx)
	}

	return transactions, nil
}

// BeginTx starts a database transaction
func (r *TransactionRepository) BeginTx(ctx context.Context) (*sql.Tx, error) {
	return r.db.BeginTx(ctx, nil)
}

// CreateWithTx creates a transaction within a database transaction
func (r *TransactionRepository) CreateWithTx(ctx context.Context, dbTx *sql.Tx, tx *model.Transaction) error {
	customDataJSON, err := json.Marshal(tx.CustomData)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO transactions (
			id, user_id, course_id, payment_method_id, amount, currency, status,
			transaction_reference, lemon_squeezy_order_id, lemon_squeezy_checkout_id,
			webhook_event_id, payment_verified_at, payment_provider, custom_data, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16)`

	_, err = dbTx.ExecContext(ctx, query,
		tx.ID, tx.UserID, tx.CourseID, tx.PaymentMethodID, tx.Amount, tx.Currency,
		tx.Status, tx.TransactionReference, tx.LemonSqueezyOrderID,
		tx.LemonSqueezyCheckoutID, tx.WebhookEventID, tx.PaymentVerifiedAt,
		tx.PaymentProvider, customDataJSON, tx.CreatedAt, tx.UpdatedAt)
	return err
}

// GetTransactionByOrderID retrieves a transaction by LemonSqueezy order ID
func (r *TransactionRepository) GetTransactionByOrderID(orderID string) (*model.Transaction, error) {
	tx := &model.Transaction{}
	var customDataJSON []byte

	query := `
		SELECT id, user_id, course_id, payment_method_id, amount, currency, status,
		       transaction_reference, lemon_squeezy_order_id, lemon_squeezy_checkout_id,
		       webhook_event_id, payment_verified_at, payment_provider, custom_data, created_at, updated_at
		FROM transactions WHERE lemon_squeezy_order_id = $1`

	err := r.db.QueryRow(query, orderID).Scan(
		&tx.ID, &tx.UserID, &tx.CourseID, &tx.PaymentMethodID, &tx.Amount, &tx.Currency,
		&tx.Status, &tx.TransactionReference, &tx.LemonSqueezyOrderID,
		&tx.LemonSqueezyCheckoutID, &tx.WebhookEventID, &tx.PaymentVerifiedAt,
		&tx.PaymentProvider, &customDataJSON, &tx.CreatedAt, &tx.UpdatedAt)

	if err != nil {
		return nil, err
	}

	if len(customDataJSON) > 0 {
		err = json.Unmarshal(customDataJSON, &tx.CustomData)
		if err != nil {
			return nil, err
		}
	}

	return tx, nil
}

// UpdateTransaction updates an existing transaction
func (r *TransactionRepository) UpdateTransaction(tx *model.Transaction) error {
	customDataJSON, err := json.Marshal(tx.CustomData)
	if err != nil {
		return err
	}

	query := `
		UPDATE transactions
		SET status = $1, payment_verified_at = $2, webhook_event_id = $3,
		    custom_data = $4, updated_at = $5
		WHERE id = $6`

	_, err = r.db.Exec(query, tx.Status, tx.PaymentVerifiedAt, tx.WebhookEventID,
		customDataJSON, tx.UpdatedAt, tx.ID)
	return err
}
// GetByCheckoutID retrieves a transaction by checkout ID
func (r *TransactionRepository) GetByCheckoutID(ctx context.Context, checkoutID string) (*model.Transaction, error) {
	query := `
		SELECT id, user_id, course_id, lemon_squeezy_checkout_id, amount, currency, status,
		       payment_provider, transaction_reference, custom_data, created_at, updated_at
		FROM transactions WHERE lemon_squeezy_checkout_id = $1`

	row := r.db.QueryRowContext(ctx, query, checkoutID)

	transaction := &model.Transaction{}
	var providerDataBytes []byte

	err := row.Scan(
		&transaction.ID, &transaction.UserID, &transaction.CourseID,
		&transaction.LemonSqueezyCheckoutID, &transaction.Amount, &transaction.Currency,
		&transaction.Status, &transaction.PaymentProvider, &transaction.TransactionReference,
		&providerDataBytes, &transaction.CreatedAt, &transaction.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("transaction not found")
		}
		return nil, err
	}

	if len(providerDataBytes) > 0 {
		if err := json.Unmarshal(providerDataBytes, &transaction.CustomData); err != nil {
			return nil, err
		}
	}

	return transaction, nil
}

// GetUserTransactions retrieves transactions for a specific user
func (r *TransactionRepository) GetUserTransactions(ctx context.Context, userID string) ([]*model.Transaction, error) {
	transactions, _, err := r.GetUserTransactionsWithFilters(ctx, userID, 0, 0, "", "")
	return transactions, err
}

// GetUserTransactionsWithFilters retrieves transactions for a specific user with pagination and filters
func (r *TransactionRepository) GetUserTransactionsWithFilters(ctx context.Context, userID string, limit, offset int, statusFilter, courseIDFilter string) ([]*model.Transaction, int64, error) {
	// Build query with filters
	baseQuery := `
		SELECT id, user_id, course_id, payment_method_id, amount, currency, status,
		       transaction_reference, lemon_squeezy_order_id, lemon_squeezy_checkout_id,
		       webhook_event_id, payment_verified_at, payment_provider, custom_data, created_at, updated_at
		FROM transactions WHERE user_id = $1`

	countQuery := `SELECT COUNT(*) FROM transactions WHERE user_id = $1`

	args := []interface{}{userID}
	argIndex := 2

	// Add status filter
	if statusFilter != "" {
		baseQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND status = $%d", argIndex)
		args = append(args, statusFilter)
		argIndex++
	}

	// Add course ID filter
	if courseIDFilter != "" {
		baseQuery += fmt.Sprintf(" AND course_id = $%d", argIndex)
		countQuery += fmt.Sprintf(" AND course_id = $%d", argIndex)
		args = append(args, courseIDFilter)
		argIndex++
	}

	// Add ordering and pagination
	baseQuery += " ORDER BY created_at DESC"
	if limit > 0 {
		baseQuery += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, limit)
		argIndex++

		if offset > 0 {
			baseQuery += fmt.Sprintf(" OFFSET $%d", argIndex)
			args = append(args, offset)
		}
	}

	// Get total count (exclude limit/offset args)
	countArgs := args
	if limit > 0 {
		if offset > 0 {
			countArgs = args[:len(args)-2] // Remove limit and offset
		} else {
			countArgs = args[:len(args)-1] // Remove limit only
		}
	}

	var totalCount int64
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	// Get transactions
	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var transactions []*model.Transaction
	for rows.Next() {
		tx := &model.Transaction{}
		var customDataJSON []byte

		err := rows.Scan(
			&tx.ID, &tx.UserID, &tx.CourseID, &tx.PaymentMethodID, &tx.Amount, &tx.Currency,
			&tx.Status, &tx.TransactionReference, &tx.LemonSqueezyOrderID,
			&tx.LemonSqueezyCheckoutID, &tx.WebhookEventID, &tx.PaymentVerifiedAt,
			&tx.PaymentProvider, &customDataJSON, &tx.CreatedAt, &tx.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}

		if customDataJSON != nil {
			if err := json.Unmarshal(customDataJSON, &tx.CustomData); err != nil {
				return nil, 0, err
			}
		}

		transactions = append(transactions, tx)
	}

	return transactions, totalCount, nil
}
