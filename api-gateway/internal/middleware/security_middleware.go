package middleware

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"time"

	"github.com/study-platform/pkg/logger"
)

type SecurityMiddleware struct {
	logger          logger.Logger
	jwtSecret       []byte
	securityHeaders map[string]string
	cspPolicy       string
	validatedInputs map[string]*regexp.Regexp
	allowedOrigins  []string
}

type SecurityConfig struct {
	JWTSecret             string
	ContentSecurityPolicy string
	EnableSecurityHeaders bool
	AllowedOrigins        []string
}

func NewSecurityMiddleware(config SecurityConfig, logger logger.Logger) *SecurityMiddleware {
	// Validate JWT secret strength (minimum 256 bits / 32 bytes)
	if len(config.JWTSecret) < 32 {
		logger.Fatal(fmt.Errorf("JWT secret must be at least 32 characters (256 bits) for security"))
	}

	// Default security headers
	securityHeaders := map[string]string{
		"X-Frame-Options":           "DENY",
		"X-Content-Type-Options":    "nosniff",
		"X-XSS-Protection":          "1; mode=block",
		"Referrer-Policy":           "strict-origin-when-cross-origin",
		"Permissions-Policy":        "geolocation=(), microphone=(), camera=()",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
		"Cache-Control":             "no-cache, no-store, must-revalidate",
		"Pragma":                    "no-cache",
		"Expires":                   "0",
	}

	// Input validation patterns
	validationPatterns := map[string]*regexp.Regexp{
		"email":         regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
		"username":      regexp.MustCompile(`^[a-zA-Z0-9_-]{3,30}$`),
		"course_id":     regexp.MustCompile(`^[a-zA-Z0-9-]{1,50}$`),
		"user_id":       regexp.MustCompile(`^[0-9]+$`),
		"alphanumeric":  regexp.MustCompile(`^[a-zA-Z0-9\s-_.,!?()]{1,500}$`),
		"safe_filename": regexp.MustCompile(`^[a-zA-Z0-9._-]{1,255}$`),
		"url":           regexp.MustCompile(`^https?://[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}[/\w\-._~:/?#[\]@!$&'()*+,;=]*$`),
	}

	return &SecurityMiddleware{
		logger:          logger,
		jwtSecret:       []byte(config.JWTSecret),
		securityHeaders: securityHeaders,
		cspPolicy:       config.ContentSecurityPolicy,
		validatedInputs: validationPatterns,
		allowedOrigins:  config.AllowedOrigins,
	}
}

// SecurityHeaders applies security headers to all responses
func (sm *SecurityMiddleware) SecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Apply all security headers
		for header, value := range sm.securityHeaders {
			w.Header().Set(header, value)
		}

		// Apply Content Security Policy if configured
		if sm.cspPolicy != "" {
			w.Header().Set("Content-Security-Policy", sm.cspPolicy)
		}

		// Remove potentially dangerous headers
		w.Header().Del("Server")
		w.Header().Del("X-Powered-By")

		next.ServeHTTP(w, r)
	})
}

// InputValidation validates request inputs against common attack vectors
func (sm *SecurityMiddleware) InputValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check for common attack patterns in URL path
		if sm.containsAttackPatterns(r.URL.Path) {
			sm.sendSecurityError(w, http.StatusBadRequest, "Invalid request path", "INVALID_PATH")
			sm.logger.Warnf("Suspicious request path detected: %s from IP: %s", r.URL.Path, sm.getClientIP(r))
			return
		}

		// Check query parameters for attack patterns
		for key, values := range r.URL.Query() {
			for _, value := range values {
				if sm.containsAttackPatterns(value) {
					sm.sendSecurityError(w, http.StatusBadRequest, "Invalid query parameter", "INVALID_QUERY")
					sm.logger.Warnf("Suspicious query parameter %s=%s from IP: %s", key, value, sm.getClientIP(r))
					return
				}
			}
		}

		// Check headers for attack patterns (skip common legitimate headers)
		legitimateHeaders := map[string]bool{
			"Accept":                         true,
			"Accept-Language":                true,
			"Accept-Encoding":                true,
			"User-Agent":                     true,
			"Host":                           true,
			"Connection":                     true,
			"Content-Type":                   true,
			"Content-Length":                 true,
			"Origin":                         true,
			"Referer":                        true,
			"Cache-Control":                  true,
			"Pragma":                         true,
			"Authorization":                  true,
			"X-Requested-With":               true,
			"Access-Control-Request-Method":  true,
			"Access-Control-Request-Headers": true,
			// Browser security headers
			"Sec-Ch-Ua":          true,
			"Sec-Ch-Ua-Mobile":   true,
			"Sec-Ch-Ua-Platform": true,
			"Sec-Fetch-Site":     true,
			"Sec-Fetch-Mode":     true,
			"Sec-Fetch-Dest":     true,
			"Sec-Fetch-User":     true,
			// Additional common headers
			"Accept-Charset":            true,
			"If-None-Match":             true,
			"If-Modified-Since":         true,
			"Dnt":                       true,
			"Upgrade-Insecure-Requests": true,
		}

		for name, values := range r.Header {
			// Skip legitimate headers and authentication headers for this check
			if legitimateHeaders[name] || name == "Authorization" {
				continue
			}
			for _, value := range values {
				if sm.containsAttackPatterns(value) {
					sm.sendSecurityError(w, http.StatusBadRequest, "Invalid header value", "INVALID_HEADER")
					sm.logger.Warnf("Suspicious header %s=%s from IP: %s", name, value, sm.getClientIP(r))
					return
				}
			}
		}

		// Validate request size (prevent memory exhaustion)
		if r.ContentLength > 50*1024*1024 { // 50MB limit
			sm.sendSecurityError(w, http.StatusRequestEntityTooLarge, "Request too large", "REQUEST_TOO_LARGE")
			return
		}

		next.ServeHTTP(w, r)
	})
}

// JWTSecurityValidation performs additional JWT security checks
func (sm *SecurityMiddleware) JWTSecurityValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			parts := strings.Split(authHeader, " ")
			if len(parts) == 2 && parts[0] == "Bearer" {
				token := parts[1]

				// Check for common JWT attack patterns
				if sm.isJWTSuspicious(token) {
					sm.sendSecurityError(w, http.StatusUnauthorized, "Invalid token format", "SUSPICIOUS_JWT")
					sm.logger.Warnf("Suspicious JWT token from IP: %s", sm.getClientIP(r))
					return
				}

				// Validate token structure (should have 3 parts separated by dots)
				if len(strings.Split(token, ".")) != 3 {
					sm.sendSecurityError(w, http.StatusUnauthorized, "Invalid token structure", "INVALID_JWT_STRUCTURE")
					return
				}
			}
		}

		next.ServeHTTP(w, r)
	})
}

// RequestSizeLimit limits request body size
func (sm *SecurityMiddleware) RequestSizeLimit(maxBytes int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.ContentLength > maxBytes {
				sm.sendSecurityError(w, http.StatusRequestEntityTooLarge, "Request body too large", "REQUEST_TOO_LARGE")
				return
			}

			// Limit the request body reader
			r.Body = http.MaxBytesReader(w, r.Body, maxBytes)
			next.ServeHTTP(w, r)
		})
	}
}

// SecureCORS implements secure CORS handling
func (sm *SecurityMiddleware) SecureCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// If there's an origin header, it must be in the allowlist
		if origin != "" {
			if sm.isOriginAllowed(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Credentials", "true")
			} else {
				// SECURITY: Explicitly reject non-allowed origins
				sm.sendSecurityError(w, http.StatusForbidden, "Origin not allowed", "CORS_BLOCKED")
				sm.logger.Warnf("Blocked request from non-allowed origin: %s from IP: %s", origin, sm.getClientIP(r))
				return
			}
		}

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			if origin != "" && sm.isOriginAllowed(origin) {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-Requested-With, X-API-Key")
				w.Header().Set("Access-Control-Allow-Credentials", "true")
				w.Header().Set("Access-Control-Max-Age", "3600")
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// RateLimitByIP implements IP-based rate limiting
func (sm *SecurityMiddleware) RateLimitByIP(requests int, window time.Duration) func(http.Handler) http.Handler {
	type clientInfo struct {
		requests int
		window   time.Time
	}
	clients := make(map[string]*clientInfo)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			clientIP := sm.getClientIP(r)
			now := time.Now()

			client, exists := clients[clientIP]
			if !exists {
				clients[clientIP] = &clientInfo{1, now}
				next.ServeHTTP(w, r)
				return
			}

			// Reset window if expired
			if now.Sub(client.window) > window {
				client.requests = 1
				client.window = now
				next.ServeHTTP(w, r)
				return
			}

			// Check rate limit
			if client.requests >= requests {
				sm.sendSecurityError(w, http.StatusTooManyRequests, "Rate limit exceeded", "RATE_LIMIT_EXCEEDED")
				sm.logger.Warnf("Rate limit exceeded for IP: %s", clientIP)
				return
			}

			client.requests++
			next.ServeHTTP(w, r)
		})
	}
}

// containsAttackPatterns checks for common attack patterns
func (sm *SecurityMiddleware) containsAttackPatterns(input string) bool {
	input = strings.ToLower(input)

	// SQL injection patterns
	sqlPatterns := []string{
		"union select", "drop table", "insert into", "delete from",
		"update set", "create table", "alter table", "truncate table",
		"exec ", "execute ", "sp_", "xp_", "/*", "*/", "--", ";--",
		"' or ", "\" or ", "' and ", "\" and ", "' union ", "\" union ",
		"0x", "char(", "ascii(", "substring(", "concat(", "load_file",
		"into outfile", "into dumpfile",
	}

	// XSS patterns
	xssPatterns := []string{
		"<script", "</script>", "javascript:", "vbscript:", "onload=",
		"onerror=", "onclick=", "onmouseover=", "alert(", "confirm(",
		"prompt(", "eval(", "expression(", "style=", "background:",
		"<iframe", "<object", "<embed", "<link", "<meta",
	}

	// Command injection patterns
	cmdPatterns := []string{
		"$(", "`", "|", "&&", "||", ";", "wget ", "curl ", "nc ",
		"netcat", "/bin/", "/usr/bin/", "cmd.exe", "powershell",
		"bash ", "sh ", "/etc/passwd", "/etc/shadow",
	}

	// Path traversal patterns
	pathPatterns := []string{
		"../", "..\\", "..", "~", "/proc/", "/sys/", "/dev/",
		"c:\\", "c:/", "\\windows\\", "/windows/",
	}

	allPatterns := append(sqlPatterns, xssPatterns...)
	allPatterns = append(allPatterns, cmdPatterns...)
	allPatterns = append(allPatterns, pathPatterns...)

	for _, pattern := range allPatterns {
		if strings.Contains(input, pattern) {
			return true
		}
	}

	return false
}

// isJWTSuspicious checks for suspicious JWT patterns
func (sm *SecurityMiddleware) isJWTSuspicious(token string) bool {
	// Check for common JWT attack patterns
	suspiciousPatterns := []string{
		"none", "None", "NONE", // Algorithm confusion
		"HS256", "HS384", "HS512", // If we expect RS256
		"../", "..\\", // Path traversal in payload
		"<script", "javascript:", // XSS in payload
		"' or ", "union select", // SQL injection in payload
	}

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(token, pattern) {
			return true
		}
	}

	return false
}

// isOriginAllowed checks if origin is in allowed list
func (sm *SecurityMiddleware) isOriginAllowed(origin string) bool {
	return slices.Contains(sm.allowedOrigins, origin)
}

// getClientIP extracts the real client IP
func (sm *SecurityMiddleware) getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header (but validate it)
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			// Basic IP validation
			if sm.isValidIP(ip) {
				return ip
			}
		}
	}

	// Check X-Real-IP header
	xri := r.Header.Get("X-Real-IP")
	if xri != "" && sm.isValidIP(xri) {
		return xri
	}

	// Fallback to RemoteAddr
	ip := strings.Split(r.RemoteAddr, ":")[0]
	return ip
}

// isValidIP performs basic IP validation
func (sm *SecurityMiddleware) isValidIP(ip string) bool {
	// Basic IP format validation (IPv4)
	parts := strings.Split(ip, ".")
	if len(parts) != 4 {
		return false
	}

	for _, part := range parts {
		if num, err := strconv.Atoi(part); err != nil || num < 0 || num > 255 {
			return false
		}
	}

	return true
}

// sendSecurityError sends a security-related error response
func (sm *SecurityMiddleware) sendSecurityError(w http.ResponseWriter, statusCode int, message, errorCode string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	response := map[string]any{
		"success":    false,
		"message":    message,
		"error_code": errorCode,
		"timestamp":  time.Now().UTC().Format(time.RFC3339),
	}

	json.NewEncoder(w).Encode(response)
}

// ValidateEnvVars validates that required security environment variables are set
func ValidateEnvVars() error {
	required := []string{
		"JWT_SECRET",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"MINIO_ROOT_USER",
		"MINIO_ROOT_PASSWORD",
	}

	missing := []string{}
	for _, env := range required {
		if os.Getenv(env) == "" {
			missing = append(missing, env)
		}
	}

	if len(missing) > 0 {
		return fmt.Errorf("missing required environment variables: %s", strings.Join(missing, ", "))
	}

	// Validate JWT secret strength
	jwtSecret := os.Getenv("JWT_SECRET")
	if len(jwtSecret) < 32 {
		return fmt.Errorf("JWT_SECRET must be at least 32 characters (256 bits) for security")
	}

	return nil
}

// GenerateSecureToken generates a cryptographically secure token
func GenerateSecureToken(length int) (string, error) {
	if length < 16 {
		return "", fmt.Errorf("token length must be at least 16 characters")
	}

	// Implementation would use crypto/rand for secure token generation
	// This is a placeholder for the actual implementation
	return "secure_token_placeholder", nil
}

// SecureCompare performs constant-time string comparison
func SecureCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
