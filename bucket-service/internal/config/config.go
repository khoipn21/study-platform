package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	S3       S3Config
	Upload   UploadConfig
	Security SecurityConfig
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type S3Config struct {
	Endpoint        string
	Region          string
	AccessKeyID     string
	SecretAccessKey string
	BucketPrefix    string
	
	// Bucket Names
	BucketFiles          string
	BucketDocuments      string
	BucketImages         string
	BucketAvatars        string
	BucketCourseMaterials string
}

type UploadConfig struct {
	MaxFileSize           int64
	MultipartThreshold    int64
	PartSize              int64
	MaxUploadParts        int
	PresignedURLExpiration int
	ChunkSize             int64
	
	// Allowed file types
	AllowedDocumentTypes []string
	AllowedImageTypes    []string
	AllowedArchiveTypes  []string
}

type SecurityConfig struct {
	JWTSecret   string
	CORSOrigins []string
	CORSMethods []string
	CORSHeaders []string
}

type LoggingConfig struct {
	Level  string
	Format string
	File   string
}

func Load() (*Config, error) {
	config := &Config{
		Server: ServerConfig{
			Port: getEnv("BUCKET_SERVICE_PORT", "8084"),
			Host: getEnv("BUCKET_SERVICE_HOST", "0.0.0.0"),
			Env:  getEnv("ENVIRONMENT", "development"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "admin"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "studyplatform"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		S3: S3Config{
			Endpoint:        getEnv("S3_ENDPOINT", "https://s3.amazonaws.com"),
			Region:          getEnv("S3_REGION", "us-east-1"),
			AccessKeyID:     getEnv("S3_ACCESS_KEY_ID", ""),
			SecretAccessKey: getEnv("S3_SECRET_ACCESS_KEY", ""),
			BucketPrefix:    getEnv("S3_BUCKET_PREFIX", "study-platform"),
			
			BucketFiles:          getEnv("S3_BUCKET_FILES", "study-platform-files-khoipn-dev"),
			BucketDocuments:      getEnv("S3_BUCKET_DOCUMENTS", "study-platform-documents-khoipn-dev"),
			BucketImages:         getEnv("S3_BUCKET_IMAGES", "study-platform-images-khoipn-dev"),
			BucketAvatars:        getEnv("S3_BUCKET_AVATARS", "study-platform-avatars-khoipn-dev"),
			BucketCourseMaterials: getEnv("S3_BUCKET_COURSE_MATERIALS", "study-platform-course-materials-khoipn-dev"),
		},
		Upload: UploadConfig{
			MaxFileSize:           getEnvInt64("MAX_FILE_SIZE", 1073741824), // 1GB
			MultipartThreshold:    getEnvInt64("MULTIPART_THRESHOLD", 104857600), // 100MB
			PartSize:              getEnvInt64("PART_SIZE", 5242880), // 5MB
			MaxUploadParts:        getEnvInt("MAX_UPLOAD_PARTS", 1000),
			PresignedURLExpiration: getEnvInt("PRESIGNED_URL_EXPIRATION", 3600), // 1 hour
			ChunkSize:             getEnvInt64("CHUNK_SIZE", 1048576), // 1MB
			
			AllowedDocumentTypes: getEnvStringSlice("ALLOWED_DOCUMENT_TYPES", []string{
				"pdf", "doc", "docx", "ppt", "pptx", "xls", "xlsx", "txt", "rtf",
			}),
			AllowedImageTypes: getEnvStringSlice("ALLOWED_IMAGE_TYPES", []string{
				"jpg", "jpeg", "png", "gif", "webp", "svg", "bmp",
			}),
			AllowedArchiveTypes: getEnvStringSlice("ALLOWED_ARCHIVE_TYPES", []string{
				"zip", "rar", "7z", "tar", "gz",
			}),
		},
		Security: SecurityConfig{
			JWTSecret: getEnv("JWT_SECRET", "your-super-secret-jwt-key-min-256-bits"),
			CORSOrigins: getEnvStringSlice("CORS_ORIGINS", []string{
				"http://localhost:3000", "https://yourdomain.com",
			}),
			CORSMethods: getEnvStringSlice("CORS_METHODS", []string{
				"GET", "POST", "PUT", "DELETE", "OPTIONS",
			}),
			CORSHeaders: getEnvStringSlice("CORS_HEADERS", []string{
				"Origin", "Content-Type", "Accept", "Authorization",
			}),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
			File:   getEnv("LOG_FILE", ""),
		},
	}

	// Validate required configuration
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	return config, nil
}

func (c *Config) Validate() error {
	if c.S3.AccessKeyID == "" {
		return fmt.Errorf("S3_ACCESS_KEY_ID is required")
	}
	if c.S3.SecretAccessKey == "" {
		return fmt.Errorf("S3_SECRET_ACCESS_KEY is required")
	}
	if c.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	if c.Security.JWTSecret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	return nil
}

func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		c.Database.Host, c.Database.User, c.Database.Password,
		c.Database.Name, c.Database.Port, c.Database.SSLMode,
	)
}

func (c *Config) GetBucketForType(fileType string) string {
	switch fileType {
	case "document":
		return c.S3.BucketDocuments
	case "image":
		return c.S3.BucketImages
	case "archive":
		return c.S3.BucketCourseMaterials
	default:
		return c.S3.BucketFiles
	}
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvInt64(key string, defaultValue int64) int64 {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.ParseInt(value, 10, 64); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvStringSlice(key string, defaultValue []string) []string {
	if value := os.Getenv(key); value != "" {
		return strings.Split(value, ",")
	}
	return defaultValue
}