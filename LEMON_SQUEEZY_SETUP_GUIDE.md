# Lemon Squeezy Integration Setup Guide

This guide provides step-by-step instructions for setting up Lemon Squeezy integration with the Study Platform for course payments.

## Table of Contents

1. [Lemon Squeezy Account Setup](#lemon-squeezy-account-setup)
2. [Store Configuration](#store-configuration)
3. [Product and Variant Setup](#product-and-variant-setup)
4. [Webhook Configuration](#webhook-configuration)
5. [Environment Variables](#environment-variables)
6. [Database Migration](#database-migration)
7. [API Integration](#api-integration)
8. [Testing](#testing)
9. [Troubleshooting](#troubleshooting)

---

## Lemon Squeezy Account Setup

### 1. Create Account
1. Visit [lemonsqueezy.com](https://lemonsqueezy.com)
2. Sign up for a new account or log in to existing account
3. Complete the verification process

### 2. Create Store
1. Go to **Settings > Store**
2. Click **Create Store**
3. Fill in store details:
   - **Store Name**: Your platform name (e.g., "StudyPlatform Courses")
   - **Store Slug**: URL-friendly name (e.g., "studyplatform")
   - **Currency**: USD (or your preferred currency)
   - **Country**: Your business location
4. Complete tax information if required
5. Save store settings

### 3. Get Store ID
1. Go to **Settings > Store**
2. Note the **Store ID** (displayed in the URL or store settings)
3. Save this for environment configuration

---

## Store Configuration

### Payment Methods
1. Go to **Settings > Payment Methods**
2. Enable desired payment methods:
   - Credit/Debit Cards (enabled by default)
   - PayPal (optional)
   - Apple Pay (optional)
   - Google Pay (optional)

### Tax Settings
1. Go to **Settings > Tax**
2. Configure tax settings based on your business requirements
3. Set up tax rates for different regions if applicable

### Email Templates
1. Go to **Settings > Emails**
2. Customize email templates for:
   - Purchase confirmations
   - Receipt emails
   - Refund notifications

---

## Product and Variant Setup

### Creating Products for Courses

For each paid course in your platform, you need to create a corresponding product in Lemon Squeezy:

#### 1. Create Product
1. Go to **Products > Add Product**
2. Fill in product details:
   - **Name**: Course title (e.g., "Advanced JavaScript Masterclass")
   - **Description**: Course description
   - **Category**: "Digital Products" or "Education"
   - **Status**: "Published"

#### 2. Create Variant
1. Within the product, click **Add Variant**
2. Configure variant:
   - **Name**: "Standard Access" (or course access type)
   - **Price**: Course price in your currency
   - **Billing**: "One-time" (for course purchases)
   - **Status**: "Published"

#### 3. Note IDs
1. Save the **Product ID** and **Variant ID**
2. These will be used to link courses in your database

### Bulk Setup for Multiple Courses
For multiple courses, you can:
1. Create products manually for each course
2. Use Lemon Squeezy API to create products programmatically
3. Use CSV import if available

---

## Webhook Configuration

### 1. Create Webhook
1. Go to **Settings > Webhooks**
2. Click **Create Webhook**
3. Configure webhook:
   - **Endpoint URL**: `https://yourdomain.com/api/v1/payments/lemonsqueezy/webhook`
   - **Secret**: Generate a secure random string (save this for environment config)
   - **Events**: Select the following events:
     - `order_created`
     - `order_refunded`
     - `subscription_payment_success` (if using subscriptions)
     - `subscription_payment_failed` (if using subscriptions)

### 2. Test Webhook
1. Use the webhook testing feature in Lemon Squeezy dashboard
2. Verify your endpoint receives and processes webhooks correctly
3. Check webhook logs for any errors

### Webhook Security
- Always verify webhook signatures using the secret
- Store the webhook secret securely in environment variables
- Implement idempotency to handle duplicate webhooks
- Log webhook events for monitoring and debugging

---

## Environment Variables

Add the following environment variables to your `.env` file:

```bash
# Lemon Squeezy Configuration
LEMONSQUEEZY_API_KEY=your_api_key_here
LEMONSQUEEZY_STORE_ID=your_store_id_here
LEMONSQUEEZY_WEBHOOK_SECRET=your_webhook_secret_here
LEMONSQUEEZY_ENVIRONMENT=test  # or "production"

# Payment Service Configuration
PAYMENT_SERVICE_PORT=8088
PAYMENT_DEFAULT_CURRENCY=USD
```

### Getting API Key
1. Go to **Settings > API**
2. Click **Create API Key**
3. Copy the generated API key
4. Store securely in your environment variables

### Environment Setup
- **Test Mode**: Use test API keys and store for development
- **Production Mode**: Use production API keys and store for live payments

---

## Database Migration

Run the Lemon Squeezy database migration:

```bash
# Navigate to Study-Platform directory
cd Study-Platform

# Run the migration
go run main.go migrate up
```

This will:
- Add Lemon Squeezy columns to existing tables
- Create webhook events table
- Create Lemon Squeezy products/variants cache tables
- Update constraints for Lemon Squeezy-only payments

### Manual Migration
If automatic migration fails, run the SQL manually:

```sql
-- Run the contents of migrations/019_lemon_squeezy_integration.up.sql
```

---

## API Integration

### Course Setup
For each course that should be paid via Lemon Squeezy:

1. **Update Course in Database**:
```sql
UPDATE courses
SET
  is_paid = TRUE,
  lemon_squeezy_product_id = 'your_product_id',
  lemon_squeezy_variant_id = 'your_variant_id'
WHERE id = 'course_uuid';
```

2. **API Endpoints**:
   - `POST /api/v1/payments/courses/{course_id}/checkout` - Create checkout
   - `POST /api/v1/payments/lemonsqueezy/webhook` - Handle webhooks
   - `GET /api/v1/payments/verify/{order_id}` - Verify payment

### Payment Flow
1. User clicks "Buy Course"
2. Frontend calls checkout API with course ID
3. Backend creates Lemon Squeezy checkout session
4. User redirects to Lemon Squeezy payment page
5. After payment, Lemon Squeezy sends webhook
6. Backend processes webhook and enrolls user
7. User gets access to course

---

## Testing

### Test Mode Setup
1. Use test API keys in environment
2. Create test products in Lemon Squeezy dashboard
3. Use test payment methods provided by Lemon Squeezy

### Test Payment Methods
Lemon Squeezy provides test card numbers:
- **Success**: `4242424242424242`
- **Decline**: `4000000000000002`
- **Insufficient Funds**: `4000000000009995`

### Testing Checklist
- [ ] Create test checkout session
- [ ] Complete test payment
- [ ] Verify webhook delivery
- [ ] Confirm user enrollment
- [ ] Test refund process
- [ ] Verify error handling

### Integration Testing
```bash
# Test checkout creation
curl -X POST http://localhost:8088/api/v1/payments/courses/{course_id}/checkout \
  -H "Authorization: Bearer {jwt_token}" \
  -H "Content-Type: application/json" \
  -d '{"variant_id": "your_variant_id"}'

# Test webhook (using test payload)
curl -X POST http://localhost:8088/api/v1/payments/lemonsqueezy/webhook \
  -H "X-Signature: test_signature" \
  -H "Content-Type: application/json" \
  -d @test_webhook_payload.json
```

---

## Troubleshooting

### Common Issues

#### 1. Webhook Not Receiving Events
**Symptoms**: Payments complete but users not enrolled
**Solutions**:
- Check webhook URL is accessible from internet
- Verify webhook secret matches environment variable
- Check webhook event logs in Lemon Squeezy dashboard
- Ensure webhook endpoint returns 200 status

#### 2. Invalid API Key
**Symptoms**: 401 Unauthorized errors
**Solutions**:
- Verify API key is correct and active
- Check API key has proper permissions
- Ensure using correct environment (test vs production)

#### 3. Product/Variant Not Found
**Symptoms**: Checkout creation fails
**Solutions**:
- Verify product and variant IDs are correct
- Check product/variant status is "published"
- Ensure product belongs to correct store

#### 4. Signature Verification Failed
**Symptoms**: Webhook signature verification errors
**Solutions**:
- Check webhook secret matches configuration
- Verify signature calculation algorithm
- Ensure request body is not modified before verification

### Debug Mode
Enable debug logging by setting:
```bash
LOG_LEVEL=debug
```

### Monitoring
Monitor the following:
- Webhook delivery success rate
- Payment completion rate
- User enrollment rate after payment
- Error rates in payment processing

### Support Resources
- [Lemon Squeezy Documentation](https://docs.lemonsqueezy.com/)
- [API Reference](https://docs.lemonsqueezy.com/api)
- [Webhook Guide](https://docs.lemonsqueezy.com/guides/webhooks)
- [Testing Guide](https://docs.lemonsqueezy.com/guides/testing)

---

## Production Checklist

Before going live:

- [ ] Switch to production API keys
- [ ] Update webhook URL to production domain
- [ ] Test payment flow end-to-end
- [ ] Verify webhook security
- [ ] Set up monitoring and alerting
- [ ] Configure proper error handling
- [ ] Test refund process
- [ ] Verify tax calculations
- [ ] Check compliance requirements
- [ ] Set up backup webhook endpoints (optional)

## Security Best Practices

1. **API Key Security**:
   - Store API keys in environment variables only
   - Use different keys for development and production
   - Rotate keys regularly
   - Never commit keys to version control

2. **Webhook Security**:
   - Always verify webhook signatures
   - Use HTTPS for webhook endpoints
   - Implement proper authentication
   - Log webhook events for audit trail

3. **Error Handling**:
   - Implement proper error responses
   - Log errors for monitoring
   - Handle edge cases gracefully
   - Provide user-friendly error messages

4. **Data Protection**:
   - Follow PCI DSS guidelines
   - Minimize stored payment data
   - Encrypt sensitive information
   - Implement proper access controls

---

This completes the Lemon Squeezy integration setup for the Study Platform. Follow this guide step-by-step to ensure proper configuration and security.