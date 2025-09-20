package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/study-platform/course-service/internal/model"
	"github.com/study-platform/course-service/internal/repository"
)

// LemonSqueezySync handles synchronization between courses and LemonSqueezy products
type LemonSqueezySync struct {
	courseRepo   *repository.CourseRepository
	lectureRepo  *repository.LectureRepository
}

type LemonSqueezyProductMapping struct {
	CourseID          uuid.UUID `json:"course_id"`
	ProductID         string    `json:"product_id"`         // LemonSqueezy product ID
	VariantID         string    `json:"variant_id"`         // LemonSqueezy variant ID
	RequiresPayment   bool      `json:"requires_payment"`
	Price             float64   `json:"price"`
	Currency          string    `json:"currency"`
	ProductName       string    `json:"product_name"`       // From LemonSqueezy
	ProductStatus     string    `json:"product_status"`     // draft, published
}

type CourseCheckoutRequest struct {
	CourseID      uuid.UUID `json:"course_id" binding:"required"`
	UserID        uuid.UUID `json:"user_id" binding:"required"`
	SuccessURL    string    `json:"success_url,omitempty"`
	CancelURL     string    `json:"cancel_url,omitempty"`
	CustomerEmail string    `json:"customer_email,omitempty"`
	CustomerName  string    `json:"customer_name,omitempty"`
}

type CourseCheckoutResponse struct {
	CheckoutURL    string                     `json:"checkout_url"`
	CheckoutID     string                     `json:"checkout_id"`
	CourseInfo     *CourseCheckoutInfo        `json:"course_info"`
	ProductMapping *LemonSqueezyProductMapping `json:"product_mapping"`
	ExpiresAt      time.Time                  `json:"expires_at"`
}

type CourseCheckoutInfo struct {
	ID          uuid.UUID `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Currency    string    `json:"currency"`
	IsPaid      bool      `json:"is_paid"`
}

func NewLemonSqueezySync(courseRepo *repository.CourseRepository, lectureRepo *repository.LectureRepository) *LemonSqueezySync {
	return &LemonSqueezySync{
		courseRepo:  courseRepo,
		lectureRepo: lectureRepo,
	}
}

// LinkCourseToLemonSqueezyProduct links a course to an existing LemonSqueezy product
func (ls *LemonSqueezySync) LinkCourseToLemonSqueezyProduct(ctx context.Context, courseID uuid.UUID, productID, variantID string) error {
	// Validate course exists
	course, err := ls.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("course not found: %w", err)
	}

	// Validate that the course is not already linked
	if course.LemonSqueezyProductID != nil && *course.LemonSqueezyProductID != "" {
		return fmt.Errorf("course is already linked to LemonSqueezy product %s", *course.LemonSqueezyProductID)
	}

	// Update course with LemonSqueezy product information
	course.LemonSqueezyProductID = &productID
	course.LemonSqueezyVariantID = &variantID
	course.IsPaid = true // Linking to LemonSqueezy implies it's paid

	err = ls.courseRepo.Update(ctx, course)
	if err != nil {
		return fmt.Errorf("failed to update course with LemonSqueezy product mapping: %w", err)
	}

	log.Printf("Successfully linked course %s to LemonSqueezy product %s (variant %s)",
		courseID.String(), productID, variantID)

	return nil
}

// UnlinkCourseFromLemonSqueezy removes LemonSqueezy product linkage from a course
func (ls *LemonSqueezySync) UnlinkCourseFromLemonSqueezy(ctx context.Context, courseID uuid.UUID) error {
	course, err := ls.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("course not found: %w", err)
	}

	course.LemonSqueezyProductID = nil
	course.LemonSqueezyVariantID = nil
	course.IsPaid = false

	err = ls.courseRepo.Update(ctx, course)
	if err != nil {
		return fmt.Errorf("failed to unlink course from LemonSqueezy: %w", err)
	}

	log.Printf("Successfully unlinked course %s from LemonSqueezy", courseID.String())
	return nil
}

// GetCourseProductMapping retrieves the LemonSqueezy product mapping for a course
func (ls *LemonSqueezySync) GetCourseProductMapping(ctx context.Context, courseID uuid.UUID) (*LemonSqueezyProductMapping, error) {
	course, err := ls.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	if course.LemonSqueezyProductID == nil || *course.LemonSqueezyProductID == "" {
		return nil, fmt.Errorf("course is not linked to any LemonSqueezy product")
	}

	mapping := &LemonSqueezyProductMapping{
		CourseID:        course.ID,
		ProductID:       *course.LemonSqueezyProductID,
		VariantID:       *course.LemonSqueezyVariantID,
		RequiresPayment: course.IsPaid,
		Price:           course.Price,
		Currency:        course.Currency,
		ProductStatus:   "published", // Assume published if linked
	}

	return mapping, nil
}

// ValidateCourseForPurchase validates that a course can be purchased
func (ls *LemonSqueezySync) ValidateCourseForPurchase(ctx context.Context, courseID uuid.UUID) (*CourseCheckoutInfo, error) {
	course, err := ls.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return nil, fmt.Errorf("course not found: %w", err)
	}

	// Check if course is published
	if course.Status != "published" {
		return nil, fmt.Errorf("course is not available for purchase (status: %s)", course.Status)
	}

	// Check if course requires payment
	if !course.IsPaid {
		return nil, fmt.Errorf("course is free and does not require payment")
	}

	// Check if course is linked to LemonSqueezy product
	if course.LemonSqueezyProductID == nil || *course.LemonSqueezyProductID == "" {
		return nil, fmt.Errorf("course is not properly configured for payment (no LemonSqueezy product linked)")
	}

	courseInfo := &CourseCheckoutInfo{
		ID:          course.ID,
		Title:       course.Title,
		Description: course.Description,
		Price:       course.Price,
		Currency:    course.Currency,
		IsPaid:      course.IsPaid,
	}

	return courseInfo, nil
}

// SyncCourseWithLemonSqueezyProduct updates course information to match LemonSqueezy product
// This is useful when product details change in LemonSqueezy dashboard
func (ls *LemonSqueezySync) SyncCourseWithLemonSqueezyProduct(ctx context.Context, courseID uuid.UUID, productData map[string]interface{}) error {
	course, err := ls.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("course not found: %w", err)
	}

	// Extract relevant fields from LemonSqueezy product data
	if name, ok := productData["name"].(string); ok && name != "" {
		// Optionally sync product name with course title
		// course.Title = name (uncomment if you want to sync names)
	}

	if status, ok := productData["status"].(string); ok {
		// Update course status based on LemonSqueezy product status
		switch status {
		case "published":
			course.Status = "published"
		case "draft":
			course.Status = "draft"
		}
	}

	// Sync price if variant information is provided
	if variants, ok := productData["variants"].([]interface{}); ok && len(variants) > 0 {
		if variantData, ok := variants[0].(map[string]interface{}); ok {
			if price, ok := variantData["price"].(float64); ok {
				course.Price = price / 100 // LemonSqueezy prices are in cents
			}
			if currency, ok := variantData["currency"].(string); ok {
				course.Currency = currency
			}
		}
	}

	err = ls.courseRepo.Update(ctx, course)
	if err != nil {
		return fmt.Errorf("failed to sync course with LemonSqueezy product: %w", err)
	}

	log.Printf("Successfully synced course %s with LemonSqueezy product data", courseID.String())
	return nil
}

// ListCoursesWithLemonSqueezyProducts returns all courses that are linked to LemonSqueezy products
func (ls *LemonSqueezySync) ListCoursesWithLemonSqueezyProducts(ctx context.Context) ([]*LemonSqueezyProductMapping, error) {
	// This would require a custom query in the repository
	// For now, returning empty slice as placeholder
	return []*LemonSqueezyProductMapping{}, nil
}

// GetCourseByLemonSqueezyProduct finds a course by its LemonSqueezy product ID
func (ls *LemonSqueezySync) GetCourseByLemonSqueezyProduct(ctx context.Context, productID string) (*model.Course, error) {
	// TODO: Implement this method when the correct repository interface is available
	return nil, fmt.Errorf("GetCourseByLemonSqueezyProduct not implemented yet")
}

// CreateFreeCourseEnrollment handles enrollment for free courses
func (ls *LemonSqueezySync) CreateFreeCourseEnrollment(ctx context.Context, courseID, userID uuid.UUID) error {
	course, err := ls.courseRepo.GetByID(ctx, courseID)
	if err != nil {
		return fmt.Errorf("course not found: %w", err)
	}

	if course.IsPaid {
		return fmt.Errorf("course requires payment and cannot be enrolled for free")
	}

	// For free courses, we would typically call the progress service to create enrollment
	// This is a placeholder for that integration
	log.Printf("Creating free enrollment for user %s in course %s", userID.String(), courseID.String())

	return nil
}