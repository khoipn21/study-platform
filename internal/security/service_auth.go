package security

import (
	"crypto/rand"
	"crypto/subtle"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/dgrijalva/jwt-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

// ServiceAuthConfig holds service authentication configuration
type ServiceAuthConfig struct {
	ServiceName    string
	SharedSecret   string
	TLSConfig      *tls.Config
	TokenDuration  time.Duration
	AllowedServices []string
}

// ServiceAuthenticator handles inter-service authentication
type ServiceAuthenticator struct {
	config        ServiceAuthConfig
	serviceTokens map[string]string
	mutex         sync.RWMutex
}

// ServiceClaims represents JWT claims for service-to-service communication
type ServiceClaims struct {
	ServiceName string   `json:"service_name"`
	Permissions []string `json:"permissions"`
	jwt.StandardClaims
}

// NewServiceAuthenticator creates a new service authenticator
func NewServiceAuthenticator(config ServiceAuthConfig) *ServiceAuthenticator {
	if config.TokenDuration == 0 {
		config.TokenDuration = 1 * time.Hour
	}

	return &ServiceAuthenticator{
		config:        config,
		serviceTokens: make(map[string]string),
	}
}

// GenerateServiceToken generates a JWT token for service-to-service communication
func (sa *ServiceAuthenticator) GenerateServiceToken(targetService string, permissions []string) (string, error) {
	now := time.Now()
	
	claims := ServiceClaims{
		ServiceName: sa.config.ServiceName,
		Permissions: permissions,
		StandardClaims: jwt.StandardClaims{
			Subject:   sa.config.ServiceName,
			Audience:  targetService,
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(sa.config.TokenDuration).Unix(),
			Issuer:    "study-platform-auth",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(sa.config.SharedSecret))
}

// ValidateServiceToken validates a JWT token from another service
func (sa *ServiceAuthenticator) ValidateServiceToken(tokenString string) (*ServiceClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &ServiceClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(sa.config.SharedSecret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	claims, ok := token.Claims.(*ServiceClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Verify the token is intended for this service
	if claims.Audience != sa.config.ServiceName && claims.Audience != "all" {
		return nil, fmt.Errorf("token not intended for this service")
	}

	// Check if the source service is allowed
	if !sa.isServiceAllowed(claims.ServiceName) {
		return nil, fmt.Errorf("service %s not allowed to communicate with %s", claims.ServiceName, sa.config.ServiceName)
	}

	return claims, nil
}

// isServiceAllowed checks if a service is in the allowed list
func (sa *ServiceAuthenticator) isServiceAllowed(serviceName string) bool {
	for _, allowed := range sa.config.AllowedServices {
		if allowed == serviceName || allowed == "*" {
			return true
		}
	}
	return false
}

// gRPC Interceptors

// UnaryServerInterceptor validates service tokens for unary gRPC calls
func (sa *ServiceAuthenticator) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		md, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			return nil, fmt.Errorf("missing metadata")
		}

		tokens := md.Get("authorization")
		if len(tokens) == 0 {
			return nil, fmt.Errorf("missing authorization token")
		}

		tokenString := tokens[0]
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = tokenString[7:]
		}

		claims, err := sa.ValidateServiceToken(tokenString)
		if err != nil {
			return nil, fmt.Errorf("invalid service token: %w", err)
		}

		// Add claims to context for use by handlers
		ctx = context.WithValue(ctx, "service_claims", claims)
		return handler(ctx, req)
	}
}

// UnaryClientInterceptor adds service tokens to outgoing unary gRPC calls
func (sa *ServiceAuthenticator) UnaryClientInterceptor(targetService string, permissions []string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		token, err := sa.GenerateServiceToken(targetService, permissions)
		if err != nil {
			return fmt.Errorf("failed to generate service token: %w", err)
		}

		ctx = metadata.AppendToOutgoingContext(ctx, "authorization", "Bearer "+token)
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// HTTP Middleware

// HTTPServiceAuthMiddleware validates service tokens for HTTP requests
func (sa *ServiceAuthenticator) HTTPServiceAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("X-Service-Auth")
		if authHeader == "" {
			http.Error(w, "Missing service authentication", http.StatusUnauthorized)
			return
		}

		claims, err := sa.ValidateServiceToken(authHeader)
		if err != nil {
			http.Error(w, "Invalid service token: "+err.Error(), http.StatusUnauthorized)
			return
		}

		// Add service information to request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "service_claims", claims)
		ctx = context.WithValue(ctx, "service_name", claims.ServiceName)
		ctx = context.WithValue(ctx, "service_permissions", claims.Permissions)

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

// HTTPServiceAuthClient adds service authentication to outgoing HTTP requests
type HTTPServiceAuthClient struct {
	authenticator  *ServiceAuthenticator
	targetService  string
	permissions    []string
	client         *http.Client
}

// NewHTTPServiceAuthClient creates a new authenticated HTTP client
func NewHTTPServiceAuthClient(authenticator *ServiceAuthenticator, targetService string, permissions []string) *HTTPServiceAuthClient {
	return &HTTPServiceAuthClient{
		authenticator: authenticator,
		targetService: targetService,
		permissions:   permissions,
		client:        &http.Client{Timeout: 30 * time.Second},
	}
}

// Do executes an HTTP request with service authentication
func (c *HTTPServiceAuthClient) Do(req *http.Request) (*http.Response, error) {
	token, err := c.authenticator.GenerateServiceToken(c.targetService, c.permissions)
	if err != nil {
		return nil, fmt.Errorf("failed to generate service token: %w", err)
	}

	req.Header.Set("X-Service-Auth", token)
	req.Header.Set("User-Agent", "StudyPlatform-Service/"+c.authenticator.config.ServiceName)

	return c.client.Do(req)
}

// TLS Configuration

// GenerateTLSConfig creates a secure TLS configuration
func GenerateTLSConfig(serverMode bool, certFile, keyFile, caFile string) (*tls.Config, error) {
	config := &tls.Config{
		MinVersion:               tls.VersionTLS12,
		CurvePreferences:         []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
	}

	if serverMode {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load server certificates: %w", err)
		}
		config.Certificates = []tls.Certificate{cert}
	}

	if caFile != "" {
		caCert, err := tls.LoadX509KeyPair(caFile, caFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load CA certificate: %w", err)
		}
		config.RootCAs.AddCert(caCert.Leaf)
	}

	return config, nil
}

// CreateGRPCServerCredentials creates gRPC server credentials with TLS
func CreateGRPCServerCredentials(certFile, keyFile string) (credentials.TransportCredentials, error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load TLS credentials: %w", err)
	}

	config := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
		NextProtos:   []string{"h2"},
	}

	return credentials.NewTLS(config), nil
}

// CreateGRPCClientCredentials creates gRPC client credentials with TLS
func CreateGRPCClientCredentials(serverName string, insecure bool) (credentials.TransportCredentials, error) {
	if insecure {
		return credentials.NewTLS(&tls.Config{
			InsecureSkipVerify: true,
			ServerName:         serverName,
		}), nil
	}

	return credentials.NewTLS(&tls.Config{
		ServerName: serverName,
		MinVersion: tls.VersionTLS12,
	}), nil
}

// Service Discovery and Health Checking

// ServiceRegistry maintains a registry of available services
type ServiceRegistry struct {
	services map[string]*ServiceInfo
	mutex    sync.RWMutex
}

// ServiceInfo contains information about a registered service
type ServiceInfo struct {
	Name        string    `json:"name"`
	Address     string    `json:"address"`
	Port        int       `json:"port"`
	Protocol    string    `json:"protocol"` // "grpc" or "http"
	HealthCheck string    `json:"health_check"`
	LastSeen    time.Time `json:"last_seen"`
	Status      string    `json:"status"` // "healthy", "unhealthy", "unknown"
}

// NewServiceRegistry creates a new service registry
func NewServiceRegistry() *ServiceRegistry {
	return &ServiceRegistry{
		services: make(map[string]*ServiceInfo),
	}
}

// RegisterService registers a service in the registry
func (sr *ServiceRegistry) RegisterService(info *ServiceInfo) {
	sr.mutex.Lock()
	defer sr.mutex.Unlock()

	info.LastSeen = time.Now()
	if info.Status == "" {
		info.Status = "unknown"
	}

	sr.services[info.Name] = info
}

// GetService retrieves service information by name
func (sr *ServiceRegistry) GetService(name string) (*ServiceInfo, error) {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()

	service, exists := sr.services[name]
	if !exists {
		return nil, fmt.Errorf("service %s not found", name)
	}

	// Check if service is stale (not seen in last 5 minutes)
	if time.Since(service.LastSeen) > 5*time.Minute {
		service.Status = "unknown"
	}

	return service, nil
}

// ListServices returns all registered services
func (sr *ServiceRegistry) ListServices() map[string]*ServiceInfo {
	sr.mutex.RLock()
	defer sr.mutex.RUnlock()

	result := make(map[string]*ServiceInfo)
	for name, service := range sr.services {
		serviceCopy := *service
		result[name] = &serviceCopy
	}

	return result
}

// HealthChecker performs health checks on registered services
type HealthChecker struct {
	registry *ServiceRegistry
	interval time.Duration
	client   *http.Client
}

// NewHealthChecker creates a new health checker
func NewHealthChecker(registry *ServiceRegistry, interval time.Duration) *HealthChecker {
	return &HealthChecker{
		registry: registry,
		interval: interval,
		client:   &http.Client{Timeout: 5 * time.Second},
	}
}

// Start begins the health checking process
func (hc *HealthChecker) Start() {
	ticker := time.NewTicker(hc.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			hc.checkAllServices()
		}
	}
}

// checkAllServices performs health checks on all registered services
func (hc *HealthChecker) checkAllServices() {
	services := hc.registry.ListServices()

	for _, service := range services {
		go hc.checkService(service)
	}
}

// checkService performs a health check on a specific service
func (hc *HealthChecker) checkService(service *ServiceInfo) {
	if service.HealthCheck == "" {
		return
	}

	resp, err := hc.client.Get(service.HealthCheck)
	if err != nil {
		hc.updateServiceStatus(service.Name, "unhealthy")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		hc.updateServiceStatus(service.Name, "healthy")
	} else {
		hc.updateServiceStatus(service.Name, "unhealthy")
	}
}

// updateServiceStatus updates the status of a service
func (hc *HealthChecker) updateServiceStatus(serviceName, status string) {
	hc.registry.mutex.Lock()
	defer hc.registry.mutex.Unlock()

	if service, exists := hc.registry.services[serviceName]; exists {
		service.Status = status
		service.LastSeen = time.Now()
	}
}

// Security Utilities

// GenerateServiceSecret generates a cryptographically secure secret for service authentication
func GenerateServiceSecret() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("failed to generate random bytes: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// SecureStringCompare performs a constant-time string comparison
func SecureStringCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

// CreateServiceWhitelist creates a whitelist of allowed services for each service
func CreateServiceWhitelist() map[string][]string {
	return map[string][]string{
		"api-gateway": {
			"auth-service",
			"course-service", 
			"progress-service",
			"video-service",
			"bucket-service",
			"chatbot-service",
			"forum-service",
			"payment-service",
		},
		"auth-service": {
			"api-gateway",
		},
		"course-service": {
			"api-gateway",
			"progress-service", // For enrollment checks
		},
		"progress-service": {
			"api-gateway",
			"course-service", // For course validation
			"payment-service", // For enrollment activation
		},
		"video-service": {
			"api-gateway",
			"course-service", // For course video validation
		},
		"bucket-service": {
			"api-gateway",
			"course-service", // For course materials
			"video-service",  // For video uploads
		},
		"chatbot-service": {
			"api-gateway",
			"course-service", // For course context
		},
		"forum-service": {
			"api-gateway",
			"course-service", // For course forums
		},
		"payment-service": {
			"api-gateway",
			"progress-service", // For enrollment activation
		},
	}
}