#!/usr/bin/env python3
"""
Batch add Swagger annotations to handler files
Usage: python3 batch-add-swagger.py
"""

# Annotations for Progress Handler
progress_annotations = {
    "UpdateProgress": """// UpdateProgress godoc
// @Summary      Update learning progress
// @Description  Update user's progress for a lecture
// @Tags         Progress
// @Accept       json
// @Produce      json
// @Param        request body object true "Progress data"
// @Success      200 {object} APIResponse "Progress updated"
// @Failure      400 {object} APIResponse "Invalid request"
// @Security     BearerAuth
// @Router       /progress/update [post]""",
    
    "GetProgress": """// GetProgress godoc
// @Summary      Get lecture progress
// @Description  Get user's progress for a specific lecture
// @Tags         Progress
// @Accept       json
// @Produce      json
// @Param        course_id path string true "Course ID"
// @Param        lecture_id path string true "Lecture ID"
// @Success      200 {object} APIResponse "Progress data"
// @Failure      404 {object} APIResponse "Not found"
// @Security     BearerAuth
// @Router       /progress/courses/{course_id}/lectures/{lecture_id} [get]""",
    
    "MarkLectureComplete": """// MarkLectureComplete godoc
// @Summary      Mark lecture as complete
// @Description  Mark a lecture as completed by the user
// @Tags         Progress
// @Accept       json
// @Produce      json
// @Param        request body object true "Completion data"
// @Success      200 {object} APIResponse "Lecture marked complete"
// @Security     BearerAuth
// @Router       /progress/lectures/complete [post]"""
}

# Annotations for Payment Handler
payment_annotations = {
    "CreateStripePaymentIntent": """// CreateStripePaymentIntent godoc
// @Summary      Create payment intent
// @Description  Create a Stripe payment intent for course purchase
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        request body object true "Payment details"
// @Success      201 {object} APIResponse "Payment intent created"
// @Failure      400 {object} APIResponse "Invalid request"
// @Security     BearerAuth
// @Router       /payments/stripe/payment-intents [post]""",
    
    "ListStripeTransactions": """// ListStripeTransactions godoc
// @Summary      List transactions
// @Description  Get list of user's payment transactions
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Success      200 {object} APIResponse "List of transactions"
// @Security     BearerAuth
// @Router       /payments/stripe/transactions [get]""",
    
    "PurchaseCourse": """// PurchaseCourse godoc
// @Summary      Purchase course
// @Description  Purchase a course using stored payment method
// @Tags         Payments
// @Accept       json
// @Produce      json
// @Param        course_id path string true "Course ID"
// @Param        request body object true "Payment method details"
// @Success      200 {object} APIResponse "Purchase successful"
// @Failure      402 {object} APIResponse "Payment required"
// @Security     BearerAuth
// @Router       /payments/purchase/course/{course_id} [post]"""
}

print("✅ Swagger annotation templates ready!")
print(f"📊 Progress annotations: {len(progress_annotations)}")
print(f"💳 Payment annotations: {len(payment_annotations)}")
print("\n💡 Apply these annotations manually to handler files")
print("📝 Then run: swag init -g cmd/main.go -o docs")
