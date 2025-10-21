#!/bin/bash

# Swagger Documentation Generation Script
# Generates OpenAPI/Swagger documentation for all microservices

set -e

echo "🚀 Generating Swagger documentation for all services..."

# Colors for output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "❌ swag CLI not found. Installing..."
    go install github.com/swaggo/swag/cmd/swag@latest
fi

# Add swag path to PATH if needed
export PATH=$PATH:$(go env GOPATH)/bin

# Services to generate docs for
SERVICES=(
    "api-gateway"
    "auth-service"
    "course-service"
    "progress-service"
    "video-service"
    "bucket-service"
    "chatbot-service"
    "forum-service"
    "payment-service"
    "instructor-dashboard-service"
)

# Project root directory
PROJECT_ROOT=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)

# Generate docs for each service
for service in "${SERVICES[@]}"; do
    echo -e "${BLUE}📝 Generating docs for ${service}...${NC}"
    
    cd "$PROJECT_ROOT/$service"
    
    # Create docs directory if it doesn't exist
    mkdir -p docs
    
    # Generate Swagger docs
    if [ -f "cmd/main.go" ]; then
        swag init -g cmd/main.go -o docs --parseDependency --parseInternal
        echo -e "${GREEN}✅ ${service} docs generated${NC}"
    else
        echo "⚠️  Skipping ${service} (no cmd/main.go found)"
    fi
done

echo -e "${GREEN}✨ All Swagger documentation generated successfully!${NC}"
echo ""
echo "📚 Access documentation at:"
echo "  • API Gateway: http://localhost:8080/api/v1/docs"
echo "  • Or individual services: http://localhost:PORT/swagger/index.html"
