package service

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"video-service/internal/config"
	"video-service/internal/model"
	"video-service/internal/queue"
	"video-service/internal/repository"
)

type VideoService struct {
	config          *config.Config
	videoRepo       *repository.VideoRepository
	sessionRepo     *repository.SessionRepository
	cloudflareService *CloudflareService
	redisClient     *queue.RedisClient
	networkService  *NetworkIntelligenceService
}

func NewVideoService(cfg *config.Config, videoRepo *repository.VideoRepository, sessionRepo *repository.SessionRepository, cloudflareService *CloudflareService, redisClient *queue.RedisClient, networkService *NetworkIntelligenceService) *VideoService {
	return &VideoService{
		config:          cfg,
		videoRepo:       videoRepo,
		sessionRepo:     sessionRepo,
		cloudflareService: cloudflareService,
		redisClient:     redisClient,
		networkService:  networkService,
	}
}

// UploadVideo handles video upload process
func (vs *VideoService) UploadVideo(ctx context.Context, userID uuid.UUID, req *model.UploadVideoRequest) (*model.UploadVideoResponse, error) {
	// Create direct upload URL from Cloudflare
	maxDurationSeconds := 3600 // 1 hour default
	uploadResponse, err := vs.cloudflareService.CreateDirectUploadURL(maxDurationSeconds)
	if err != nil {
		return nil, fmt.Errorf("failed to create upload URL: %w", err)
	}

	// Create video record in database
	video := &model.Video{
		CloudflareUID:    uploadResponse.Result.UID,
		Title:           req.Title,
		Description:     req.Description,
		UploadUserID:    userID,
		CourseID:        req.CourseID,
		LectureID:       req.LectureID,
		Status:          "uploading",
		Visibility:      getVisibility(req.Visibility),
		Metadata:        model.JSONB{},
	}

	if err := vs.videoRepo.CreateVideo(video); err != nil {
		return nil, fmt.Errorf("failed to create video record: %w", err)
	}

	return &model.UploadVideoResponse{
		VideoID:                 video.ID,
		CloudflareUID:           video.CloudflareUID,
		Title:                  video.Title,
		Status:                 video.Status,
		UploadURL:              uploadResponse.Result.UploadURL,
		EstimatedProcessingTime: "5-10 minutes",
	}, nil
}

// GetVideo retrieves video details
func (vs *VideoService) GetVideo(ctx context.Context, videoID uuid.UUID, userID *uuid.UUID) (*model.Video, error) {
	video, err := vs.videoRepo.GetVideoByID(videoID)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if !vs.canUserAccessVideo(video, userID) {
		return nil, fmt.Errorf("access denied")
	}

	// Update URLs from Cloudflare if needed
	if video.StreamURL == "" || video.ThumbnailURL == "" {
		video.StreamURL = vs.cloudflareService.GetStreamURL(video.CloudflareUID)
		video.ThumbnailURL = vs.cloudflareService.GetThumbnailURL(video.CloudflareUID)

		// Update in database
		vs.videoRepo.UpdateVideoURLs(video.ID, video.ThumbnailURL, video.StreamURL, video.PreviewURL)
	}

	return video, nil
}

// CreateViewingSession creates a new viewing session
func (vs *VideoService) CreateViewingSession(ctx context.Context, videoID uuid.UUID, userID uuid.UUID, req *model.CreateSessionRequest) (*model.CreateSessionResponse, error) {
	// Get video details
	video, err := vs.videoRepo.GetVideoByID(videoID)
	if err != nil {
		return nil, err
	}

	// Check permissions
	if !vs.canUserAccessVideo(video, &userID) {
		return nil, fmt.Errorf("access denied")
	}

	// Check if video is ready for streaming
	if video.Status != "ready" {
		return nil, fmt.Errorf("video is not ready for streaming")
	}

	// Generate session ID
	sessionID := uuid.New().String()

	// Create viewing session
	session := &model.ViewingSession{
		SessionID:    sessionID,
		UserID:       userID,
		VideoID:      videoID,
		StartedAt:    time.Now(),
		LastHeartbeat: time.Now(),
		UserAgent:    req.DeviceInfo.UserAgent,
		IPAddress:    "", // Will be set by handler
	}

	if err := vs.sessionRepo.CreateViewingSession(session); err != nil {
		return nil, fmt.Errorf("failed to create viewing session: %w", err)
	}

	// Get video qualities
	qualities, err := vs.videoRepo.GetVideoQualitiesByVideoID(videoID)
	if err != nil {
		// If no qualities found, create default ones based on Cloudflare Stream
		qualities = vs.generateDefaultQualities(video)
	}

	// Get recommended quality based on connection type
	recommendedQuality := vs.networkService.GetQualityRecommendation(ctx, req.DeviceInfo.ConnectionType)

	// Generate URLs
	streamURL := vs.cloudflareService.GetStreamURL(video.CloudflareUID)
	thumbnailURL := vs.cloudflareService.GetThumbnailURL(video.CloudflareUID)
	websocketURL := fmt.Sprintf("ws://localhost:%s/api/videos/ws/%s", vs.config.Port, sessionID)

	response := &model.CreateSessionResponse{
		SessionID:          sessionID,
		StreamURL:          streamURL,
		ThumbnailURL:       thumbnailURL,
		RecommendedQuality: recommendedQuality,
		WebSocketURL:       websocketURL,
		ExpiresAt:          time.Now().Add(6 * time.Hour),
	}

	// Convert pointer slice to value slice
	qualityValues := make([]model.VideoQuality, len(qualities))
	for i, q := range qualities {
		qualityValues[i] = *q
	}
	response.Qualities = qualityValues

	return response, nil
}

// UpdateNetworkStatus processes network status updates
func (vs *VideoService) UpdateNetworkStatus(ctx context.Context, sessionID string, update *model.NetworkStatusUpdate) (*model.NetworkStatusResponse, error) {
	// Verify session exists
	session, err := vs.sessionRepo.GetViewingSessionBySessionID(sessionID)
	if err != nil {
		return nil, fmt.Errorf("session not found: %w", err)
	}

	// Update session heartbeat
	vs.sessionRepo.UpdateSessionHeartbeat(sessionID)

	// Process through network intelligence service
	response, err := vs.networkService.ProcessNetworkUpdate(ctx, sessionID, update)
	if err != nil {
		return nil, fmt.Errorf("failed to process network update: %w", err)
	}

	// Store network metrics in database
	metrics := &model.NetworkMetrics{
		SessionID:           sessionID,
		UserID:              session.UserID,
		BandwidthMbps:       update.BandwidthMbps,
		LatencyMs:           update.LatencyMs,
		PacketLossPercent:   update.PacketLoss,
		ConnectionType:      update.ConnectionType,
		QualityScore:        response.QualityScore,
		RecommendedQuality:  response.RecommendedQuality,
		BufferHealthSeconds: update.BufferHealth,
	}

	if err := vs.sessionRepo.CreateNetworkMetrics(metrics); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to store network metrics: %v\n", err)
	}

	return response, nil
}

// UpdateSessionProgress updates viewing progress
func (vs *VideoService) UpdateSessionProgress(ctx context.Context, sessionID string, currentTime int, quality string) error {
	// Get session
	session, err := vs.sessionRepo.GetViewingSessionBySessionID(sessionID)
	if err != nil {
		return err
	}

	// Calculate watch time increment
	// Simple calculation - in practice, you'd want more sophisticated tracking
	watchTimeIncrement := 10 // seconds, assuming regular updates

	newWatchTime := session.TotalWatchTimeSeconds + watchTimeIncrement

	// Update session progress
	if err := vs.sessionRepo.UpdateSessionProgress(sessionID, currentTime, quality, newWatchTime); err != nil {
		return err
	}

	// Check if video is completed (90% watched)
	video, err := vs.videoRepo.GetVideoByID(session.VideoID)
	if err == nil && video.DurationSeconds != nil {
		completionThreshold := float64(*video.DurationSeconds) * 0.9
		if float64(currentTime) >= completionThreshold && !session.Completed {
			vs.sessionRepo.MarkSessionCompleted(sessionID)
		}
	}

	return nil
}

// ListUserVideos lists videos uploaded by a user
func (vs *VideoService) ListUserVideos(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*model.Video, error) {
	return vs.videoRepo.ListVideosByUser(userID, limit, offset)
}

// ListCourseVideos lists videos for a course
func (vs *VideoService) ListCourseVideos(ctx context.Context, courseID uuid.UUID, limit, offset int) ([]*model.Video, error) {
	return vs.videoRepo.ListVideosByCourse(courseID, limit, offset)
}

// SearchVideos searches videos
func (vs *VideoService) SearchVideos(ctx context.Context, query string, limit, offset int) ([]*model.Video, error) {
	return vs.videoRepo.SearchVideos(query, limit, offset)
}

// ProcessCloudflareWebhook processes webhooks from Cloudflare Stream
func (vs *VideoService) ProcessCloudflareWebhook(ctx context.Context, webhookData map[string]interface{}) error {
	// Extract UID from webhook
	uid, ok := webhookData["uid"].(string)
	if !ok {
		return fmt.Errorf("missing uid in webhook data")
	}

	// Get video from database
	video, err := vs.videoRepo.GetVideoByCloudflareUID(uid)
	if err != nil {
		return fmt.Errorf("video not found: %w", err)
	}

	// Get video details from Cloudflare
	streamVideo, err := vs.cloudflareService.GetVideoDetails(uid)
	if err != nil {
		return fmt.Errorf("failed to get video details: %w", err)
	}

	// Update video status and metadata
	video.Status = vs.mapCloudflareStatus(streamVideo)
	video.DurationSeconds = vs.getDurationFromCloudflare(streamVideo)
	video.ThumbnailURL = vs.cloudflareService.GetThumbnailURL(uid)
	video.StreamURL = vs.cloudflareService.GetStreamURL(uid)

	// Update video in database
	if err := vs.videoRepo.UpdateVideo(video); err != nil {
		return fmt.Errorf("failed to update video: %w", err)
	}

	// Create default quality variants if video is ready
	if video.Status == "ready" {
		vs.createDefaultQualityVariants(video)
	}

	return nil
}

// DeleteVideo deletes a video
func (vs *VideoService) DeleteVideo(ctx context.Context, videoID uuid.UUID, userID uuid.UUID) error {
	// Get video
	video, err := vs.videoRepo.GetVideoByID(videoID)
	if err != nil {
		return err
	}

	// Check permissions (only owner can delete)
	if video.UploadUserID != userID {
		return fmt.Errorf("access denied")
	}

	// Delete from Cloudflare
	if err := vs.cloudflareService.DeleteVideo(video.CloudflareUID); err != nil {
		// Log error but continue with database deletion
		fmt.Printf("Failed to delete video from Cloudflare: %v\n", err)
	}

	// Soft delete from database
	return vs.videoRepo.DeleteVideo(videoID)
}

// Helper methods

func (vs *VideoService) canUserAccessVideo(video *model.Video, userID *uuid.UUID) bool {
	// Public videos are accessible to all
	if video.Visibility == "public" {
		return true
	}

	// Private videos require authentication and ownership/enrollment
	if userID == nil {
		return false
	}

	// Owner can always access
	if video.UploadUserID == *userID {
		return true
	}

	// TODO: Add course enrollment check here
	// For now, allow access to unlisted videos if user is authenticated
	if video.Visibility == "unlisted" {
		return true
	}

	return false
}

func getVisibility(visibility string) string {
	switch visibility {
	case "public", "private", "unlisted":
		return visibility
	default:
		return "private"
	}
}

func (vs *VideoService) generateDefaultQualities(video *model.Video) []*model.VideoQuality {
	// Generate default quality variants based on typical streaming qualities
	qualities := []*model.VideoQuality{
		{
			VideoID:      video.ID,
			QualityLabel: "360p",
			BitrateKbps:  800,
			Width:        640,
			Height:       360,
			FPS:          30,
			Codec:        "h264",
			URL:          vs.cloudflareService.GetStreamURL(video.CloudflareUID),
		},
		{
			VideoID:      video.ID,
			QualityLabel: "720p",
			BitrateKbps:  2500,
			Width:        1280,
			Height:       720,
			FPS:          30,
			Codec:        "h264",
			URL:          vs.cloudflareService.GetStreamURL(video.CloudflareUID),
		},
		{
			VideoID:      video.ID,
			QualityLabel: "1080p",
			BitrateKbps:  5000,
			Width:        1920,
			Height:       1080,
			FPS:          30,
			Codec:        "h264",
			URL:          vs.cloudflareService.GetStreamURL(video.CloudflareUID),
		},
	}

	return qualities
}

func (vs *VideoService) createDefaultQualityVariants(video *model.Video) {
	qualities := vs.generateDefaultQualities(video)
	for _, quality := range qualities {
		vs.videoRepo.CreateVideoQuality(quality)
	}
}

func (vs *VideoService) mapCloudflareStatus(streamVideo *StreamVideo) string {
	if streamVideo.ReadyToStream {
		return "ready"
	}

	// Check status from Cloudflare
	if status, ok := streamVideo.Status["state"].(string); ok {
		switch status {
		case "pendingupload":
			return "uploading"
		case "downloading":
			return "processing"
		case "queued":
			return "processing"
		case "inprogress":
			return "processing"
		case "ready":
			return "ready"
		case "error":
			return "error"
		default:
			return "processing"
		}
	}

	return "processing"
}

func (vs *VideoService) getDurationFromCloudflare(streamVideo *StreamVideo) *int {
	if streamVideo.Duration > 0 {
		duration := int(streamVideo.Duration)
		return &duration
	}
	return nil
}