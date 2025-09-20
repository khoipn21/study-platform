package service

import (
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"mime"
	"path/filepath"
	"strings"
	"time"

	"bucket-service/internal/config"
	"bucket-service/internal/model"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type S3Service struct {
	client   *s3.Client
	uploader *manager.Uploader
	config   *config.Config
}

func NewS3Service(cfg *config.Config) (*S3Service, error) {
	awsCfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithRegion(cfg.S3.Region),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.S3.AccessKeyID,
			cfg.S3.SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load AWS config: %w", err)
	}

	// Create S3 client with custom endpoint if provided
	var s3Client *s3.Client
	if cfg.S3.Endpoint != "https://s3.amazonaws.com" {
		s3Client = s3.NewFromConfig(awsCfg, func(o *s3.Options) {
			o.BaseEndpoint = aws.String(cfg.S3.Endpoint)
			o.UsePathStyle = true
		})
	} else {
		s3Client = s3.NewFromConfig(awsCfg)
	}

	uploader := manager.NewUploader(s3Client, func(u *manager.Uploader) {
		u.PartSize = cfg.Upload.PartSize
	})

	service := &S3Service{
		client:   s3Client,
		uploader: uploader,
		config:   cfg,
	}

	// Validate S3 connectivity during initialization (non-blocking in development)
	if err := service.ValidateConnection(context.TODO()); err != nil {
		// In development, just log the warning but don't fail
		// In production, you might want to fail fast
		fmt.Printf("Warning: S3 connection validation failed: %v\n", err)
		// Uncomment the next line for stricter validation in production
		// return nil, fmt.Errorf("S3 connection validation failed: %w", err)
	}

	return service, nil
}

type UploadInput struct {
	Reader      io.Reader
	Filename    string
	ContentType string
	FileType    model.FileType
	UserID      uuid.UUID
	IsPublic    bool
	Metadata    map[string]string
}

type UploadResult struct {
	FileID      uuid.UUID
	ObjectKey   string
	BucketName  string
	URL         string
	Size        int64
	Checksum    string
	ContentType string
}

func (s *S3Service) UploadFile(ctx context.Context, input *UploadInput) (*UploadResult, error) {
	// Generate unique object key
	fileID := uuid.New()
	objectKey := s.generateObjectKey(fileID, input.Filename, input.FileType)
	bucketName := s.GetBucketName(input.FileType)

	// Create metadata
	metadata := map[string]string{
		"user-id":     input.UserID.String(),
		"file-id":     fileID.String(),
		"filename":    input.Filename,
		"upload-time": time.Now().UTC().Format(time.RFC3339),
	}
	for k, v := range input.Metadata {
		metadata[k] = v
	}

	// Calculate size and checksum
	var size int64
	var checksum string
	if sizeReader, ok := input.Reader.(*SizeReader); ok {
		size = sizeReader.Size()
		checksum = sizeReader.Checksum()
	} else {
		// For readers without size info, we need to buffer
		data, err := io.ReadAll(input.Reader)
		if err != nil {
			return nil, fmt.Errorf("failed to read data: %w", err)
		}
		size = int64(len(data))
		checksum = fmt.Sprintf("%x", md5.Sum(data))
		input.Reader = strings.NewReader(string(data))
	}

	// Set ACL based on public flag
	var acl types.ObjectCannedACL
	if input.IsPublic {
		acl = types.ObjectCannedACLPublicRead
	} else {
		acl = types.ObjectCannedACLPrivate
	}

	// Upload to S3
	result, err := s.uploader.Upload(ctx, &s3.PutObjectInput{
		Bucket:        aws.String(bucketName),
		Key:           aws.String(objectKey),
		Body:          input.Reader,
		ContentType:   aws.String(input.ContentType),
		Metadata:      metadata,
		ACL:           acl,
		ContentLength: aws.Int64(size),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to upload to S3: %w", err)
	}

	// Construct the URL manually if Location is empty
	url := result.Location
	if url == "" {
		if s.config.S3.Endpoint == "https://s3.amazonaws.com" {
			url = fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, s.config.S3.Region, objectKey)
		} else {
			// For MinIO or custom endpoints
			url = fmt.Sprintf("%s/%s/%s", s.config.S3.Endpoint, bucketName, objectKey)
		}
	}

	return &UploadResult{
		FileID:      fileID,
		ObjectKey:   objectKey,
		BucketName:  bucketName,
		URL:         url,
		Size:        size,
		Checksum:    checksum,
		ContentType: input.ContentType,
	}, nil
}

func (s *S3Service) GetPresignedURL(ctx context.Context, bucketName, objectKey string, expires time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	request, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	return request.URL, nil
}

func (s *S3Service) GetPresignedUploadURL(ctx context.Context, bucketName, objectKey, contentType string, expires time.Duration) (string, error) {
	presignClient := s3.NewPresignClient(s.client)

	request, err := presignClient.PresignPutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
	}, func(opts *s3.PresignOptions) {
		opts.Expires = expires
	})
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload URL: %w", err)
	}

	return request.URL, nil
}

func (s *S3Service) DeleteFile(ctx context.Context, bucketName, objectKey string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}
	return nil
}

func (s *S3Service) HeadObject(ctx context.Context, bucketName, objectKey string) (*s3.HeadObjectOutput, error) {
	result, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get object metadata: %w", err)
	}
	return result, nil
}

// Multipart upload support
type MultipartUpload struct {
	UploadID   string
	BucketName string
	ObjectKey  string
	Parts      []types.CompletedPart
}

func (s *S3Service) CreateMultipartUpload(ctx context.Context, bucketName, objectKey, contentType string, metadata map[string]string) (*MultipartUpload, error) {
	result, err := s.client.CreateMultipartUpload(ctx, &s3.CreateMultipartUploadInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(objectKey),
		ContentType: aws.String(contentType),
		Metadata:    metadata,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create multipart upload: %w", err)
	}

	return &MultipartUpload{
		UploadID:   *result.UploadId,
		BucketName: bucketName,
		ObjectKey:  objectKey,
		Parts:      make([]types.CompletedPart, 0),
	}, nil
}

func (s *S3Service) GetMultipartPresignedURLs(ctx context.Context, upload *MultipartUpload, partNumbers []int32, expires time.Duration) ([]string, error) {
	presignClient := s3.NewPresignClient(s.client)
	urls := make([]string, len(partNumbers))

	for i, partNumber := range partNumbers {
		request, err := presignClient.PresignUploadPart(ctx, &s3.UploadPartInput{
			Bucket:     aws.String(upload.BucketName),
			Key:        aws.String(upload.ObjectKey),
			UploadId:   aws.String(upload.UploadID),
			PartNumber: aws.Int32(partNumber),
		}, func(opts *s3.PresignOptions) {
			opts.Expires = expires
		})
		if err != nil {
			return nil, fmt.Errorf("failed to generate presigned URL for part %d: %w", partNumber, err)
		}
		urls[i] = request.URL
	}

	return urls, nil
}

func (s *S3Service) CompleteMultipartUpload(ctx context.Context, upload *MultipartUpload) (*s3.CompleteMultipartUploadOutput, error) {
	result, err := s.client.CompleteMultipartUpload(ctx, &s3.CompleteMultipartUploadInput{
		Bucket:   aws.String(upload.BucketName),
		Key:      aws.String(upload.ObjectKey),
		UploadId: aws.String(upload.UploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: upload.Parts,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to complete multipart upload: %w", err)
	}
	return result, nil
}

func (s *S3Service) AbortMultipartUpload(ctx context.Context, upload *MultipartUpload) error {
	_, err := s.client.AbortMultipartUpload(ctx, &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(upload.BucketName),
		Key:      aws.String(upload.ObjectKey),
		UploadId: aws.String(upload.UploadID),
	})
	if err != nil {
		return fmt.Errorf("failed to abort multipart upload: %w", err)
	}
	return nil
}

// Helper methods
func (s *S3Service) generateObjectKey(fileID uuid.UUID, filename string, fileType model.FileType) string {
	ext := filepath.Ext(filename)
	timestamp := time.Now().Format("2006/01/02")
	return fmt.Sprintf("%s/%s/%s%s", fileType, timestamp, fileID.String(), ext)
}

func (s *S3Service) GetBucketName(fileType model.FileType) string {
	switch fileType {
	case model.FileTypeDocument:
		return s.config.S3.BucketDocuments
	case model.FileTypeImage:
		return s.config.S3.BucketImages
	case model.FileTypeArchive:
		return s.config.S3.BucketCourseMaterials
	default:
		return s.config.S3.BucketFiles
	}
}

func (s *S3Service) GetContentType(filename string) string {
	ext := filepath.Ext(filename)
	contentType := mime.TypeByExtension(ext)
	if contentType == "" {
		return "application/octet-stream"
	}
	return contentType
}

// SizeReader wraps an io.Reader to track size and calculate checksum
type SizeReader struct {
	reader   io.Reader
	size     int64
	checksum string
}

func NewSizeReader(reader io.Reader, size int64) *SizeReader {
	return &SizeReader{
		reader: reader,
		size:   size,
	}
}

func (sr *SizeReader) Read(p []byte) (n int, err error) {
	return sr.reader.Read(p)
}

func (sr *SizeReader) Size() int64 {
	return sr.size
}

func (sr *SizeReader) Checksum() string {
	if sr.checksum == "" {
		// Calculate checksum if not already done
		hash := md5.New()
		io.Copy(hash, sr.reader)
		sr.checksum = fmt.Sprintf("%x", hash.Sum(nil))
	}
	return sr.checksum
}

// ValidateConnection tests S3 connectivity by checking if we can access a specific bucket
func (s *S3Service) ValidateConnection(ctx context.Context) error {
	// Test basic connectivity by checking if we can head a bucket
	// This requires less permissions than listing all buckets
	bucketName := s.config.S3.BucketFiles
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		// Try to create the bucket if it doesn't exist (for development)
		_, createErr := s.client.CreateBucket(ctx, &s3.CreateBucketInput{
			Bucket: aws.String(bucketName),
		})
		if createErr != nil {
			return fmt.Errorf("cannot connect to S3 or create bucket '%s': %w", bucketName, err)
		}
	}
	return nil
}