package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"video-service/internal/config"
)

type CloudflareService struct {
	config     *config.Config
	httpClient *http.Client
}

type DirectUploadResponse struct {
	Result struct {
		UploadURL string `json:"uploadURL"`
		UID       string `json:"uid"`
	} `json:"result"`
	Success  bool     `json:"success"`
	Errors   []string `json:"errors"`
	Messages []string `json:"messages"`
}

type StreamVideo struct {
	UID              string                 `json:"uid"`
	Status           map[string]interface{} `json:"status"`
	Meta             map[string]interface{} `json:"meta"`
	Preview          string                 `json:"preview"`
	Thumbnail        string                 `json:"thumbnail"`
	ReadyToStream    bool                   `json:"readyToStream"`
	Duration         float64                `json:"duration"`
	Input            map[string]interface{} `json:"input"`
	Playback         map[string]interface{} `json:"playback"`
	WaterMark        map[string]interface{} `json:"watermark"`
	Created          time.Time              `json:"created"`
	Modified         time.Time              `json:"modified"`
}

type StreamVideoResponse struct {
	Result  StreamVideo `json:"result"`
	Success bool        `json:"success"`
	Errors  []string    `json:"errors"`
}

func NewCloudflareService(cfg *config.Config) *CloudflareService {
	return &CloudflareService{
		config:     cfg,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// CreateDirectUploadURL creates a one-time upload URL for direct uploads
func (cs *CloudflareService) CreateDirectUploadURL(maxDurationSeconds int) (*DirectUploadResponse, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/stream/direct_upload", cs.config.CloudflareAccountID)
	
	payload := map[string]interface{}{
		"maxDurationSeconds": maxDurationSeconds,
	}
	
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cs.config.CloudflareStreamToken)
	
	resp, err := cs.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var response DirectUploadResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	if !response.Success {
		return nil, fmt.Errorf("cloudflare API error: %v", response.Errors)
	}
	
	return &response, nil
}

// GetVideoDetails retrieves video details from Cloudflare Stream
func (cs *CloudflareService) GetVideoDetails(uid string) (*StreamVideo, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/stream/%s", cs.config.CloudflareAccountID, uid)
	
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+cs.config.CloudflareStreamToken)
	
	resp, err := cs.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	var response StreamVideoResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}
	
	if !response.Success {
		return nil, fmt.Errorf("cloudflare API error: %v", response.Errors)
	}
	
	return &response.Result, nil
}

// CreateResumableUploadURL creates a resumable upload URL using TUS protocol
func (cs *CloudflareService) CreateResumableUploadURL(uploadLength int64, metadata map[string]string) (string, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/stream?direct_user=true", cs.config.CloudflareAccountID)
	
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+cs.config.CloudflareStreamToken)
	req.Header.Set("Tus-Resumable", "1.0.0")
	req.Header.Set("Upload-Length", strconv.FormatInt(uploadLength, 10))
	
	// Add metadata if provided
	if len(metadata) > 0 {
		var metadataPairs []string
		for key, value := range metadata {
			metadataPairs = append(metadataPairs, fmt.Sprintf("%s %s", key, value))
		}
		req.Header.Set("Upload-Metadata", strings.Join(metadataPairs, ","))
	}
	
	resp, err := cs.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	location := resp.Header.Get("Location")
	if location == "" {
		return "", fmt.Errorf("no location header in response")
	}
	
	return location, nil
}

// DeleteVideo deletes a video from Cloudflare Stream
func (cs *CloudflareService) DeleteVideo(uid string) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/stream/%s", cs.config.CloudflareAccountID, uid)
	
	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Authorization", "Bearer "+cs.config.CloudflareStreamToken)
	
	resp, err := cs.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	return nil
}

// UpdateVideoMetadata updates video metadata in Cloudflare Stream
func (cs *CloudflareService) UpdateVideoMetadata(uid string, metadata map[string]interface{}) error {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/stream/%s", cs.config.CloudflareAccountID, uid)
	
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return fmt.Errorf("failed to marshal metadata: %w", err)
	}
	
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cs.config.CloudflareStreamToken)
	
	resp, err := cs.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}
	
	return nil
}

// GetStreamURL generates the streaming URL for a video
func (cs *CloudflareService) GetStreamURL(uid string) string {
	return fmt.Sprintf("https://customer-mwulubub36waz34d.cloudflarestream.com/%s/manifest/video.m3u8", uid)
}

// GetThumbnailURL generates the thumbnail URL for a video
func (cs *CloudflareService) GetThumbnailURL(uid string) string {
	return fmt.Sprintf("https://customer-mwulubub36waz34d.cloudflarestream.com/%s/thumbnails/thumbnail.jpg", uid)
}

// GetEmbedURL generates the embed URL for a video
func (cs *CloudflareService) GetEmbedURL(uid string) string {
	return fmt.Sprintf("https://customer-mwulubub36waz34d.cloudflarestream.com/%s/iframe", uid)
}

// CreateDirectUploadURLWithExpiration creates a direct upload URL with custom expiration
func (cs *CloudflareService) CreateDirectUploadURLWithExpiration(maxDurationSeconds int) (*DirectUploadResponse, error) {
	url := fmt.Sprintf("https://api.cloudflare.com/client/v4/accounts/%s/stream/direct_upload", cs.config.CloudflareAccountID)

	payload := map[string]interface{}{
		"maxDurationSeconds": maxDurationSeconds,
		"requireSignedURLs": true, // Enable signed URLs for security
		"expiry": time.Now().Add(time.Duration(maxDurationSeconds) * time.Second).UTC().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+cs.config.CloudflareStreamToken)

	resp, err := cs.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var response DirectUploadResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if !response.Success {
		return nil, fmt.Errorf("cloudflare API error: %v", response.Errors)
	}

	return &response, nil
}