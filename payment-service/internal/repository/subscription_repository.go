package repository

import (
	"database/sql"
	"fmt"

	"payment-service/internal/model"
)

type SubscriptionRepository struct {
	db *sql.DB
}

func NewSubscriptionRepository(db *sql.DB) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

func (r *SubscriptionRepository) Create(sub *model.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, user_id, payment_method_id, plan_name, status, billing_period, next_billing_date, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.db.Exec(query, sub.ID, sub.UserID, sub.PaymentMethodID, sub.PlanName, sub.Status, sub.BillingPeriod, sub.NextBillingDate, sub.Price, sub.CreatedAt, sub.UpdatedAt)
	return err
}

func (r *SubscriptionRepository) GetByID(id string) (*model.Subscription, error) {
	sub := &model.Subscription{}
	query := `
		SELECT id, user_id, payment_method_id, plan_name, status, billing_period, next_billing_date, price, created_at, updated_at
		FROM subscriptions WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(
		&sub.ID, &sub.UserID, &sub.PaymentMethodID, &sub.PlanName, &sub.Status, &sub.BillingPeriod, &sub.NextBillingDate, &sub.Price, &sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func (r *SubscriptionRepository) GetByUserID(userID string) ([]*model.Subscription, error) {
	query := `
		SELECT id, user_id, payment_method_id, plan_name, status, billing_period, next_billing_date, price, created_at, updated_at
		FROM subscriptions WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*model.Subscription
	for rows.Next() {
		sub := &model.Subscription{}
		err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.PaymentMethodID, &sub.PlanName, &sub.Status, &sub.BillingPeriod, &sub.NextBillingDate, &sub.Price, &sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) Update(sub *model.Subscription) error {
	query := `
		UPDATE subscriptions 
		SET payment_method_id = $2, plan_name = $3, status = $4, billing_period = $5, next_billing_date = $6, price = $7, updated_at = $8
		WHERE id = $1`

	result, err := r.db.Exec(query, sub.ID, sub.PaymentMethodID, sub.PlanName, sub.Status, sub.BillingPeriod, sub.NextBillingDate, sub.Price, sub.UpdatedAt)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (r *SubscriptionRepository) Delete(id string) error {
	query := `DELETE FROM subscriptions WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("subscription not found")
	}

	return nil
}

func (r *SubscriptionRepository) GetByStatus(status string) ([]*model.Subscription, error) {
	query := `
		SELECT id, user_id, payment_method_id, plan_name, status, billing_period, next_billing_date, price, created_at, updated_at
		FROM subscriptions WHERE status = $1 ORDER BY created_at DESC`

	rows, err := r.db.Query(query, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*model.Subscription
	for rows.Next() {
		sub := &model.Subscription{}
		err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.PaymentMethodID, &sub.PlanName, &sub.Status, &sub.BillingPeriod, &sub.NextBillingDate, &sub.Price, &sub.CreatedAt, &sub.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, sub)
	}

	return subscriptions, nil
}

func (r *SubscriptionRepository) CheckOwnership(id, userID string) (bool, error) {
	var count int
	query := `SELECT COUNT(*) FROM subscriptions WHERE id = $1 AND user_id = $2`
	err := r.db.QueryRow(query, id, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}