#!/bin/bash

# Quick Network Monitoring API Test
# Tests key network monitoring endpoints

set -e

# Configuration
VIDEO_SERVICE_URL="http://localhost:8084/api/videos"
SESSION_ID="quick-test-session-$(date +%s)"

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

test_endpoint() {
    local method=$1
    local endpoint=$2
    local data=$3
    local description=$4

    print_info "Testing: $description"

    if [ "$method" = "GET" ]; then
        response=$(curl -s -w "\n%{http_code}" "$VIDEO_SERVICE_URL$endpoint" 2>/dev/null)
    else
        response=$(curl -s -w "\n%{http_code}" -X "$method" -H "Content-Type: application/json" -d "$data" "$VIDEO_SERVICE_URL$endpoint" 2>/dev/null)
    fi

    http_code=$(echo "$response" | tail -n1)
    response_body=$(echo "$response" | head -n -1)

    if [ "$http_code" -ge 200 ] && [ "$http_code" -lt 300 ]; then
        print_success "✓ $description (HTTP $http_code)"
        echo "$response_body" | head -c 200
        echo "..."
    else
        print_error "✗ $description (HTTP $http_code)"
    fi
    echo "----------------------------------------"
}

echo "========================================"
echo "  Quick Network Monitoring Test"
echo "========================================"
echo "Service URL: $VIDEO_SERVICE_URL"
echo "Session ID: $SESSION_ID"
echo "========================================"

# Test 1: Health Check
test_endpoint "GET" "/../../health" "" "Health Check"

# Test 2: Network Status Update (no auth needed for direct access)
network_update='{
    "bandwidth_mbps": 8.5,
    "latency_ms": 65,
    "packet_loss": 0.02,
    "connection_type": "wifi",
    "buffer_health": 12,
    "current_time": 150,
    "current_quality": "720p"
}'
test_endpoint "POST" "/sessions/$SESSION_ID/network" "$network_update" "Network Status Update"

# Test 3: WebSocket Stats (public endpoint)
test_endpoint "GET" "/ws/stats" "" "WebSocket Statistics"

# Test 4: Quality Recommendation (might need auth but worth testing)
test_endpoint "GET" "/sessions/$SESSION_ID/quality-recommendation?connection_type=wifi" "" "Quality Recommendation"

# Test 5: Network Diagnostics (might need auth but worth testing)
diagnostics_request='{
    "test_duration_seconds": 5,
    "test_type": "bandwidth"
}'
test_endpoint "POST" "/sessions/$SESSION_ID/network-diagnostics" "$diagnostics_request" "Network Diagnostics"

echo "========================================"
echo "  Quick Test Complete"
echo "========================================"
echo "Note: Some endpoints may require authentication"
echo "For full testing, use test_network_monitoring_apis.sh"