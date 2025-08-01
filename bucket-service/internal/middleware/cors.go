package middleware

import (
	"bucket-service/internal/config"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS(cfg *config.Config) gin.HandlerFunc {
	corsConfig := cors.Config{
		AllowOrigins:     cfg.Security.CORSOrigins,
		AllowMethods:     cfg.Security.CORSMethods,
		AllowHeaders:     cfg.Security.CORSHeaders,
		ExposeHeaders:    []string{"Content-Length", "Content-Type", "ETag", "x-amz-version-id"},
		AllowCredentials: true,
		MaxAge:           12 * 3600, // 12 hours
	}

	return cors.New(corsConfig)
}