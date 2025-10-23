package handler

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"bucket-service/internal/config"
	"bucket-service/internal/middleware"
	"bucket-service/internal/model"
	"bucket-service/internal/repository"
	"bucket-service/internal/service"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type FileHandler struct {
	fileRepo     *repository.FileRepository
	s3Service    *service.S3Service
	imageService *service.ImageService
	config       *config.Config
}

func NewFileHandler(fileRepo *repository.FileRepository, s3Service *service.S3Service, imageService *service.ImageService, cfg *config.Config) *FileHandler {
	return &FileHandler{
		fileRepo:     fileRepo,
		s3Service:    s3Service,
		imageService: imageService,
		config:       cfg,
	}
}

type UploadResponse struct {
	FileID       uuid.UUID `json:"file_id"`
	Filename     string    `json:"filename"`
	Size         int64     `json:"size"`
	ContentType  string    `json:"content_type"`
	URL          string    `json:"url"`
	ThumbnailURL *string   `json:"thumbnail_url,omitempty"`
}

func (h *FileHandler) UploadFile(c *gin.Context) {
	fmt.Printf("DEBUG: Upload request received\n")
	userID, exists := middleware.GetUserID(c)
	if !exists {
		fmt.Printf("DEBUG: User not authenticated\n")
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	fmt.Printf("DEBUG: User authenticated: %s\n", userID.String())

	// Parse multipart form
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse file: %v\n", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}
	defer file.Close()
	fmt.Fprintf(os.Stderr, "File parsed successfully: %s\n", header.Filename)

	// Get optional parameters
	bucketType := c.PostForm("bucket")
	if bucketType == "" {
		bucketType = "general"
	}
	isPublic := c.PostForm("is_public") == "true"
	metadata := c.PostForm("metadata")
	customKey := c.PostForm("key") // Custom key for specific file paths like {userId}/avatar

	// Validate file
	if err := h.validateFile(header, file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Determine file type and content type
	contentType := h.s3Service.GetContentType(header.Filename)
	fileType := model.GetFileType(contentType)

	// Process image if it's an image file
	var processedData []byte
	var thumbnailURL *string

	if h.imageService != nil && h.imageService.IsImageFile(contentType) {
		// Process main image
		processed, err := h.imageService.ProcessImage(file, contentType)
		if err != nil {
			fmt.Printf("Image processing error: %v\n", err)
			// For now, skip image processing and treat as regular file
			processedData = make([]byte, header.Size)
			file.Seek(0, 0)
			file.Read(processedData)
		} else if processed != nil {
			processedData = processed.Data
			if processed.ContentType != "" {
				contentType = processed.ContentType
			}

			// Try to generate thumbnails - completely isolated to prevent any panics
			func() {
				defer func() {
					if r := recover(); r != nil {
						fmt.Printf("Thumbnail generation panic recovered: %v\n", r)
					}
				}()

				// Reset file pointer for thumbnail generation
				file.Seek(0, 0)
				thumbnails, err := h.imageService.GenerateThumbnails(file, contentType)
				if err != nil {
					fmt.Printf("Thumbnail generation error: %v\n", err)
					return
				}

				if thumbnails == nil || len(thumbnails) == 0 {
					fmt.Printf("No thumbnails generated\n")
					return
				}

				// Upload thumbnail (using medium size)
				thumbnail, exists := thumbnails["medium"]
				if !exists || thumbnail == nil || len(thumbnail.Data) == 0 {
					fmt.Printf("Medium thumbnail not available or empty\n")
					return
				}

				thumbReader := bytes.NewReader(thumbnail.Data)
				if thumbReader == nil {
					fmt.Printf("Failed to create thumbnail reader\n")
					return
				}

				thumbInput := &service.UploadInput{
					Reader:      thumbReader,
					Filename:    "thumb_" + header.Filename,
					ContentType: thumbnail.ContentType,
					FileType:    model.FileTypeImage,
					UserID:      userID,
					IsPublic:    isPublic,
				}

				thumbResult, err := h.s3Service.UploadFile(c.Request.Context(), thumbInput)
				if err != nil {
					fmt.Printf("Thumbnail upload error: %v\n", err)
					return
				}

				if thumbResult != nil && thumbResult.URL != "" {
					thumbnailURL = &thumbResult.URL
					fmt.Printf("Thumbnail uploaded successfully: %s\n", thumbResult.URL)
				}
			}()
		}
	} else {
		// Read file data for non-images
		processedData = make([]byte, header.Size)
		file.Seek(0, 0)
		file.Read(processedData)
	}

	// If custom key provided, check if file exists and delete it first
	if customKey != "" {
		fmt.Printf("DEBUG: Using custom key: %s\n", customKey)
		// Get the actual bucket name used by S3Service
		actualBucketName := h.s3Service.GetBucketName(fileType)
		// Try to find and delete existing file at this key
		existingFile, err := h.fileRepo.GetFileByObjectKey(c.Request.Context(), actualBucketName, customKey)
		if err == nil && existingFile != nil {
			fmt.Printf("DEBUG: Found existing file at key %s, deleting...\n", customKey)
			// Delete from S3
			h.s3Service.DeleteFile(c.Request.Context(), existingFile.BucketName, existingFile.ObjectKey)
			// Delete from database
			h.fileRepo.DeleteFile(c.Request.Context(), existingFile.ID)
		}
	}

	// Upload to S3
	uploadInput := &service.UploadInput{
		Reader:      bytes.NewReader(processedData),
		Filename:    header.Filename,
		ContentType: contentType,
		FileType:    fileType,
		UserID:      userID,
		IsPublic:    isPublic,
		CustomKey:   customKey, // Pass custom key to S3 service
		Metadata: map[string]string{
			"original_filename": header.Filename,
			"bucket_type":       bucketType,
			"metadata":          metadata,
		},
	}

	uploadResult, err := h.s3Service.UploadFile(c.Request.Context(), uploadInput)
	if err != nil {
		fmt.Printf("DEBUG: Upload error details: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file", "details": err.Error()})
		return
	}

	// Save file metadata to database
	fileModel := &model.File{
		ID:               uploadResult.FileID,
		Filename:         header.Filename,
		OriginalFilename: header.Filename,
		ContentType:      uploadResult.ContentType,
		SizeBytes:        uploadResult.Size,
		BucketName:       uploadResult.BucketName,
		ObjectKey:        uploadResult.ObjectKey,
		UploadUserID:     userID,
		IsPublic:         isPublic,
		Checksum:         &uploadResult.Checksum,
		ThumbnailURL:     thumbnailURL,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	if err := h.fileRepo.CreateFile(c.Request.Context(), fileModel); err != nil {
		// Cleanup S3 file on database error
		h.s3Service.DeleteFile(c.Request.Context(), uploadResult.BucketName, uploadResult.ObjectKey)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})
		return
	}

	response := &UploadResponse{
		FileID:       uploadResult.FileID,
		Filename:     header.Filename,
		Size:         uploadResult.Size,
		ContentType:  uploadResult.ContentType,
		URL:          uploadResult.URL,
		ThumbnailURL: thumbnailURL,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *FileHandler) GetFile(c *gin.Context) {
	fileIDStr := c.Param("fileId")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get file from database
	file, err := h.fileRepo.GetFileByID(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Check permissions
	if !file.IsPublic && file.UploadUserID != userID {
		hasPermission, err := h.fileRepo.HasFilePermission(c.Request.Context(), fileID, userID, "read")
		if err != nil || !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	// Generate presigned URL
	expires := time.Duration(h.config.Upload.PresignedURLExpiration) * time.Second
	url, err := h.s3Service.GetPresignedURL(c.Request.Context(), file.BucketName, file.ObjectKey, expires)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate download URL"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// GetPresignedURL returns a presigned URL as JSON instead of redirecting
func (h *FileHandler) GetPresignedURL(c *gin.Context) {
	type PresignedRequest struct {
		BucketName string `json:"bucket_name" binding:"required"`
		ObjectKey  string `json:"object_key" binding:"required"`
		ExpiresIn  int    `json:"expires_in"` // seconds, defaults to config value
		Operation  string `json:"operation"`  // GetObject, PutObject, defaults to GetObject
	}

	var req PresignedRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body", "details": err.Error()})
		return
	}

	// Default expiration from config
	expiresIn := req.ExpiresIn
	if expiresIn <= 0 {
		expiresIn = h.config.Upload.PresignedURLExpiration
	}
	expires := time.Duration(expiresIn) * time.Second

	// Default operation
	operation := req.Operation
	if operation == "" {
		operation = "GetObject"
	}

	var url string
	var err error

	switch operation {
	case "GetObject":
		url, err = h.s3Service.GetPresignedURL(c.Request.Context(), req.BucketName, req.ObjectKey, expires)
	case "PutObject":
		// For PutObject, we'd need content type, but for now just support GetObject
		c.JSON(http.StatusBadRequest, gin.H{"error": "PutObject operation not yet supported in this endpoint"})
		return
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported operation. Use 'GetObject'"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate presigned URL", "details": err.Error()})
		return
	}

	expiresAt := time.Now().Add(expires)
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": gin.H{
			"url":        url,
			"expires_at": expiresAt.Format(time.RFC3339),
		},
	})
}

func (h *FileHandler) GetFileMetadata(c *gin.Context) {
	fileIDStr := c.Param("fileId")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get file from database
	file, err := h.fileRepo.GetFileByID(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Check permissions
	if !file.IsPublic && file.UploadUserID != userID {
		hasPermission, err := h.fileRepo.HasFilePermission(c.Request.Context(), fileID, userID, "read")
		if err != nil || !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	c.JSON(http.StatusOK, file)
}

func (h *FileHandler) DeleteFile(c *gin.Context) {
	fileIDStr := c.Param("fileId")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get file from database
	file, err := h.fileRepo.GetFileByID(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Check permissions (only owner or admin can delete)
	userRole, _ := middleware.GetUserRole(c)
	if file.UploadUserID != userID && userRole != "admin" {
		hasPermission, err := h.fileRepo.HasFilePermission(c.Request.Context(), fileID, userID, "delete")
		if err != nil || !hasPermission {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access denied"})
			return
		}
	}

	// Delete from S3
	if err := h.s3Service.DeleteFile(c.Request.Context(), file.BucketName, file.ObjectKey); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file from storage"})
		return
	}

	// Soft delete from database
	if err := h.fileRepo.SoftDeleteFile(c.Request.Context(), fileID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete file record"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "File deleted successfully"})
}

func (h *FileHandler) ListFiles(c *gin.Context) {
	userID, exists := middleware.GetUserID(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Parse query parameters
	options := &repository.ListFilesOptions{
		UserID:      &userID,
		BucketName:  c.Query("bucket"),
		ContentType: c.Query("content_type"),
		Search:      c.Query("search"),
		SortBy:      c.Query("sort"),
		SortOrder:   c.Query("order"),
	}

	if isPublicStr := c.Query("is_public"); isPublicStr != "" {
		isPublic := isPublicStr == "true"
		options.IsPublic = &isPublic
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil {
			options.Page = page
		}
	}

	if limitStr := c.Query("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil {
			options.Limit = limit
		}
	}

	result, err := h.fileRepo.ListFiles(c.Request.Context(), options)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list files"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *FileHandler) GetPublicFile(c *gin.Context) {
	fileIDStr := c.Param("fileId")
	fileID, err := uuid.Parse(fileIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid file ID"})
		return
	}

	// Get file from database
	file, err := h.fileRepo.GetFileByID(c.Request.Context(), fileID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "File not found"})
		return
	}

	// Check if file is public
	if !file.IsPublic {
		c.JSON(http.StatusForbidden, gin.H{"error": "File is not public"})
		return
	}

	// Generate presigned URL
	expires := time.Duration(h.config.Upload.PresignedURLExpiration) * time.Second
	url, err := h.s3Service.GetPresignedURL(c.Request.Context(), file.BucketName, file.ObjectKey, expires)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate download URL"})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Helper methods
func (h *FileHandler) validateFile(header *multipart.FileHeader, file multipart.File) error {
	// Check for empty file
	if header.Size == 0 {
		return fmt.Errorf("file is empty, cannot upload zero-byte file")
	}

	// Check file size
	if header.Size > h.config.Upload.MaxFileSize {
		return fmt.Errorf("file size exceeds limit of %d bytes", h.config.Upload.MaxFileSize)
	}

	// Check file extension
	ext := strings.ToLower(filepath.Ext(header.Filename)[1:]) // Remove the dot
	if ext == "" {
		return fmt.Errorf("file must have an extension")
	}

	// Check if file type is allowed
	allowedTypes := append(h.config.Upload.AllowedDocumentTypes, h.config.Upload.AllowedImageTypes...)
	allowedTypes = append(allowedTypes, h.config.Upload.AllowedArchiveTypes...)
	
	allowed := false
	for _, allowedType := range allowedTypes {
		if ext == allowedType {
			allowed = true
			break
		}
	}

	if !allowed {
		return fmt.Errorf("file type '%s' not allowed", ext)
	}

	return nil
}