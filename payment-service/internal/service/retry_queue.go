package service

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/study-platform/payment-service/internal/model"
)

// RetryQueue manages webhook retry logic
type RetryQueue struct {
	mu        sync.RWMutex
	retries   map[string]*WebhookRetryInfo
	maxRetries int
	logger    log.Logger
}

// NewRetryQueue creates a new retry queue
func NewRetryQueue(maxRetries int) *RetryQueue {
	return &RetryQueue{
		retries:    make(map[string]*WebhookRetryInfo),
		maxRetries: maxRetries,
	}
}

// AddWebhookRetry adds a webhook to the retry queue
func (rq *RetryQueue) AddWebhookRetry(ctx context.Context, payload *model.LemonSqueezyWebhookPayload, signature string, retryCount int) error {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	eventID := payload.Data.ID

	// Calculate next retry time with exponential backoff
	nextRetryAt := time.Now().Add(time.Duration(retryCount*retryCount) * time.Minute)

	rq.retries[eventID] = &WebhookRetryInfo{
		EventID:     eventID,
		Payload:     payload,
		Signature:   signature,
		RetryCount:  retryCount,
		NextRetryAt: nextRetryAt,
	}

	return nil
}

// GetPendingRetries returns webhooks that are ready for retry
func (rq *RetryQueue) GetPendingRetries(ctx context.Context) []*WebhookRetryInfo {
	rq.mu.RLock()
	defer rq.mu.RUnlock()

	var pending []*WebhookRetryInfo
	now := time.Now()

	for _, retry := range rq.retries {
		if retry.RetryCount < rq.maxRetries && now.After(retry.NextRetryAt) {
			pending = append(pending, retry)
		}
	}

	return pending
}

// RemoveRetry removes a retry from the queue
func (rq *RetryQueue) RemoveRetry(eventID string) {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	delete(rq.retries, eventID)
}

// UpdateRetryCount updates the retry count for a webhook
func (rq *RetryQueue) UpdateRetryCount(eventID string, retryCount int) {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	if retry, exists := rq.retries[eventID]; exists {
		retry.RetryCount = retryCount
		retry.NextRetryAt = time.Now().Add(time.Duration(retryCount*retryCount) * time.Minute)
	}
}

// GetRetryInfo returns retry information for a specific event
func (rq *RetryQueue) GetRetryInfo(eventID string) *WebhookRetryInfo {
	rq.mu.RLock()
	defer rq.mu.RUnlock()

	return rq.retries[eventID]
}

// CleanupExpiredRetries removes retries that have exceeded max attempts
func (rq *RetryQueue) CleanupExpiredRetries() {
	rq.mu.Lock()
	defer rq.mu.Unlock()

	for eventID, retry := range rq.retries {
		if retry.RetryCount >= rq.maxRetries {
			delete(rq.retries, eventID)
		}
	}
}