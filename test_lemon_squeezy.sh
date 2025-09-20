#!/bin/bash

# Test script for Lemon Squeezy integration
echo "Testing Lemon Squeezy Integration"
echo "=================================="

# Load environment variables
source .env

# Test variables
USER_ID="test-user-123"
COURSE_ID="test-course-456"
BASE_URL="http://localhost:8088"

echo "Configuration:"
echo "API Key: ${LEMON_SQUEEZY_API_KEY:0:20}..."
echo "Store ID: $LEMON_SQUEEZY_STORE_ID"
echo "Product ID: $LEMON_SQUEEZY_PRODUCT_ID"
echo "Variant ID: $LEMON_SQUEEZY_VARIANT_ID"
echo "Webhook Secret: ${LEMON_SQUEEZY_WEBHOOK_SECRET:0:10}..."
echo "Webhook URL: $LEMON_SQUEEZY_WEBHOOK_URL"
echo ""

# Test 1: Health check
echo "Test 1: Health Check"
echo "===================="
curl -s "$BASE_URL/health" | jq '.' || echo "Health check failed"
echo ""

# Test 2: Get products
echo "Test 2: Get Products"
echo "==================="
curl -s -H "X-User-ID: $USER_ID" -H "X-User-Role: student" \
     "$BASE_URL/api/v1/lemonsqueezy/products" | jq '.' || echo "Get products failed"
echo ""

# Test 3: Get variants
echo "Test 3: Get Variants"
echo "==================="
curl -s -H "X-User-ID: $USER_ID" -H "X-User-Role: student" \
     "$BASE_URL/api/v1/lemonsqueezy/variants?product_id=$LEMON_SQUEEZY_PRODUCT_ID" | jq '.' || echo "Get variants failed"
echo ""

# Test 4: Create checkout
echo "Test 4: Create Checkout"
echo "======================="
CHECKOUT_DATA='{
  "variant_id": "'$LEMON_SQUEEZY_VARIANT_ID'",
  "custom_data": {
    "course_name": "Test Course",
    "user_email": "test@example.com"
  },
  "success_url": "http://localhost:3000/payment/success",
  "cancel_url": "http://localhost:3000/payment/cancel"
}'

echo "Checkout data: $CHECKOUT_DATA"
echo ""

CHECKOUT_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "X-User-ID: $USER_ID" \
  -H "X-User-Role: student" \
  -d "$CHECKOUT_DATA" \
  "$BASE_URL/api/v1/lemonsqueezy/checkout/course/$COURSE_ID")

echo "Checkout response:"
echo "$CHECKOUT_RESPONSE" | jq '.' || echo "Create checkout failed"
echo ""

# Test 5: Test webhook signature verification
echo "Test 5: Test Webhook Signature"
echo "=============================="
WEBHOOK_DATA='{
  "meta": {
    "event_name": "order_created",
    "custom_data": {
      "user_id": "'$USER_ID'",
      "course_id": "'$COURSE_ID'"
    }
  },
  "data": {
    "id": "test-order-123",
    "type": "orders",
    "attributes": {
      "store_id": '$LEMON_SQUEEZY_STORE_ID',
      "customer_id": 12345,
      "identifier": "test-order-123",
      "order_number": 1001,
      "user_name": "Test User",
      "user_email": "test@example.com",
      "currency": "USD",
      "currency_rate": "1.00000000",
      "subtotal": 2999,
      "discount_total": 0,
      "tax": 0,
      "total": 2999,
      "subtotal_usd": 2999,
      "discount_total_usd": 0,
      "tax_usd": 0,
      "total_usd": 2999,
      "tax_name": null,
      "tax_rate": "0.00",
      "status": "paid",
      "status_formatted": "Paid",
      "refunded": false,
      "refunded_at": null,
      "subtotal_formatted": "$29.99",
      "discount_total_formatted": "$0.00",
      "tax_formatted": "$0.00",
      "total_formatted": "$29.99",
      "created_at": "2024-01-01T00:00:00.000000Z",
      "updated_at": "2024-01-01T00:00:00.000000Z"
    }
  }
}'

# Create signature for webhook test
SIGNATURE=$(echo -n "$WEBHOOK_DATA" | openssl dgst -sha256 -hmac "$LEMON_SQUEEZY_WEBHOOK_SECRET" -hex | sed 's/^.* //')

echo "Testing webhook with signature: $SIGNATURE"
echo ""

WEBHOOK_RESPONSE=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -H "X-Signature: $SIGNATURE" \
  -d "$WEBHOOK_DATA" \
  "$BASE_URL/api/v1/payments/lemonsqueezy/webhook")

echo "Webhook response status: $?"
echo "Webhook response: $WEBHOOK_RESPONSE"
echo ""

echo "Lemon Squeezy Integration Test Complete!"
echo "========================================"