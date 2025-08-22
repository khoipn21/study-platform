package config

import (
	"log"
	"os"
)

type Config struct {
	// Service Configuration
	Port string
	Host string

	// Database Configuration
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	// Cloudflare Stream Configuration
	CloudflareAccountID   string
	CloudflareStreamToken string
	CloudflareAPIEmail    string
	CloudflareAPIKey      string

	// Redis Configuration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
	RedisPoolSize int

	// WebSocket Configuration
	WSReadBufferSize    int
	WSWriteBufferSize   int
	WSHeartbeatInterval string

	// Video Processing
	MaxVideoSize     int64
	AllowedFormats   []string
	DefaultThumbnailTime int

	// Network Intelligence
	BandwidthCheckInterval   string
	QualityChangeThreshold   float64
	BufferTargetSeconds      int
	MinBufferSeconds         int

	// Security
	JWTSecret           string
	MaxSessionDuration  string
	CORSOrigins         []string

	// Analytics
	AnalyticsBatchSize     int
	AnalyticsFlushInterval string
}

func Load() *Config {
	config := &Config{
		// Service Configuration
		Port: getEnv("VIDEO_SERVICE_PORT", "8084"),
		Host: getEnv("VIDEO_SERVICE_HOST", "0.0.0.0"),

		// Database Configuration
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "2345"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "studyplatform"),

		// Cloudflare Stream Configuration
		CloudflareAccountID:   getEnv("CLOUDFLARE_ACCOUNT_ID", "db8adad78e6e5907f175d8048cc5391a"),
		CloudflareStreamToken: getEnv("CLOUDFLARE_STREAM_TOKEN", "9KJ0x_M53_acH1NTxuyAw2QE5x_RkzK754cZ5OSD"),
		CloudflareAPIEmail:    getEnv("CLOUDFLARE_API_EMAIL", ""),
		CloudflareAPIKey:      getEnv("CLOUDFLARE_API_KEY", ""),

		// Redis Configuration
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       0,
		RedisPoolSize: 10,

		// WebSocket Configuration
		WSReadBufferSize:    1024,
		WSWriteBufferSize:   1024,
		WSHeartbeatInterval: getEnv("WS_HEARTBEAT_INTERVAL", "30s"),

		// Video Processing
		MaxVideoSize:         5368709120, // 5GB
		AllowedFormats:       []string{"mp4", "avi", "mov", "wmv", "flv", "webm"},
		DefaultThumbnailTime: 10, // seconds

		// Network Intelligence
		BandwidthCheckInterval:   getEnv("BANDWIDTH_CHECK_INTERVAL", "10s"),
		QualityChangeThreshold:   0.8,
		BufferTargetSeconds:      10,
		MinBufferSeconds:         5,

		// Security
		JWTSecret:          getEnv("JWT_SECRET", "your-secret-key"),
		MaxSessionDuration: getEnv("MAX_SESSION_DURATION", "24h"),
		CORSOrigins:        []string{"http://localhost:3000", "https://yourdomain.com"},

		// Analytics
		AnalyticsBatchSize:     100,
		AnalyticsFlushInterval: getEnv("ANALYTICS_FLUSH_INTERVAL", "60s"),
	}

	// Validate required configuration
	if config.CloudflareAccountID == "" || config.CloudflareStreamToken == "" {
		log.Fatal("Cloudflare Stream configuration is required")
	}

	return config
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}