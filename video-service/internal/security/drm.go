package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
)

type DRMManager struct {
	secretKey             string
	defaultExpiryDuration time.Duration
	cloudflareAccountID   string
	cloudflareAPIToken    string
}

type SignedURLClaims struct {
	VideoID   uuid.UUID `json:"video_id"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt int64     `json:"expires_at"`
	AccessLevel string  `json:"access_level"` // "view", "download", "admin"
}

type VideoAccessPolicy struct {
	AllowedCountries []string `json:"allowed_countries,omitempty"`
	BlockedCountries []string `json:"blocked_countries,omitempty"`
	RequireAuth      bool     `json:"require_auth"`
	MaxViews         int      `json:"max_views,omitempty"`
	TimeRestrictions *TimeRestriction `json:"time_restrictions,omitempty"`
}

type TimeRestriction struct {
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
}

func NewDRMManager(secretKey, cloudflareAccountID, cloudflareAPIToken string) *DRMManager {
	return &DRMManager{
		secretKey:             secretKey,
		defaultExpiryDuration: 2 * time.Hour, // Default 2-hour expiry
		cloudflareAccountID:   cloudflareAccountID,
		cloudflareAPIToken:    cloudflareAPIToken,
	}
}

// GenerateSignedURL creates a signed URL for video access
func (drm *DRMManager) GenerateSignedURL(videoID, userID uuid.UUID, accessLevel string, customExpiry *time.Duration) (string, error) {
	expiry := drm.defaultExpiryDuration
	if customExpiry != nil {
		expiry = *customExpiry
	}

	claims := SignedURLClaims{
		VideoID:     videoID,
		UserID:      userID,
		ExpiresAt:   time.Now().Add(expiry).Unix(),
		AccessLevel: accessLevel,
	}

	// Encode claims
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %w", err)
	}

	claimsB64 := base64.URLEncoding.EncodeToString(claimsJSON)

	// Generate signature
	signature := drm.generateSignature(claimsB64)

	// Create signed URL
	baseURL := fmt.Sprintf("/api/videos/%s/stream", videoID.String())
	signedURL := fmt.Sprintf("%s?token=%s&signature=%s", baseURL,
		url.QueryEscape(claimsB64),
		url.QueryEscape(signature))

	return signedURL, nil
}

// ValidateSignedURL validates a signed URL and returns the claims
func (drm *DRMManager) ValidateSignedURL(token, signature string) (*SignedURLClaims, error) {
	// Verify signature
	expectedSignature := drm.generateSignature(token)
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return nil, fmt.Errorf("invalid signature")
	}

	// Decode token
	claimsJSON, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("failed to decode token: %w", err)
	}

	var claims SignedURLClaims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, fmt.Errorf("failed to unmarshal claims: %w", err)
	}

	// Check expiry
	if time.Now().Unix() > claims.ExpiresAt {
		return nil, fmt.Errorf("token expired")
	}

	return &claims, nil
}

// GenerateCloudflareSignedURL creates a signed URL using Cloudflare Stream's signing
func (drm *DRMManager) GenerateCloudflareSignedURL(videoID string, userID uuid.UUID, expiry time.Duration) (string, error) {
	// For Cloudflare Stream, we need to use their signed URL format
	// This is a simplified implementation - in production, you'd use Cloudflare's SDK

	expiryTimestamp := time.Now().Add(expiry).Unix()

	// Create the resource path
	resourcePath := fmt.Sprintf("/%s/%s", drm.cloudflareAccountID, videoID)

	// Generate signature using Cloudflare's method
	// Note: This is a simplified version - actual implementation would use Cloudflare's exact algorithm
	h := hmac.New(sha256.New, []byte(drm.secretKey))
	h.Write([]byte(fmt.Sprintf("%s%d%s", resourcePath, expiryTimestamp, userID.String())))
	signature := base64.URLEncoding.EncodeToString(h.Sum(nil))

	// Construct the signed URL
	baseURL := fmt.Sprintf("https://customer-m033z5x00ks6nunl.cloudflarestream.com/%s", videoID)
	signedURL := fmt.Sprintf("%s?expires=%d&signature=%s&user=%s",
		baseURL, expiryTimestamp, url.QueryEscape(signature), userID.String())

	return signedURL, nil
}

// ValidateAccess checks if a user has access to a video based on policy
func (drm *DRMManager) ValidateAccess(userID uuid.UUID, videoID uuid.UUID, policy *VideoAccessPolicy, userIP, userCountry string) error {
	// Check authentication requirement
	if policy.RequireAuth && userID == uuid.Nil {
		return fmt.Errorf("authentication required")
	}

	// Check geographic restrictions
	if len(policy.AllowedCountries) > 0 {
		allowed := false
		for _, country := range policy.AllowedCountries {
			if country == userCountry {
				allowed = true
				break
			}
		}
		if !allowed {
			return fmt.Errorf("access not allowed from country: %s", userCountry)
		}
	}

	if len(policy.BlockedCountries) > 0 {
		for _, country := range policy.BlockedCountries {
			if country == userCountry {
				return fmt.Errorf("access blocked from country: %s", userCountry)
			}
		}
	}

	// Check time restrictions
	if policy.TimeRestrictions != nil {
		now := time.Now()
		if policy.TimeRestrictions.StartTime != nil && now.Before(*policy.TimeRestrictions.StartTime) {
			return fmt.Errorf("video not available yet")
		}
		if policy.TimeRestrictions.EndTime != nil && now.After(*policy.TimeRestrictions.EndTime) {
			return fmt.Errorf("video no longer available")
		}
	}

	return nil
}

// CreateViewingSession creates a secure viewing session
func (drm *DRMManager) CreateViewingSession(userID, videoID uuid.UUID, sessionDuration time.Duration) (*ViewingSession, error) {
	sessionID := uuid.New()

	session := &ViewingSession{
		SessionID:   sessionID,
		UserID:      userID,
		VideoID:     videoID,
		CreatedAt:   time.Now(),
		ExpiresAt:   time.Now().Add(sessionDuration),
		Status:      "active",
		ViewCount:   0,
		MaxViews:    1, // Default to single view
	}

	// Generate session token
	sessionToken, err := drm.generateSessionToken(session)
	if err != nil {
		return nil, fmt.Errorf("failed to generate session token: %w", err)
	}

	session.Token = sessionToken
	return session, nil
}

// ValidateViewingSession validates a viewing session
func (drm *DRMManager) ValidateViewingSession(sessionToken string) (*ViewingSession, error) {
	// This is a simplified version - in production, you'd store sessions in Redis/database
	// and validate against stored data

	// Decode and validate session token
	claimsJSON, err := base64.URLEncoding.DecodeString(sessionToken)
	if err != nil {
		return nil, fmt.Errorf("invalid session token")
	}

	var session ViewingSession
	if err := json.Unmarshal(claimsJSON, &session); err != nil {
		return nil, fmt.Errorf("failed to decode session")
	}

	// Check expiry
	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
	}

	// Check view limits
	if session.MaxViews > 0 && session.ViewCount >= session.MaxViews {
		return nil, fmt.Errorf("view limit exceeded")
	}

	return &session, nil
}

// Helper methods

func (drm *DRMManager) generateSignature(data string) string {
	h := hmac.New(sha256.New, []byte(drm.secretKey))
	h.Write([]byte(data))
	return base64.URLEncoding.EncodeToString(h.Sum(nil))
}

func (drm *DRMManager) generateSessionToken(session *ViewingSession) (string, error) {
	sessionJSON, err := json.Marshal(session)
	if err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(sessionJSON), nil
}

// ViewingSession represents a secure video viewing session
type ViewingSession struct {
	SessionID uuid.UUID `json:"session_id"`
	UserID    uuid.UUID `json:"user_id"`
	VideoID   uuid.UUID `json:"video_id"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
	Status    string    `json:"status"` // "active", "expired", "revoked"
	ViewCount int       `json:"view_count"`
	MaxViews  int       `json:"max_views"`
	Token     string    `json:"token"`
	IPAddress string    `json:"ip_address,omitempty"`
	UserAgent string    `json:"user_agent,omitempty"`
}

// Security logging for audit trails
type SecurityEvent struct {
	EventType   string    `json:"event_type"`
	UserID      uuid.UUID `json:"user_id"`
	VideoID     uuid.UUID `json:"video_id"`
	IPAddress   string    `json:"ip_address"`
	UserAgent   string    `json:"user_agent"`
	Timestamp   time.Time `json:"timestamp"`
	Description string    `json:"description"`
	Severity    string    `json:"severity"` // "low", "medium", "high", "critical"
}

func (drm *DRMManager) LogSecurityEvent(event *SecurityEvent) {
	// In production, this would send to a security logging service
	// For now, we'll just log to stdout
	eventJSON, _ := json.Marshal(event)
	fmt.Printf("SECURITY_EVENT: %s\n", eventJSON)
}