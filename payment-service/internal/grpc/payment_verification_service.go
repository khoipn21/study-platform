package grpc

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/study-platform/payment-service/internal/model"
	"github.com/study-platform/payment-service/internal/service"
	"github.com/study-platform/payment-service/internal/repository"
	pb "github.com/study-platform/payment-service/proto"
	"github.com/study-platform/pkg/logger"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// PaymentVerificationServer implements the PaymentVerificationService
type PaymentVerificationServer struct {
	pb.UnimplementedPaymentVerificationServiceServer
	accessValidator  *service.AccessValidator
	paymentService   *service.PaymentService
	transactionRepo  *repository.TransactionRepository
	checkoutService  *service.CheckoutService
	previewRepo      *repository.PreviewSessionRepository
	statisticsService *service.StatisticsService
	logger           logger.Logger
}

// NewPaymentVerificationServer creates a new payment verification gRPC server
func NewPaymentVerificationServer(
	accessValidator *service.AccessValidator,
	paymentService *service.PaymentService,
	transactionRepo *repository.TransactionRepository,
	checkoutService *service.CheckoutService,
	previewRepo *repository.PreviewSessionRepository,
	statisticsService *service.StatisticsService,
	logger logger.Logger,
) *PaymentVerificationServer {
	return &PaymentVerificationServer{
		accessValidator:   accessValidator,
		paymentService:    paymentService,
		transactionRepo:   transactionRepo,
		checkoutService:   checkoutService,
		previewRepo:       previewRepo,
		statisticsService: statisticsService,
		logger:            logger,
	}
}

// VerifyCourseAccess validates user access to a course
func (s *PaymentVerificationServer) VerifyCourseAccess(ctx context.Context, req *pb.VerifyCourseAccessRequest) (*pb.VerifyCourseAccessResponse, error) {
	s.logger.Infof("gRPC VerifyCourseAccess called for user %s, course %s", req.UserId, req.CourseId)

	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}
	if req.CourseId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "course_id is required")
	}

	// Convert audit info
	auditInfo := convertAuditInfo(req.AuditInfo)

	// Validate course access
	result, err := s.accessValidator.ValidateCourseAccess(ctx, req.UserId, req.CourseId, auditInfo)
	if err != nil {
		s.logger.Errorf("Failed to validate course access: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to validate course access: %v", err)
	}

	// Convert response
	response := &pb.VerifyCourseAccessResponse{
		UserId:          result.UserID,
		CourseId:        result.CourseID,
		HasAccess:       result.HasAccess,
		AccessLevel:     result.AccessLevel,
		PaymentRequired: result.PaymentRequired,
		PaymentVerified: result.PaymentVerified,
		CoursePrice:     result.CoursePrice,
		Currency:        result.Currency,
		CheckoutUrl:     result.CheckoutURL,
		Message:         result.Message,
		CachedUntil:     timestamppb.New(result.CachedUntil),
		RiskScore:       result.RiskScore,
	}

	if result.TransactionID != nil {
		response.TransactionId = *result.TransactionID
	}

	if result.AccessExpiresAt != nil {
		response.AccessExpiresAt = timestamppb.New(*result.AccessExpiresAt)
	}

	return response, nil
}

// VerifyLectureAccess validates user access to a lecture
func (s *PaymentVerificationServer) VerifyLectureAccess(ctx context.Context, req *pb.VerifyLectureAccessRequest) (*pb.VerifyLectureAccessResponse, error) {
	s.logger.Infof("gRPC VerifyLectureAccess called for user %s, lecture %s", req.UserId, req.LectureId)

	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}
	if req.LectureId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "lecture_id is required")
	}

	// Convert audit info
	auditInfo := convertAuditInfo(req.AuditInfo)

	// For now, we'll pass a dummy course ID since the proto doesn't include it
	// In a real implementation, we'd lookup the course ID from the lecture ID
	courseID := "unknown" // TODO: Get actual course ID from lecture ID

	// Validate lecture access
	result, err := s.accessValidator.ValidateLectureAccess(ctx, req.UserId, courseID, req.LectureId, auditInfo)
	if err != nil {
		s.logger.Errorf("Failed to validate lecture access: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to validate lecture access: %v", err)
	}

	// Convert course access - result IS the course access validation
	courseAccessResponse := &pb.VerifyCourseAccessResponse{
		UserId:          result.UserID,
		CourseId:        result.CourseID,
		HasAccess:       result.HasAccess,
		AccessLevel:     result.AccessLevel,
		PaymentRequired: result.PaymentRequired,
		PaymentVerified: result.PaymentVerified,
		CoursePrice:     result.CoursePrice,
		Currency:        result.Currency,
		CheckoutUrl:     result.CheckoutURL,
		Message:         result.Message,
		CachedUntil:     timestamppb.New(result.CachedUntil),
	}

	if result.TransactionID != nil {
		courseAccessResponse.TransactionId = *result.TransactionID
	}

	response := &pb.VerifyLectureAccessResponse{
		UserId:               req.UserId,
		LectureId:            req.LectureId,
		CourseId:             result.CourseID,
		HasAccess:            result.HasAccess,
		AccessLevel:          result.AccessLevel,
		IsPreview:            result.IsPreview,
		PreviewTimeRemaining: int32(result.PreviewTimeRemaining),
		CourseAccess:         courseAccessResponse,
		Message:              result.Message,
	}

	return response, nil
}

// GetUserTransactions retrieves user's transaction history
func (s *PaymentVerificationServer) GetUserTransactions(ctx context.Context, req *pb.GetUserTransactionsRequest) (*pb.GetUserTransactionsResponse, error) {
	s.logger.Infof("gRPC GetUserTransactions called for user %s", req.UserId)

	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	limit := int(req.Limit)
	if limit <= 0 || limit > 100 {
		limit = 20 // Default limit
	}

	offset := int(req.Offset)
	if offset < 0 {
		offset = 0
	}

	// Get transactions from repository
	transactions, totalCount, err := s.transactionRepo.GetUserTransactionsWithFilters(ctx, req.UserId, limit, offset, req.StatusFilter, req.CourseIdFilter)
	if err != nil {
		s.logger.Errorf("Failed to get user transactions: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve transactions: %v", err)
	}

	// Convert transactions to protobuf
	pbTransactions := make([]*pb.Transaction, len(transactions))
	for i, tx := range transactions {
		pbTransactions[i] = convertTransactionToProto(tx)
	}

	response := &pb.GetUserTransactionsResponse{
		Transactions: pbTransactions,
		TotalCount:   int32(totalCount),
		HasMore:      int64(offset+limit) < totalCount,
	}

	return response, nil
}

// GetTransaction retrieves a specific transaction
func (s *PaymentVerificationServer) GetTransaction(ctx context.Context, req *pb.GetTransactionRequest) (*pb.GetTransactionResponse, error) {
	s.logger.Infof("gRPC GetTransaction called for transaction %s", req.TransactionId)

	if req.TransactionId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "transaction_id is required")
	}
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	// Get transaction
	transaction, err := s.transactionRepo.GetByID(ctx, req.TransactionId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "transaction not found")
	}

	// Verify ownership
	if transaction.UserID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "access denied")
	}

	// Get verification details if available
	var verificationDetails *pb.PaymentVerificationDetails
	if transaction.LemonSqueezyOrderID != nil {
		verificationDetails = &pb.PaymentVerificationDetails{
			ProviderOrderId:    *transaction.LemonSqueezyOrderID,
			ProviderStatus:     transaction.Status,
			ProviderAmount:     transaction.Amount,
			ProviderCurrency:   transaction.Currency,
			VerificationSource: "database",
		}

		if transaction.PaymentVerifiedAt != nil {
			verificationDetails.VerifiedAt = timestamppb.New(*transaction.PaymentVerifiedAt)
		}
	}

	response := &pb.GetTransactionResponse{
		Transaction:         convertTransactionToProto(transaction),
		VerificationDetails: verificationDetails,
	}

	return response, nil
}

// VerifyPaymentStatus provides comprehensive payment verification
func (s *PaymentVerificationServer) VerifyPaymentStatus(ctx context.Context, req *pb.VerifyPaymentStatusRequest) (*pb.VerifyPaymentStatusResponse, error) {
	s.logger.Infof("gRPC VerifyPaymentStatus called for transaction %s", req.TransactionId)

	if req.TransactionId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "transaction_id is required")
	}
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	// Get transaction
	transaction, err := s.transactionRepo.GetByID(ctx, req.TransactionId)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "transaction not found")
	}

	// Verify ownership
	if transaction.UserID != req.UserId {
		return nil, status.Errorf(codes.PermissionDenied, "access denied")
	}

	// Get verification details
	verificationDetails, err := s.paymentService.VerifyPaymentWithProvider(ctx, req.TransactionId, req.ForceProviderCheck)
	if err != nil {
		s.logger.Warnf("Failed to verify payment with provider: %v", err)
	}

	// Get access status if transaction is for a course
	var accessStatus *pb.VerifyCourseAccessResponse
	if transaction.CourseID != nil {
		accessResult, err := s.accessValidator.ValidateCourseAccess(ctx, req.UserId, *transaction.CourseID, &service.AccessAuditInfo{
			UserID:       req.UserId,
			ResourceID:   *transaction.CourseID,
			ResourceType: "course",
		})
		if err == nil {
			accessStatus = &pb.VerifyCourseAccessResponse{
				UserId:          accessResult.UserID,
				CourseId:        accessResult.CourseID,
				HasAccess:       accessResult.HasAccess,
				AccessLevel:     accessResult.AccessLevel,
				PaymentRequired: accessResult.PaymentRequired,
				PaymentVerified: accessResult.PaymentVerified,
				Message:         accessResult.Message,
			}
		}
	}

	response := &pb.VerifyPaymentStatusResponse{
		Transaction:         convertTransactionToProto(transaction),
		VerificationDetails: convertVerificationDetailsToProto(verificationDetails),
		AccessStatus:        accessStatus,
	}

	return response, nil
}

// CreateCheckoutSession creates a new checkout session
func (s *PaymentVerificationServer) CreateCheckoutSession(ctx context.Context, req *pb.CreateCheckoutSessionRequest) (*pb.CreateCheckoutSessionResponse, error) {
	s.logger.Infof("gRPC CreateCheckoutSession called for user %s, course %s", req.UserId, req.CourseId)

	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}
	if req.CourseId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "course_id is required")
	}

	// Convert checkout options
	options := &service.CheckoutOptions{
		SuccessURL:    req.Options.SuccessUrl,
		CancelURL:     req.Options.CancelUrl,
		CustomerEmail: req.Options.CustomerEmail,
		CustomerName:  req.Options.CustomerName,
	}

	if req.Options.CustomData != "" {
		var customData map[string]interface{}
		if err := json.Unmarshal([]byte(req.Options.CustomData), &customData); err == nil {
			options.CustomData = customData
		}
	}

	// Create checkout session
	checkout, err := s.checkoutService.CreateCourseCheckout(ctx, req.UserId, req.CourseId, options)
	if err != nil {
		s.logger.Errorf("Failed to create checkout session: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create checkout session: %v", err)
	}

	response := &pb.CreateCheckoutSessionResponse{
		CheckoutUrl: checkout.CheckoutURL,
		CheckoutId:  checkout.CheckoutID,
		ExpiresAt:   timestamppb.New(time.Now().Add(24 * time.Hour)),
	}

	if checkout.TransactionID != "" {
		response.TransactionId = checkout.TransactionID
	}

	return response, nil
}

// UpdatePreviewSession updates preview session duration
func (s *PaymentVerificationServer) UpdatePreviewSession(ctx context.Context, req *pb.UpdatePreviewSessionRequest) (*pb.UpdatePreviewSessionResponse, error) {
	s.logger.Infof("gRPC UpdatePreviewSession called for user %s, lecture %s", req.UserId, req.LectureId)

	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}
	if req.LectureId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "lecture_id is required")
	}

	err := s.accessValidator.UpdatePreviewSession(ctx, req.UserId, req.LectureId, int(req.DurationSeconds))
	if err != nil {
		s.logger.Errorf("Failed to update preview session: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update preview session: %v", err)
	}

	// Get updated session info
	session, err := s.previewRepo.GetByUserAndLecture(ctx, req.UserId, req.LectureId)
	if err != nil {
		return &pb.UpdatePreviewSessionResponse{
			Success: true,
			Message: "Session updated successfully",
		}, nil
	}

	remainingTime := session.PreviewLimitSeconds - session.SessionDurationSeconds
	if remainingTime < 0 {
		remainingTime = 0
	}

	response := &pb.UpdatePreviewSessionResponse{
		Success:          true,
		RemainingTime:    int32(remainingTime),
		PreviewExhausted: session.PreviewExhausted,
		Message:          "Session updated successfully",
	}

	return response, nil
}

// GetPreviewSession retrieves preview session information
func (s *PaymentVerificationServer) GetPreviewSession(ctx context.Context, req *pb.GetPreviewSessionRequest) (*pb.GetPreviewSessionResponse, error) {
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}
	if req.LectureId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "lecture_id is required")
	}

	session, err := s.previewRepo.GetByUserAndLecture(ctx, req.UserId, req.LectureId)
	if err != nil {
		return &pb.GetPreviewSessionResponse{
			Exists: false,
		}, nil
	}

	response := &pb.GetPreviewSessionResponse{
		Session: convertPreviewSessionToProto(session),
		Exists:  true,
	}

	return response, nil
}

// ClearUserAccessCache clears cached access for a user
func (s *PaymentVerificationServer) ClearUserAccessCache(ctx context.Context, req *pb.ClearUserAccessCacheRequest) (*pb.ClearUserAccessCacheResponse, error) {
	if req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "user_id is required")
	}

	err := s.accessValidator.ClearUserCache(ctx, req.UserId, req.CourseId)
	if err != nil {
		s.logger.Errorf("Failed to clear user cache: %v", err)
		return &pb.ClearUserAccessCacheResponse{
			Success: false,
			Message: fmt.Sprintf("Failed to clear cache: %v", err),
		}, nil
	}

	return &pb.ClearUserAccessCacheResponse{
		Success: true,
		Message: "Cache cleared successfully",
	}, nil
}

// GetAccessStatistics retrieves access analytics
func (s *PaymentVerificationServer) GetAccessStatistics(ctx context.Context, req *pb.GetAccessStatisticsRequest) (*pb.GetAccessStatisticsResponse, error) {
	if req.CourseId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "course_id is required")
	}

	days := int(req.Days)
	if days <= 0 || days > 365 {
		days = 30 // Default to 30 days
	}

	statistics, err := s.statisticsService.GetAccessStatistics(ctx, req.CourseId)
	if err != nil {
		s.logger.Errorf("Failed to get access statistics: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to retrieve statistics: %v", err)
	}

	response := &pb.GetAccessStatisticsResponse{
		Statistics: convertAccessStatisticsToProto(statistics),
	}

	return response, nil
}

// Helper functions

func convertAuditInfo(pbAuditInfo *pb.AccessAuditInfo) *service.AccessAuditInfo {
	if pbAuditInfo == nil {
		return nil
	}

	return &service.AccessAuditInfo{
		UserID:       pbAuditInfo.UserId,
		ResourceID:   pbAuditInfo.ResourceId,
		ResourceType: pbAuditInfo.ResourceType,
		ClientIP:     pbAuditInfo.ClientIp,
		UserAgent:    pbAuditInfo.UserAgent,
	}
}

func convertTransactionToProto(tx *model.Transaction) *pb.Transaction {
	pbTx := &pb.Transaction{
		Id:              tx.ID,
		UserId:          tx.UserID,
		Amount:          tx.Amount,
		Currency:        tx.Currency,
		Status:          tx.Status,
		PaymentProvider: tx.PaymentProvider,
		CreatedAt:       timestamppb.New(tx.CreatedAt),
		UpdatedAt:       timestamppb.New(tx.UpdatedAt),
	}

	if tx.CourseID != nil {
		pbTx.CourseId = *tx.CourseID
	}

	if tx.LemonSqueezyOrderID != nil {
		pbTx.LemonSqueezyOrderId = *tx.LemonSqueezyOrderID
	}

	if tx.LemonSqueezyCheckoutID != nil {
		pbTx.LemonSqueezyCheckoutId = *tx.LemonSqueezyCheckoutID
	}

	if tx.WebhookEventID != nil {
		pbTx.WebhookEventId = *tx.WebhookEventID
	}

	if tx.PaymentVerifiedAt != nil {
		pbTx.PaymentVerifiedAt = timestamppb.New(*tx.PaymentVerifiedAt)
	}

	if tx.RefundedAt != nil {
		pbTx.RefundedAt = timestamppb.New(*tx.RefundedAt)
	}

	if tx.ExpiresAt != nil {
		pbTx.ExpiresAt = timestamppb.New(*tx.ExpiresAt)
	}

	if tx.CustomData != nil {
		if data, err := json.Marshal(tx.CustomData); err == nil {
			pbTx.CustomData = string(data)
		}
	}

	return pbTx
}

func convertVerificationDetailsToProto(details *service.PaymentVerificationDetails) *pb.PaymentVerificationDetails {
	if details == nil {
		return nil
	}

	pbDetails := &pb.PaymentVerificationDetails{
		ProviderOrderId:    details.ProviderReference,
		ProviderStatus:     details.ProviderStatus,
		ProviderAmount:     0, // Extract from details map
		ProviderCurrency:   "", // Extract from details map
		VerificationSource: "provider_api",
		VerifiedAt:         timestamppb.New(details.VerifiedAt),
	}

	// Extract amount and currency from details map
	if details.Details != nil {
		if amount, ok := details.Details["amount"].(float64); ok {
			pbDetails.ProviderAmount = amount
		}
		if currency, ok := details.Details["currency"].(string); ok {
			pbDetails.ProviderCurrency = currency
		}

		if data, err := json.Marshal(details.Details); err == nil {
			pbDetails.ProviderDetails = string(data)
		}
	}

	return pbDetails
}

func convertPreviewSessionToProto(session *repository.PreviewSession) *pb.PreviewSession {
	pbSession := &pb.PreviewSession{
		Id:                    session.ID,
		UserId:                session.UserID,
		LectureId:             session.LectureID,
		SessionStartedAt:      timestamppb.New(session.SessionStartedAt),
		SessionDurationSeconds: int32(session.SessionDurationSeconds),
		PreviewLimitSeconds:   int32(session.PreviewLimitSeconds),
		PreviewExhausted:      session.PreviewExhausted,
		IpAddress:             session.IPAddress,
		CreatedAt:             timestamppb.New(session.CreatedAt),
		UpdatedAt:             timestamppb.New(session.UpdatedAt),
	}

	if session.LastAccessedAt != nil {
		pbSession.LastAccessedAt = timestamppb.New(*session.LastAccessedAt)
	}

	return pbSession
}

func convertAccessStatisticsToProto(stats *service.AccessStatistics) *pb.AccessStatistics {
	pbStats := &pb.AccessStatistics{
		TotalAttempts:         stats.TotalAttempts,
		SuccessfulAccess:      stats.SuccessfulAccess,
		PreviewAccess:         stats.PreviewAccess,
		DeniedAccess:          stats.DeniedAccess,
		UniqueUsers:           stats.UniqueUsers,
		SuccessRate:           stats.SuccessRate,
		ConversionRate:        stats.ConversionRate,
		TotalRevenue:          stats.TotalRevenue,
		AverageRevenuePerUser: stats.AverageRevenuePerUser,
	}

	// Convert daily stats
	if stats.DailyStats != nil {
		pbStats.DailyStats = make([]*pb.DailyStatistic, len(stats.DailyStats))
		for i, daily := range stats.DailyStats {
			pbStats.DailyStats[i] = &pb.DailyStatistic{
				Date:        daily.Date,
				Attempts:    daily.Attempts,
				Successes:   daily.Successes,
				Conversions: daily.Conversions,
				Revenue:     daily.Revenue,
				UniqueUsers: daily.UniqueUsers,
			}
		}
	}

	return pbStats
}