package repository

import (
	"context"
	"time"

	"github.com/study-platform/payment-service/internal/model"
)

// Simple SubscriptionRepository implementation for missing methods
func (r *SubscriptionRepository) CreateSubscription(ctx context.Context, subscription *model.Subscription) error {
	query := `
		INSERT INTO subscriptions (id, user_id, provider_id, plan_name, status,
			billing_period, next_billing_date, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	now := time.Now()
	if subscription.CreatedAt.IsZero() {
		subscription.CreatedAt = now
	}
	if subscription.UpdatedAt.IsZero() {
		subscription.UpdatedAt = now
	}

	_, err := r.db.ExecContext(ctx, query,
		subscription.ID, subscription.UserID, subscription.ProviderID,
		subscription.PlanName, subscription.Status, subscription.BillingPeriod,
		subscription.NextBillingDate, subscription.Price,
		subscription.CreatedAt, subscription.UpdatedAt)

	return err
}