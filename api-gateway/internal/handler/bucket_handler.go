package handler

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/study-platform/pkg/logger"
)

type BucketHandler struct {
	bucketServiceURL string
	logger           logger.Logger
	client           *http.Client
}

func NewBucketHandler(bucketServiceURL string, log logger.Logger) *BucketHandler {
	return &BucketHandler{
		bucketServiceURL: bucketServiceURL,
		logger:           log,
		client:           &http.Client{},
	}
}

// UploadFile handles file upload
func (bh *BucketHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	bh.logger.Info("Handling file upload request")

	// Create request to bucket service
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/api/files/upload", bh.bucketServiceURL), r.Body)
	if err != nil {
		bh.logger.Error(fmt.Errorf("failed to create request: %w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Copy headers
	req.Header.Set("Content-Type", r.Header.Get("Content-Type"))
	req.Header.Set("Authorization", r.Header.Get("Authorization"))
	req.ContentLength = r.ContentLength

	// Forward request to bucket service
	bh.forwardRequest(w, req)
}

// DownloadFile handles file download
func (bh *BucketHandler) DownloadFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["fileId"]

	bh.logger.Info(fmt.Sprintf("Handling file download request for file %s", fileID))

	// Create request to bucket service
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/api/files/%s", bh.bucketServiceURL, fileID), nil)
	if err != nil {
		bh.logger.Error(fmt.Errorf("failed to create request: %w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Copy authorization header
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	// Forward request to bucket service
	bh.forwardRequest(w, req)
}

// DeleteFile handles file deletion
func (bh *BucketHandler) DeleteFile(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["fileId"]

	bh.logger.Info(fmt.Sprintf("Handling file deletion request for file %s", fileID))

	// Create request to bucket service
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/api/files/%s", bh.bucketServiceURL, fileID), nil)
	if err != nil {
		bh.logger.Error(fmt.Errorf("failed to create request: %w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Copy authorization header
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	// Forward request to bucket service
	bh.forwardRequest(w, req)
}

// ListFiles handles file listing
func (bh *BucketHandler) ListFiles(w http.ResponseWriter, r *http.Request) {
	bh.logger.Info("Handling file list request")

	// Build query parameters
	params := url.Values{}
	for key, values := range r.URL.Query() {
		for _, value := range values {
			params.Add(key, value)
		}
	}

	// Create request to bucket service
	reqURL := fmt.Sprintf("http://%s/api/files?%s", bh.bucketServiceURL, params.Encode())
	req, err := http.NewRequest("GET", reqURL, nil)
	if err != nil {
		bh.logger.Error(fmt.Errorf("failed to create request: %w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Copy authorization header
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	// Forward request to bucket service
	bh.forwardRequest(w, req)
}

// GetFileMetadata handles file metadata retrieval
func (bh *BucketHandler) GetFileMetadata(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fileID := vars["fileId"]

	bh.logger.Info(fmt.Sprintf("Handling file metadata request for file %s", fileID))

	// Create request to bucket service
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/api/files/%s/metadata", bh.bucketServiceURL, fileID), nil)
	if err != nil {
		bh.logger.Error(fmt.Errorf("failed to create request: %w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Copy authorization header
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	// Forward request to bucket service
	bh.forwardRequest(w, req)
}

// StartMultipartUpload handles multipart upload initialization
func (bh *BucketHandler) StartMultipartUpload(w http.ResponseWriter, r *http.Request) {
	bh.logger.Info("Handling multipart upload start request")

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		bh.logger.Error(fmt.Errorf("failed to read request body: %w", err))
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Create request to bucket service
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/api/files/upload/start", bh.bucketServiceURL), bytes.NewReader(body))
	if err != nil {
		bh.logger.Error(fmt.Errorf("failed to create request: %w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Copy headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	// Forward request to bucket service
	bh.forwardRequest(w, req)
}

// CompleteMultipartUpload handles multipart upload completion
func (bh *BucketHandler) CompleteMultipartUpload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	bh.logger.Info(fmt.Sprintf("Handling multipart upload completion for session %s", sessionID))

	// Read request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		bh.logger.Error(fmt.Errorf("failed to read request body: %w", err))
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Create request to bucket service
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s/api/files/upload/%s/complete", bh.bucketServiceURL, sessionID), bytes.NewReader(body))
	if err != nil {
		bh.logger.Error(fmt.Errorf("failed to create request: %w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Copy headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	// Forward request to bucket service
	bh.forwardRequest(w, req)
}

// AbortMultipartUpload handles multipart upload abortion
func (bh *BucketHandler) AbortMultipartUpload(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	bh.logger.Info(fmt.Sprintf("Handling multipart upload abort for session %s", sessionID))

	// Create request to bucket service
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://%s/api/files/upload/%s", bh.bucketServiceURL, sessionID), nil)
	if err != nil {
		bh.logger.Error(fmt.Errorf("failed to create request: %w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Copy authorization header
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	// Forward request to bucket service
	bh.forwardRequest(w, req)
}

// GetUploadProgress handles upload progress tracking
func (bh *BucketHandler) GetUploadProgress(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	sessionID := vars["sessionId"]

	bh.logger.Info(fmt.Sprintf("Handling upload progress request for session %s", sessionID))

	// Create request to bucket service
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s/api/files/upload/%s/progress", bh.bucketServiceURL, sessionID), nil)
	if err != nil {
		bh.logger.Error(fmt.Errorf("failed to create request: %w", err))
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Copy authorization header
	req.Header.Set("Authorization", r.Header.Get("Authorization"))

	// Forward request to bucket service
	bh.forwardRequest(w, req)
}

// forwardRequest forwards the HTTP request to bucket service
func (bh *BucketHandler) forwardRequest(w http.ResponseWriter, req *http.Request) {
	// Execute request
	resp, err := bh.client.Do(req)
	if err != nil {
		bh.logger.Error(fmt.Errorf("failed to forward request to bucket service: %w", err))
		http.Error(w, "Service unavailable", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	// Copy response headers
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// Set status code
	w.WriteHeader(resp.StatusCode)

	// Copy response body
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		bh.logger.Error(fmt.Errorf("failed to copy response body: %w", err))
	}
}