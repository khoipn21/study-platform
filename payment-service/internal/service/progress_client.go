package service

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "payment-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ProgressClient struct {
	client pb.ProgressServiceClient
	conn   *grpc.ClientConn
}

func NewProgressClient(serviceURL string) (*ProgressClient, error) {
	conn, err := grpc.Dial(serviceURL, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to progress service: %w", err)
	}

	client := pb.NewProgressServiceClient(conn)

	return &ProgressClient{
		client: client,
		conn:   conn,
	}, nil
}

func (pc *ProgressClient) Close() error {
	return pc.conn.Close()
}

func (pc *ProgressClient) CreateEnrollment(ctx context.Context, userID, courseID string) (*pb.Enrollment, error) {
	req := &pb.CreateEnrollmentRequest{
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

func (pc *ProgressClient) GetEnrollment(ctx context.Context, userID, courseID string) (*pb.Enrollment, error) {
	req := &pb.GetEnrollmentRequest{
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

func (pc *ProgressClient) ListUserEnrollments(ctx context.Context, userID string, page, pageSize int32) ([]*pb.Enrollment, error) {
	req := &pb.ListEnrollmentsRequest{
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

func (pc *ProgressClient) UpdateEnrollmentStatus(ctx context.Context, userID, courseID, status string) (*pb.Enrollment, error) {
	req := &pb.UpdateEnrollmentStatusRequest{
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