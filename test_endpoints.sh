#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Base URL
BASE_URL="http://localhost:8080/api/v1"

# Test variables
TEST_EMAIL="test@example.com"
TEST_PASSWORD="password123"
TEST_USERNAME="testuser"
TOKEN=""

echo -e "${BLUE}=== Study Platform API Endpoint Testing ===${NC}"
echo ""

# Function to make HTTP requests and display results
test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local expected_status=$4
    local description=$5
    local auth_header=$6

    echo -e "${YELLOW}Testing: $description${NC}"
    echo "  $method $endpoint"

    if [ -n "$auth_header" ]; then
        if [ -n "$data" ]; then
            response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
                -H "Content-Type: application/json" \
                -H "Authorization: Bearer $TOKEN" \
                -d "$data")
        else
            response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
                -H "Authorization: Bearer $TOKEN")
        fi
    else
        if [ -n "$data" ]; then
            response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint" \
                -H "Content-Type: application/json" \
                -d "$data")
        else
            response=$(curl -s -w "\n%{http_code}" -X $method "$BASE_URL$endpoint")
        fi
    fi

    # Split response and status code
    status_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | head -n -1)

    if [ "$status_code" -eq "$expected_status" ]; then
        echo -e "  ${GREEN}✓ Status: $status_code (Expected: $expected_status)${NC}"
    else
        echo -e "  ${RED}✗ Status: $status_code (Expected: $expected_status)${NC}"
    fi

    if [ -n "$response_body" ]; then
        echo "  Response: $response_body"
    fi
    echo ""
}

# 1. Health Check Endpoints
echo -e "${BLUE}1. Health Check Endpoints${NC}"
test_endpoint "GET" "/health" "" 200 "Basic health check"
test_endpoint "GET" "/health/circuit-breakers" "" 200 "Circuit breaker status"

# 2. Authentication Endpoints
echo -e "${BLUE}2. Authentication Endpoints${NC}"

# User Registration
test_endpoint "POST" "/auth/register" \
    "{\"username\":\"$TEST_USERNAME\",\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}" \
    201 "User registration"

# User Login
login_response=$(curl -s -X POST "$BASE_URL/auth/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}")

echo -e "${YELLOW}Testing: User login${NC}"
echo "  POST /auth/login"
echo "  Response: $login_response"

# Extract token if login was successful
if echo "$login_response" | grep -q "token"; then
    TOKEN=$(echo "$login_response" | grep -o '"token":"[^"]*"' | cut -d'"' -f4)
    echo -e "  ${GREEN}✓ Login successful, token extracted${NC}"
else
    echo -e "  ${RED}✗ Login failed or token not found${NC}"
fi
echo ""

# Token validation
if [ -n "$TOKEN" ]; then
    test_endpoint "POST" "/auth/validate" \
        "{\"token\":\"$TOKEN\"}" \
        200 "Token validation"

    # User profile
    test_endpoint "GET" "/auth/profile" "" 200 "Get user profile" "auth"
fi

# OAuth endpoints (these will return URLs)
test_endpoint "GET" "/auth/oauth/google/url" "" 200 "Get Google OAuth URL"
test_endpoint "GET" "/auth/oauth/github/url" "" 200 "Get GitHub OAuth URL"
test_endpoint "GET" "/auth/oauth/facebook/url" "" 200 "Get Facebook OAuth URL"

# 3. Course Management Endpoints
echo -e "${BLUE}3. Course Management Endpoints${NC}"

# List courses (should work without authentication)
test_endpoint "GET" "/courses" "" 200 "List all courses"
test_endpoint "GET" "/courses?page=1&page_size=5" "" 200 "List courses with pagination"

# Search courses
test_endpoint "GET" "/courses/search?q=programming" "" 200 "Search courses"

# Create course (requires authentication)
if [ -n "$TOKEN" ]; then
    course_data="{\"title\":\"Test Course\",\"description\":\"A test course\",\"category\":\"Programming\",\"level\":\"beginner\",\"price\":99.99,\"currency\":\"USD\",\"tags\":[\"test\",\"programming\"]}"
    test_endpoint "POST" "/courses" "$course_data" 201 "Create new course" "auth"
fi

# 4. Enrollment Endpoints
echo -e "${BLUE}4. Enrollment Endpoints${NC}"

if [ -n "$TOKEN" ]; then
    test_endpoint "GET" "/enrollments" "" 200 "List user enrollments" "auth"

    # Try to enroll in a course (this may fail if no courses exist)
    test_endpoint "POST" "/enrollments" \
        "{\"course_id\":\"123e4567-e89b-12d3-a456-426614174000\"}" \
        400 "Create enrollment (expected to fail with invalid course ID)" "auth"
fi

# 5. Progress Tracking Endpoints
echo -e "${BLUE}5. Progress Tracking Endpoints${NC}"

if [ -n "$TOKEN" ]; then
    # Update progress (may fail without valid course/lecture IDs)
    progress_data="{\"course_id\":\"123e4567-e89b-12d3-a456-426614174000\",\"lecture_id\":\"123e4567-e89b-12d3-a456-426614174001\",\"progress_percentage\":85.5,\"watch_time_seconds\":1200,\"completed\":false}"
    test_endpoint "POST" "/progress/update" "$progress_data" 400 "Update progress (expected to fail)" "auth"

    # Complete lecture
    complete_data="{\"course_id\":\"123e4567-e89b-12d3-a456-426614174000\",\"lecture_id\":\"123e4567-e89b-12d3-a456-426614174001\"}"
    test_endpoint "POST" "/progress/lectures/complete" "$complete_data" 400 "Complete lecture (expected to fail)" "auth"
fi

# 6. Analytics Endpoints
echo -e "${BLUE}6. Analytics Endpoints${NC}"

if [ -n "$TOKEN" ]; then
    test_endpoint "GET" "/analytics/user" "" 200 "Get user analytics" "auth"
fi

# 7. Video Service Endpoints (through API Gateway)
echo -e "${BLUE}7. Video Service Endpoints${NC}"

# Direct video service health check
video_health=$(curl -s http://localhost:8084/health)
echo -e "${YELLOW}Testing: Video service health (direct)${NC}"
echo "  GET http://localhost:8084/health"
echo "  Response: $video_health"
echo ""

# 8. Bucket Service Endpoints
echo -e "${BLUE}8. Bucket Service Endpoints${NC}"

# Direct bucket service health check
bucket_health=$(curl -s http://localhost:8085/health)
echo -e "${YELLOW}Testing: Bucket service health (direct)${NC}"
echo "  GET http://localhost:8085/health"
echo "  Response: $bucket_health"
echo ""

# 9. Forum Service Endpoints
echo -e "${BLUE}9. Forum Service Endpoints${NC}"

# Direct forum service health check
forum_health=$(curl -s http://localhost:8087/health)
echo -e "${YELLOW}Testing: Forum service health (direct)${NC}"
echo "  GET http://localhost:8087/health"
echo "  Response: $forum_health"
echo ""

# 10. Chatbot Service Endpoints
echo -e "${BLUE}10. Chatbot Service Endpoints${NC}"

# Direct chatbot service health check
chatbot_health=$(curl -s http://localhost:8086/health)
echo -e "${YELLOW}Testing: Chatbot service health (direct)${NC}"
echo "  GET http://localhost:8086/health"
echo "  Response: $chatbot_health"
echo ""

# 11. Payment Service Endpoints
echo -e "${BLUE}11. Payment Service Endpoints${NC}"

# Direct payment service health check
payment_health=$(curl -s http://localhost:8088/health)
echo -e "${YELLOW}Testing: Payment service health (direct)${NC}"
echo "  GET http://localhost:8088/health"
echo "  Response: $payment_health"
echo ""

# 12. Error Handling Tests
echo -e "${BLUE}12. Error Handling Tests${NC}"

test_endpoint "GET" "/nonexistent" "" 404 "Non-existent endpoint"
test_endpoint "POST" "/auth/login" \
    "{\"invalid\":\"json\"}" \
    400 "Invalid login data"
test_endpoint "GET" "/auth/profile" "" 401 "Unauthorized access (no token)"

# Summary
echo -e "${BLUE}=== Testing Complete ===${NC}"
echo ""
echo "Notes:"
echo "- Some endpoints may fail if they require specific data setup"
echo "- Authentication-required endpoints need valid tokens"
echo "- Resource-specific endpoints need valid IDs"
echo "- This script tests the API structure and basic functionality"
echo ""
echo "For production use, ensure all environment variables are properly set:"
echo "- Database connections"
echo "- JWT secrets"
echo "- External service credentials (OAuth, payment gateways, etc.)"