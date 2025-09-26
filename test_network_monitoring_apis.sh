#!/bin/bash

# Comprehensive Network Monitoring API Test Script
# This script tests all network monitoring endpoints for the video service

set -e

# Configuration
BASE_URL="http://localhost:8084/api/videos"
JWT_TOKEN="your-jwt-token-here"  # Replace with actual JWT token
SESSION_ID="test-session-123"
VIDEO_ID="test-video-456"
USER_ID="test-user-789"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

# Function to make API calls with proper headers
make_request() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4

    print_status "Testing: $description"
    echo "Endpoint: $method $endpoint"

    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -H "Content-Type: application/json" \
            "$BASE_URL$endpoint" 2>/dev/null)
    else
        response=$(curl -s -w "\n%{http_code}" \
            -X "$method" \
            -H "Authorization: Bearer $JWT_TOKEN" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$BASE_URL$endpoint" 2>/dev/null)
    fi

    http_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | head -n -1)

    echo "HTTP Status: $http_code"
    echo "Response: $response_body" | jq . 2>/dev/null || echo "$response_body"

    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        print_success "✓ $description - OK"
    else
        print_error "✗ $description - Failed (HTTP $http_code)"
    fi

    echo "----------------------------------------"
    echo
}

# Check if jq is installed for JSON formatting
if ! command -v jq &> /dev/null; then
    print_warning "jq is not installed. JSON responses will not be formatted."
fi

echo "========================================"
echo "  Network Monitoring API Test Suite"
echo "========================================"
echo "Base URL: $BASE_URL"
echo "Session ID: $SESSION_ID"
echo "Video ID: $VIDEO_ID"
echo "========================================"
echo

# Test 1: Health Check
print_status "Starting API tests..."
make_request "GET" "/../../health" "" "Health Check"

# Test 2: Create a video session first
print_status "Setting up test session..."
session_data='{
    "device_info": {
        "user_agent": "Mozilla/5.0 (Test Browser)",
        "screen_resolution": "1920x1080",
        "connection_type": "wifi"
    }
}'
make_request "POST" "/$VIDEO_ID/sessions" "$session_data" "Create Video Session"

# Test 3: Update Network Status
print_status "Testing network status updates..."
network_update='{
    "bandwidth_mbps": 8.5,
    "latency_ms": 65,
    "packet_loss": 0.02,
    "connection_type": "wifi",
    "buffer_health": 12,
    "current_time": 150,
    "current_quality": "720p"
}'
make_request "POST" "/sessions/$SESSION_ID/network" "$network_update" "Update Network Status"

# Test 4: Get Network Status
make_request "GET" "/sessions/$SESSION_ID/network-status" "" "Get Current Network Status"

# Test 5: Get Quality Recommendation
make_request "GET" "/sessions/$SESSION_ID/quality-recommendation?connection_type=wifi" "" "Get Quality Recommendation"

# Test 6: Get Network Pattern Analysis
make_request "GET" "/sessions/$SESSION_ID/network-pattern?window_minutes=5" "" "Get Network Pattern Analysis"

# Test 7: Network Diagnostics
diagnostics_request='{
    "test_duration_seconds": 10,
    "test_type": "full"
}'
make_request "POST" "/sessions/$SESSION_ID/network-diagnostics" "$diagnostics_request" "Run Network Diagnostics"

# Test 8: Get Session Analytics
make_request "GET" "/sessions/$SESSION_ID/analytics" "" "Get Comprehensive Session Analytics"

# Test 9: Update Bandwidth Estimate
bandwidth_estimate='{
    "estimated_bandwidth_mbps": 12.5,
    "measurement_method": "throughput_test",
    "confidence": 0.85,
    "test_duration_seconds": 5
}'
make_request "POST" "/sessions/$SESSION_ID/bandwidth-estimate" "$bandwidth_estimate" "Update Bandwidth Estimate"

# Test 10: Get Network History
make_request "GET" "/sessions/$SESSION_ID/network-history?hours=1" "" "Get Network Metrics History"

# Test 11: Get Active Network Sessions
make_request "GET" "/network/active-sessions" "" "Get Active Network Sessions"

# Test 12: WebSocket Stats (no auth needed)
make_request "GET" "/ws/stats" "" "Get WebSocket Statistics"

# Test 13: Update Session Progress
progress_update='{
    "current_time": 180,
    "quality": "720p"
}'
make_request "PUT" "/sessions/$SESSION_ID/progress" "$progress_update" "Update Session Progress"

# Test 14: Test WebSocket Connection
print_status "Testing WebSocket Connection..."
echo "WebSocket URL: ws://localhost:8084/api/videos/ws/$SESSION_ID"
print_warning "WebSocket testing requires manual testing or wscat tool"
echo

# Test 15: Test different connection types
print_status "Testing different connection scenarios..."

# 4G Connection
network_4g='{
    "bandwidth_mbps": 5.2,
    "latency_ms": 120,
    "packet_loss": 0.05,
    "connection_type": "4g",
    "buffer_health": 8,
    "current_time": 200,
    "current_quality": "480p"
}'
make_request "POST" "/sessions/$SESSION_ID/network" "$network_4g" "4G Network Update"

# Ethernet Connection
network_ethernet='{
    "bandwidth_mbps": 25.0,
    "latency_ms": 15,
    "packet_loss": 0.001,
    "connection_type": "ethernet",
    "buffer_health": 20,
    "current_time": 220,
    "current_quality": "1080p"
}'
make_request "POST" "/sessions/$SESSION_ID/network" "$network_ethernet" "Ethernet Network Update"

# Poor Connection
network_poor='{
    "bandwidth_mbps": 1.8,
    "latency_ms": 300,
    "packet_loss": 0.08,
    "connection_type": "3g",
    "buffer_health": 3,
    "current_time": 240,
    "current_quality": "360p"
}'
make_request "POST" "/sessions/$SESSION_ID/network" "$network_poor" "Poor Network Update"

# Test 16: Error scenarios
print_status "Testing error scenarios..."

# Invalid session ID
make_request "GET" "/sessions/invalid-session/network-status" "" "Invalid Session ID (Should Fail)"

# Invalid data format
invalid_data='{"invalid": "data"}'
make_request "POST" "/sessions/$SESSION_ID/network" "$invalid_data" "Invalid Network Data (Should Fail)"

# Missing required fields
incomplete_data='{
    "bandwidth_mbps": 5.0
}'
make_request "POST" "/sessions/$SESSION_ID/network" "$incomplete_data" "Incomplete Network Data"

# Test 17: Performance testing with rapid updates
print_status "Performance testing with rapid updates..."
for i in {1..5}; do
    bandwidth=$((RANDOM % 20 + 1))
    latency=$((RANDOM % 200 + 20))
    rapid_update="{
        \"bandwidth_mbps\": $bandwidth,
        \"latency_ms\": $latency,
        \"packet_loss\": 0.01,
        \"connection_type\": \"wifi\",
        \"buffer_health\": 10,
        \"current_time\": $((250 + i * 10)),
        \"current_quality\": \"720p\"
    }"
    make_request "POST" "/sessions/$SESSION_ID/network" "$rapid_update" "Rapid Update #$i"
    sleep 1
done

echo "========================================"
echo "  Network Monitoring API Tests Complete"
echo "========================================"
echo

print_status "Summary of tested endpoints:"
echo "1. Health Check"
echo "2. Create Video Session"
echo "3. Update Network Status"
echo "4. Get Network Status"
echo "5. Get Quality Recommendation"
echo "6. Get Network Pattern Analysis"
echo "7. Run Network Diagnostics"
echo "8. Get Session Analytics"
echo "9. Update Bandwidth Estimate"
echo "10. Get Network History"
echo "11. Get Active Sessions"
echo "12. WebSocket Stats"
echo "13. Update Session Progress"
echo "14. Different Connection Types Testing"
echo "15. Error Scenarios"
echo "16. Performance Testing"

print_success "All API tests completed!"
print_warning "Note: Some tests may fail if the service is not running or if authentication is required."
print_warning "Make sure to:"
print_warning "1. Update JWT_TOKEN with a valid token"
print_warning "2. Ensure the video service is running on port 8084"
print_warning "3. Have Redis and PostgreSQL running"
print_warning "4. Run database migrations"

echo
echo "To run this script:"
echo "chmod +x test_network_monitoring_apis.sh"
echo "./test_network_monitoring_apis.sh"