package service

import (
	"context"
	"fmt"
	"log"
	"time"

	progresspb "github.com/study-platform/payment-service/github.com/study-platform/progress-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProgressClient struct {
	client progresspb.ProgressServiceClient
	conn   *grpc.ClientConn
}

func NewProgressClient(serviceURL string) (*ProgressClient, error) {
	conn, err := grpc.Dial(serviceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to progress service: %w", err)
	}

	client := progresspb.NewProgressServiceClient(conn)

	return &ProgressClient{
		client: client,
		conn:   conn,
	}, nil
}

func (pc *ProgressClient) Close() error {
	return pc.conn.Close()
}

func (pc *ProgressClient) CreateEnrollment(ctx context.Context, userID, courseID string) (*progresspb.Enrollment, error) {
	req := &progresspb.CreateEnrollmentRequest{
		UserId:   userID,
		CourseId: courseID,
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := pc.client.CreateEnrollment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to create enrollment: %w", err)
	}

	log.Printf("Successfully created enrollment: %s for user %s in course %s", 
		resp.Enrollment.Id, userID, courseID)

	return resp.Enrollment, nil
}

func (pc *ProgressClient) GetEnrollment(ctx context.Context, userID, courseID string) (*progresspb.Enrollment, error) {
	req := &progresspb.GetEnrollmentRequest{
		UserId:   userID,
		CourseId: courseID,
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := pc.client.GetEnrollment(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to get enrollment: %w", err)
	}

	return resp.Enrollment, nil
}

func (pc *ProgressClient) ListUserEnrollments(ctx context.Context, userID string, page, pageSize int32) ([]*progresspb.Enrollment, error) {
	req := &progresspb.ListEnrollmentsRequest{
		UserId:   userID,
		Page:     page,
		PageSize: pageSize,
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := pc.client.ListEnrollments(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to list enrollments: %w", err)
	}

	return resp.Enrollments, nil
}

func (pc *ProgressClient) UpdateEnrollmentStatus(ctx context.Context, userID, courseID, status string) (*progresspb.Enrollment, error) {
	req := &progresspb.UpdateEnrollmentStatusRequest{
		UserId:   userID,
		CourseId: courseID,
		Status:   status,
	}

	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	resp, err := pc.client.UpdateEnrollmentStatus(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to update enrollment status: %w", err)
	}

	log.Printf("Successfully updated enrollment status to %s for user %s in course %s",
		status, userID, courseID)

	return resp.Enrollment, nil
}

// CompleteEnrollmentPayment completes a pending enrollment after successful payment verification
func (pc *ProgressClient) CompleteEnrollmentPayment(ctx context.Context, userID, courseID, orderID string, amount float64, currency string) error {
	// First, check if there's already an enrollment
	_, err := pc.GetEnrollment(ctx, userID, courseID)
	if err != nil {
		// No existing enrollment, create new one with payment info
		_, err = pc.CreateEnrollment(ctx, userID, courseID)
		if err != nil {
			return fmt.Errorf("failed to create paid enrollment: %w", err)
		}
		log.Printf("Created new paid enrollment for user %s in course %s (order: %s)", userID, courseID, orderID)
		return nil
	}

	// Update existing enrollment to active/paid status
	_, err = pc.UpdateEnrollmentStatus(ctx, userID, courseID, "enrolled")
	if err != nil {
		return fmt.Errorf("failed to activate enrollment after payment: %w", err)
	}

	log.Printf("Activated enrollment for user %s in course %s after payment verification (order: %s, amount: %.2f %s)",
		userID, courseID, orderID, amount, currency)

	return nil
}