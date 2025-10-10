#!/bin/bash
set -e

echo "=========================================="
echo "Setting up Study Platform EC2 Instance"
echo "=========================================="

# Update system
echo "Updating system packages..."
sudo apt-get update
sudo apt-get upgrade -y

# Install Docker
echo "Installing Docker..."
if ! command -v docker &> /dev/null; then
    curl -fsSL https://get.docker.com -o get-docker.sh
    sudo sh get-docker.sh
    sudo usermod -aG docker ubuntu
    rm get-docker.sh
    echo "Docker installed successfully"
else
    echo "Docker already installed"
fi

# Install Docker Compose
echo "Installing Docker Compose..."
if ! command -v docker compose &> /dev/null; then
    sudo apt-get install -y docker-compose-plugin
    echo "Docker Compose installed successfully"
else
    echo "Docker Compose already installed"
fi

# Install Git
echo "Installing Git..."
sudo apt-get install -y git

# Install Nginx
echo "Installing Nginx..."
sudo apt-get install -y nginx

# Install certbot for SSL (optional)
echo "Installing Certbot..."
sudo apt-get install -y certbot python3-certbot-nginx

# Create application directory
echo "Creating application directory..."
sudo mkdir -p /opt/study-platform
sudo chown -R ubuntu:ubuntu /opt/study-platform

# Configure firewall (if ufw is enabled)
if sudo ufw status | grep -q "Status: active"; then
    echo "Configuring firewall..."
    sudo ufw allow 22/tcp
    sudo ufw allow 80/tcp
    sudo ufw allow 443/tcp
    sudo ufw allow 8080:8089/tcp
fi

# Create .env file template
echo "Creating environment file template..."
cat > /opt/study-platform/.env.example << 'EOF'
# Database Configuration
DATABASE_URL=postgresql://USERNAME:PASSWORD@postgres:5432/studyplatform?sslmode=disable
POSTGRES_USER=your-db-username
POSTGRES_PASSWORD=your-db-password
POSTGRES_DB=studyplatform

# JWT Configuration
JWT_SECRET=your-super-secret-jwt-key-change-this-in-production

# Redis Configuration
REDIS_HOST=redis
REDIS_PORT=6379
REDIS_PASSWORD=

# Service URLs (Internal Docker network)
AUTH_SERVICE_URL=auth-service:8081
COURSE_SERVICE_URL=course-service:8082
PROGRESS_SERVICE_URL=progress-service:8083
VIDEO_SERVICE_URL=video-service:8084
BUCKET_SERVICE_URL=bucket-service:8085
CHATBOT_SERVICE_URL=chatbot-service:8086
FORUM_SERVICE_URL=forum-service:8087
PAYMENT_SERVICE_URL=payment-service:8088
INSTRUCTOR_DASHBOARD_SERVICE_URL=instructor-dashboard-service:8089

# AWS S3 Configuration
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=your-access-key
AWS_SECRET_ACCESS_KEY=your-secret-key
S3_BUCKET_NAME=study-platform-files

# Cloudflare Stream Configuration
CLOUDFLARE_ACCOUNT_ID=your-account-id
CLOUDFLARE_API_TOKEN=your-api-token
CLOUDFLARE_STREAM_URL=https://customer-xxx.cloudflarestream.com

# OpenAI Configuration
OPENAI_API_KEY=your-openai-api-key

# OAuth Configuration
GOOGLE_CLIENT_ID=your-google-client-id
GOOGLE_CLIENT_SECRET=your-google-client-secret
GITHUB_CLIENT_ID=your-github-client-id
GITHUB_CLIENT_SECRET=your-github-client-secret
FACEBOOK_CLIENT_ID=your-facebook-client-id
FACEBOOK_CLIENT_SECRET=your-facebook-client-secret

# Payment Gateway Configuration
STRIPE_SECRET_KEY=your-stripe-secret-key
STRIPE_WEBHOOK_SECRET=your-stripe-webhook-secret
PAYPAL_CLIENT_ID=your-paypal-client-id
PAYPAL_CLIENT_SECRET=your-paypal-client-secret
LEMON_SQUEEZY_API_KEY=your-lemon-squeezy-api-key

# MinIO Configuration (for development)
MINIO_ENDPOINT=minio:9000
MINIO_ACCESS_KEY=your-minio-access-key
MINIO_SECRET_KEY=your-minio-secret-key
MINIO_USE_SSL=false
EOF

echo "=========================================="
echo "EC2 Setup completed successfully!"
echo "=========================================="
echo ""
echo "Next steps:"
echo "1. Copy the .env.example to .env and configure with your values"
echo "2. Clone the repository: cd /opt/study-platform && git clone https://github.com/khoipn21/study-platform.git ."
echo "3. Start services: docker compose up -d"
echo "4. Configure Nginx for domain routing"
echo ""
echo "Instance is ready for deployment!"
