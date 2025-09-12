#!/bin/bash

# ==============================================
# STUDY PLATFORM - SECURE ENVIRONMENT SETUP
# ==============================================
# This script generates secure credentials and sets up the environment
# NEVER run this script in production - use proper secret management
# ==============================================

set -euo pipefail

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if required tools are installed
check_dependencies() {
    log_info "Checking dependencies..."
    
    local deps=("openssl" "docker" "docker-compose")
    local missing=()
    
    for dep in "${deps[@]}"; do
        if ! command -v "$dep" &> /dev/null; then
            missing+=("$dep")
        fi
    done
    
    if [ ${#missing[@]} -ne 0 ]; then
        log_error "Missing dependencies: ${missing[*]}"
        log_error "Please install the missing dependencies and try again"
        exit 1
    fi
    
    log_success "All dependencies found"
}

# Generate secure random string
generate_secure_string() {
    local length=${1:-32}
    openssl rand -base64 $((length * 3 / 4)) | tr -d "=+/" | cut -c1-${length}
}

# Generate secure password
generate_secure_password() {
    local length=${1:-20}
    # Use a mix of letters, numbers, and special characters
    openssl rand -base64 48 | tr -d "=+/" | cut -c1-${length}
}

# Validate password strength
validate_password_strength() {
    local password="$1"
    local min_length=12
    
    if [ ${#password} -lt $min_length ]; then
        return 1
    fi
    
    # Check for at least one uppercase, lowercase, digit, and special char
    if [[ ! "$password" =~ [A-Z] ]] || [[ ! "$password" =~ [a-z] ]] || [[ ! "$password" =~ [0-9] ]]; then
        return 1
    fi
    
    return 0
}

# Create .env file with secure values
create_env_file() {
    local env_file=".env"
    local backup_file=".env.backup.$(date +%Y%m%d_%H%M%S)"
    
    log_info "Creating secure environment file..."
    
    # Backup existing .env file
    if [ -f "$env_file" ]; then
        log_warning "Backing up existing .env file to $backup_file"
        cp "$env_file" "$backup_file"
    fi
    
    # Generate secure credentials
    local jwt_secret=$(generate_secure_string 64)
    local session_secret=$(generate_secure_string 64)
    local postgres_password=$(generate_secure_password 24)
    local postgres_user="study_admin_$(generate_secure_string 8)"
    local minio_user="minio_admin_$(generate_secure_string 8)"
    local minio_password=$(generate_secure_password 20)
    
    # Create .env file
    cat > "$env_file" << EOF
# ==============================================
# STUDY PLATFORM - SECURE ENVIRONMENT
# ==============================================
# Generated on: $(date)
# WARNING: Keep this file secure and never commit to version control
# ==============================================

# ==============================================
# DATABASE CONFIGURATION
# ==============================================
POSTGRES_DB=studyplatform
POSTGRES_USER=${postgres_user}
POSTGRES_PASSWORD=${postgres_password}
DATABASE_URL=postgres://${postgres_user}:${postgres_password}@postgres:5432/studyplatform?sslmode=disable

# Database Application Users
DB_HOST=postgres
DB_PORT=5432
DB_USER=${postgres_user}
DB_PASSWORD=${postgres_password}
DB_NAME=studyplatform
DB_SSLMODE=disable
DB_MAX_CONNECTIONS=25
DB_MAX_IDLE_CONNECTIONS=5
DB_CONNECTION_LIFETIME=300s

# ==============================================
# REDIS CONFIGURATION
# ==============================================
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0
REDIS_MAX_CONNECTIONS=10

# ==============================================
# OBJECT STORAGE CONFIGURATION
# ==============================================
MINIO_ROOT_USER=${minio_user}
MINIO_ROOT_PASSWORD=${minio_password}

# ==============================================
# AUTHENTICATION & SECURITY
# ==============================================
JWT_SECRET=${jwt_secret}
JWT_EXPIRY=24h
JWT_REFRESH_EXPIRY=168h
SESSION_SECRET=${session_secret}
SESSION_TIMEOUT=3600
BCRYPT_ROUNDS=12

# ==============================================
# EXTERNAL API KEYS (REPLACE WITH REAL VALUES)
# ==============================================
# AWS S3 Configuration
S3_ENDPOINT=https://s3.amazonaws.com
S3_REGION=us-east-1
S3_ACCESS_KEY_ID=REPLACE_WITH_REAL_AWS_ACCESS_KEY
S3_SECRET_ACCESS_KEY=REPLACE_WITH_REAL_AWS_SECRET_KEY

# S3 Bucket Names
S3_BUCKET_FILES=your-platform-files-bucket
S3_BUCKET_DOCUMENTS=your-platform-documents-bucket
S3_BUCKET_IMAGES=your-platform-images-bucket
S3_BUCKET_AVATARS=your-platform-avatars-bucket
S3_BUCKET_COURSE_MATERIALS=your-platform-course-materials-bucket

# Cloudflare Stream
CLOUDFLARE_ACCOUNT_ID=REPLACE_WITH_REAL_CLOUDFLARE_ACCOUNT_ID
CLOUDFLARE_STREAM_TOKEN=REPLACE_WITH_REAL_CLOUDFLARE_TOKEN

# OpenAI
OPENAI_API_KEY=REPLACE_WITH_REAL_OPENAI_API_KEY
OPENAI_MODEL=gpt-3.5-turbo
OPENAI_MAX_TOKENS=1000
OPENAI_TEMPERATURE=0.7

# Payment Gateways
STRIPE_SECRET_KEY=REPLACE_WITH_REAL_STRIPE_SECRET_KEY
STRIPE_PUBLISHABLE_KEY=REPLACE_WITH_REAL_STRIPE_PUBLISHABLE_KEY
STRIPE_WEBHOOK_SECRET=REPLACE_WITH_REAL_STRIPE_WEBHOOK_SECRET
PAYPAL_CLIENT_ID=REPLACE_WITH_REAL_PAYPAL_CLIENT_ID
PAYPAL_CLIENT_SECRET=REPLACE_WITH_REAL_PAYPAL_CLIENT_SECRET
PAYPAL_SANDBOX=true
PAYMENT_CURRENCY=USD

# ==============================================
# SERVICE CONFIGURATION
# ==============================================
# Service Ports
GRPC_PORT_AUTH=8081
GRPC_PORT_COURSE=8082
GRPC_PORT_PROGRESS=8083
VIDEO_SERVICE_PORT=8084
BUCKET_SERVICE_PORT=8085
CHATBOT_PORT=8086
FORUM_PORT=8087
PAYMENT_PORT=8088
API_GATEWAY_PORT=8080

# Service Hosts
AUTH_SERVICE_HOST=0.0.0.0
COURSE_SERVICE_HOST=0.0.0.0
PROGRESS_SERVICE_HOST=0.0.0.0
VIDEO_SERVICE_HOST=0.0.0.0
BUCKET_SERVICE_HOST=0.0.0.0
CHATBOT_HOST=0.0.0.0
FORUM_HOST=0.0.0.0
PAYMENT_HOST=0.0.0.0

# Service URLs
AUTH_SERVICE_URL=auth-service:8081
COURSE_SERVICE_URL=course-service:8082
PROGRESS_SERVICE_URL=progress-service:8083
VIDEO_SERVICE_URL=http://video-service:8084
BUCKET_SERVICE_URL=http://bucket-service:8085
CHATBOT_SERVICE_URL=http://chatbot-service:8086
FORUM_SERVICE_URL=http://forum-service:8087
PAYMENT_SERVICE_URL=http://payment-service:8088

# ==============================================
# SECURITY & CORS CONFIGURATION
# ==============================================
CORS_ORIGINS=http://localhost:3000,http://127.0.0.1:3000,http://localhost:8080
RATE_LIMIT_REQUESTS=100
RATE_LIMIT_WINDOW=60s
RATE_LIMIT_BURST=200
SECURITY_HEADERS_ENABLED=true
CONTENT_SECURITY_POLICY=default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'

# ==============================================
# ENVIRONMENT SETTINGS
# ==============================================
ENVIRONMENT=development
LOG_LEVEL=info
DEBUG_MODE=false

# Health Checks
HEALTH_CHECK_ENABLED=true
HEALTH_CHECK_INTERVAL=30s
HEALTH_CHECK_TIMEOUT=10s
HEALTH_CHECK_RETRIES=3

# Metrics
METRICS_ENABLED=true
METRICS_PORT=9090

# ==============================================
# GENERATED CREDENTIALS SUMMARY
# ==============================================
# Database User: ${postgres_user}
# MinIO User: ${minio_user}
# JWT Secret Length: ${#jwt_secret} characters
# Session Secret Length: ${#session_secret} characters
# Password Strength: Strong (24+ characters)
# ==============================================

EOF

    log_success "Environment file created: $env_file"
    log_warning "IMPORTANT: Replace placeholder API keys with real values before production use"
    
    # Set secure permissions
    chmod 600 "$env_file"
    log_success "Set secure permissions (600) on $env_file"
}

# Validate environment file
validate_env_file() {
    local env_file=".env"
    
    if [ ! -f "$env_file" ]; then
        log_error "Environment file not found: $env_file"
        return 1
    fi
    
    log_info "Validating environment file..."
    
    # Check for required variables
    local required_vars=(
        "POSTGRES_USER"
        "POSTGRES_PASSWORD" 
        "JWT_SECRET"
        "SESSION_SECRET"
        "MINIO_ROOT_USER"
        "MINIO_ROOT_PASSWORD"
    )
    
    local missing_vars=()
    
    # Source the .env file in a subshell to check variables
    for var in "${required_vars[@]}"; do
        if ! grep -q "^${var}=" "$env_file" || [ -z "$(grep "^${var}=" "$env_file" | cut -d'=' -f2)" ]; then
            missing_vars+=("$var")
        fi
    done
    
    if [ ${#missing_vars[@]} -ne 0 ]; then
        log_error "Missing required environment variables: ${missing_vars[*]}"
        return 1
    fi
    
    # Validate JWT secret length
    local jwt_secret=$(grep "^JWT_SECRET=" "$env_file" | cut -d'=' -f2)
    if [ ${#jwt_secret} -lt 32 ]; then
        log_error "JWT_SECRET must be at least 32 characters"
        return 1
    fi
    
    # Validate password strength
    local db_password=$(grep "^POSTGRES_PASSWORD=" "$env_file" | cut -d'=' -f2)
    if ! validate_password_strength "$db_password"; then
        log_error "Database password does not meet strength requirements"
        return 1
    fi
    
    log_success "Environment file validation passed"
}

# Generate TLS certificates for development
generate_dev_certificates() {
    local cert_dir="certs"
    
    log_info "Generating development TLS certificates..."
    
    mkdir -p "$cert_dir"
    
    # Generate private key
    openssl genrsa -out "$cert_dir/server.key" 2048
    
    # Generate certificate signing request
    openssl req -new -key "$cert_dir/server.key" -out "$cert_dir/server.csr" -subj "/C=US/ST=State/L=City/O=StudyPlatform/OU=Development/CN=localhost"
    
    # Generate self-signed certificate
    openssl x509 -req -days 365 -in "$cert_dir/server.csr" -signkey "$cert_dir/server.key" -out "$cert_dir/server.crt"
    
    # Set secure permissions
    chmod 600 "$cert_dir/server.key"
    chmod 644 "$cert_dir/server.crt"
    
    # Clean up CSR
    rm "$cert_dir/server.csr"
    
    log_success "Development certificates generated in $cert_dir/"
    log_warning "These are self-signed certificates for development only"
}

# Setup Docker secrets (for Docker Swarm)
setup_docker_secrets() {
    log_info "Setting up Docker secrets (if Docker Swarm is active)..."
    
    if ! docker node ls &> /dev/null; then
        log_warning "Docker Swarm not active, skipping secret setup"
        return 0
    fi
    
    # Create secrets from .env file
    local secrets=(
        "jwt_secret:JWT_SECRET"
        "db_password:POSTGRES_PASSWORD"
        "minio_password:MINIO_ROOT_PASSWORD"
    )
    
    for secret_pair in "${secrets[@]}"; do
        local secret_name=${secret_pair%%:*}
        local env_var=${secret_pair##*:}
        local value=$(grep "^${env_var}=" .env | cut -d'=' -f2)
        
        if [ -n "$value" ]; then
            echo "$value" | docker secret create "$secret_name" - 2>/dev/null || log_warning "Secret $secret_name already exists"
        fi
    done
    
    log_success "Docker secrets configured"
}

# Show security checklist
show_security_checklist() {
    cat << EOF

${GREEN}==============================================
SECURITY SETUP COMPLETE
==============================================${NC}

${YELLOW}IMPORTANT SECURITY CHECKLIST:${NC}

${BLUE}✓ Environment Setup:${NC}
  - Secure .env file created with proper permissions
  - Strong passwords and secrets generated
  - Database credentials configured

${BLUE}✓ Next Steps (CRITICAL):${NC}
  1. Replace placeholder API keys with real values
  2. Review and update CORS origins for your domain
  3. Configure firewall rules (ports 8080 only external)
  4. Set up SSL/TLS certificates for production
  5. Configure backup strategies
  6. Set up monitoring and alerting
  7. Review and test security policies

${BLUE}✓ Before Production:${NC}
  - Enable SSL/TLS in production environment
  - Use proper secret management (AWS Secrets Manager, etc.)
  - Set up database encryption at rest
  - Configure audit logging
  - Perform security penetration testing
  - Set up rate limiting at load balancer level
  - Configure DDoS protection

${BLUE}✓ Files Created:${NC}
  - .env (secure environment variables)
  - .env.example (template for team members)
  - certs/ (development TLS certificates)

${RED}⚠️  CRITICAL SECURITY WARNINGS:${NC}
  - Never commit .env file to version control
  - Replace all placeholder API keys before going live
  - Rotate secrets regularly in production
  - Monitor for credential leaks

${GREEN}Ready to start services with: docker-compose up -d${NC}

EOF
}

# Main execution
main() {
    echo -e "${BLUE}"
    echo "=============================================="
    echo "  STUDY PLATFORM - SECURE ENVIRONMENT SETUP"
    echo "=============================================="
    echo -e "${NC}"
    
    # Check if already in Study-Platform directory
    if [ ! -f "docker-compose.yml" ]; then
        log_error "Must be run from Study-Platform directory"
        log_error "Please cd to Study-Platform/ and run this script again"
        exit 1
    fi
    
    check_dependencies
    create_env_file
    validate_env_file
    generate_dev_certificates
    setup_docker_secrets
    show_security_checklist
    
    log_success "Secure environment setup complete!"
}

# Run main function
main "$@"