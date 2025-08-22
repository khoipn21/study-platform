package service

import (
	"context"
	"fmt"
	"math"
	"time"

	"video-service/internal/config"
	"video-service/internal/model"
	"video-service/internal/queue"
)

type NetworkIntelligenceService struct {
	config      *config.Config
	redisClient *queue.RedisClient
}

func NewNetworkIntelligenceService(cfg *config.Config, redisClient *queue.RedisClient) *NetworkIntelligenceService {
	return &NetworkIntelligenceService{
		config:      cfg,
		redisClient: redisClient,
	}
}

// CalculateQualityScore calculates a quality score from 1-10 based on network metrics
func (nis *NetworkIntelligenceService) CalculateQualityScore(metrics *model.NetworkMetrics) int {
	score := 10

	// Bandwidth factor (40% weight)
	if metrics.BandwidthMbps < 1.0 {
		score -= 4
	} else if metrics.BandwidthMbps < 3.0 {
		score -= 2
	} else if metrics.BandwidthMbps < 5.0 {
		score -= 1
	}

	// Latency factor (30% weight)
	if metrics.LatencyMs > 500 {
		score -= 3
	} else if metrics.LatencyMs > 200 {
		score -= 2
	} else if metrics.LatencyMs > 100 {
		score -= 1
	}

	// Packet loss factor (20% weight)
	if metrics.PacketLossPercent > 0.05 {
		score -= 2
	} else if metrics.PacketLossPercent > 0.01 {
		score -= 1
	}

	// Buffer health factor (10% weight)
	if metrics.BufferHealthSeconds < 3 {
		score -= 1
	}

	if score < 1 {
		score = 1
	}

	return score
}

// RecommendQuality recommends a video quality based on the quality score
func (nis *NetworkIntelligenceService) RecommendQuality(score int, currentQuality string) string {
	qualityMap := map[int]string{
		1: "240p", 2: "240p", 3: "360p",
		4: "360p", 5: "480p", 6: "480p",
		7: "720p", 8: "720p", 9: "1080p", 10: "1080p",
	}

	recommended := qualityMap[score]

	// Avoid frequent switching
	if nis.shouldPreventSwitch(currentQuality, recommended) {
		return currentQuality
	}

	return recommended
}

// shouldPreventSwitch prevents frequent quality switching
func (nis *NetworkIntelligenceService) shouldPreventSwitch(current, recommended string) bool {
	if current == "" {
		return false
	}

	qualityOrder := map[string]int{
		"240p": 1, "360p": 2, "480p": 3, "720p": 4, "1080p": 5,
	}

	currentLevel := qualityOrder[current]
	recommendedLevel := qualityOrder[recommended]

	// Only switch if the difference is significant (more than 1 level)
	return math.Abs(float64(currentLevel-recommendedLevel)) <= 1
}

// ProcessNetworkUpdate processes network status updates and provides recommendations
func (nis *NetworkIntelligenceService) ProcessNetworkUpdate(ctx context.Context, sessionID string, update *model.NetworkStatusUpdate) (*model.NetworkStatusResponse, error) {
	// Create network metrics from update
	metrics := &model.NetworkMetrics{
		SessionID:           sessionID,
		BandwidthMbps:       update.BandwidthMbps,
		LatencyMs:           update.LatencyMs,
		PacketLossPercent:   update.PacketLoss,
		ConnectionType:      update.ConnectionType,
		BufferHealthSeconds: update.BufferHealth,
		Timestamp:           time.Now(),
	}

	// Calculate quality score
	score := nis.CalculateQualityScore(metrics)
	metrics.QualityScore = score

	// Get current quality from update or cache
	currentQuality := update.CurrentQuality
	if currentQuality == "" {
		if cached, err := nis.redisClient.GetCachedQualityRecommendation(ctx, sessionID); err == nil {
			currentQuality = cached
		}
	}

	// Recommend new quality
	recommendedQuality := nis.RecommendQuality(score, currentQuality)
	metrics.RecommendedQuality = recommendedQuality

	// Cache the network metrics
	if err := nis.redisClient.CacheNetworkMetrics(ctx, sessionID, metrics, 5*time.Minute); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to cache network metrics: %v\n", err)
	}

	// Cache the quality recommendation
	if err := nis.redisClient.CacheQualityRecommendation(ctx, sessionID, recommendedQuality, 1*time.Minute); err != nil {
		// Log error but don't fail the request
		fmt.Printf("Failed to cache quality recommendation: %v\n", err)
	}

	// Determine if preloading should be enabled
	shouldPreload := nis.shouldEnablePreloading(metrics)

	// Calculate buffer target
	bufferTarget := nis.calculateBufferTarget(metrics)

	// Publish network status update to Redis
	wsNetworkStatus := &model.WSNetworkStatus{
		SessionID:          sessionID,
		Timestamp:          time.Now().Format(time.RFC3339),
		BandwidthMbps:      metrics.BandwidthMbps,
		LatencyMs:          metrics.LatencyMs,
		PacketLoss:         metrics.PacketLossPercent,
		ConnectionType:     metrics.ConnectionType,
		QualityScore:       score,
		CurrentQuality:     currentQuality,
		RecommendedQuality: recommendedQuality,
		BufferHealth:       metrics.BufferHealthSeconds,
	}

	if err := nis.redisClient.PublishNetworkStatus(ctx, sessionID, wsNetworkStatus); err != nil {
		fmt.Printf("Failed to publish network status: %v\n", err)
	}

	return &model.NetworkStatusResponse{
		RecommendedQuality: recommendedQuality,
		QualityScore:       score,
		ShouldPreload:      shouldPreload,
		BufferTarget:       bufferTarget,
	}, nil
}

// shouldEnablePreloading determines if preloading should be enabled based on network conditions
func (nis *NetworkIntelligenceService) shouldEnablePreloading(metrics *model.NetworkMetrics) bool {
	// Enable preloading for good network conditions
	return metrics.BandwidthMbps > 3.0 && 
		   metrics.LatencyMs < 200 && 
		   metrics.PacketLossPercent < 0.02 && 
		   metrics.BufferHealthSeconds > 5
}

// calculateBufferTarget calculates optimal buffer target based on network conditions
func (nis *NetworkIntelligenceService) calculateBufferTarget(metrics *model.NetworkMetrics) int {
	baseTarget := nis.config.BufferTargetSeconds

	// Adjust based on network quality
	if metrics.BandwidthMbps < 2.0 || metrics.LatencyMs > 300 {
		return baseTarget + 5 // Increase buffer for poor connections
	}

	if metrics.BandwidthMbps > 10.0 && metrics.LatencyMs < 50 {
		return baseTarget - 2 // Decrease buffer for excellent connections
	}

	return baseTarget
}

// DetectQualityChangeNeeded checks if quality change is needed based on network trends
func (nis *NetworkIntelligenceService) DetectQualityChangeNeeded(ctx context.Context, sessionID string, currentMetrics *model.NetworkMetrics) (*model.WSQualityChange, error) {
	// Get cached previous metrics
	previousMetrics, err := nis.redisClient.GetCachedNetworkMetrics(ctx, sessionID)
	if err != nil {
		// No previous metrics, no change needed
		return nil, nil
	}

	currentScore := nis.CalculateQualityScore(currentMetrics)
	previousScore := nis.CalculateQualityScore(previousMetrics)

	// Calculate score difference
	scoreDiff := float64(currentScore-previousScore) / float64(previousScore)

	// Only recommend change if the score difference exceeds threshold
	if math.Abs(scoreDiff) < nis.config.QualityChangeThreshold {
		return nil, nil
	}

	currentQuality := currentMetrics.RecommendedQuality
	if currentQuality == "" {
		currentQuality = "720p" // Default
	}

	newQuality := nis.RecommendQuality(currentScore, currentQuality)
	if newQuality == currentQuality {
		return nil, nil
	}

	// Determine reason for quality change
	reason := "network_optimization"
	if currentScore > previousScore {
		reason = "network_improvement"
	} else {
		reason = "network_degradation"
	}

	return &model.WSQualityChange{
		SessionID:   sessionID,
		FromQuality: currentQuality,
		ToQuality:   newQuality,
		Reason:      reason,
		Timestamp:   time.Now().Format(time.RFC3339),
	}, nil
}

// GetQualityRecommendation provides quality recommendation for a new session
func (nis *NetworkIntelligenceService) GetQualityRecommendation(ctx context.Context, connectionType string) string {
	// Default recommendations based on connection type
	switch connectionType {
	case "wifi":
		return "720p"
	case "5g":
		return "1080p"
	case "4g":
		return "480p"
	case "3g":
		return "360p"
	case "ethernet":
		return "1080p"
	default:
		return "480p" // Safe default
	}
}

// GeneratePreloadSegments generates list of segments to preload based on current position and network conditions
func (nis *NetworkIntelligenceService) GeneratePreloadSegments(currentTimeSeconds int, bufferHealth int, quality string) []string {
	// Calculate how many segments to preload based on buffer health
	segmentDuration := 10 // seconds per segment
	preloadCount := 3     // default

	if bufferHealth < 5 {
		preloadCount = 5 // More aggressive preloading for low buffer
	} else if bufferHealth > 15 {
		preloadCount = 2 // Less preloading for healthy buffer
	}

	segments := make([]string, preloadCount)
	startSegment := (currentTimeSeconds / segmentDuration) + 1

	for i := 0; i < preloadCount; i++ {
		segmentNum := startSegment + i
		segments[i] = fmt.Sprintf("seg_%d_%s.ts", segmentNum, quality)
	}

	return segments
}

// AnalyzeNetworkPattern analyzes network patterns and provides insights
func (nis *NetworkIntelligenceService) AnalyzeNetworkPattern(ctx context.Context, sessionID string, windowMinutes int) (*NetworkPattern, error) {
	// This would typically query historical data from database
	// For now, return a basic pattern analysis
	
	cachedMetrics, err := nis.redisClient.GetCachedNetworkMetrics(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("no network data available for analysis")
	}

	pattern := &NetworkPattern{
		SessionID:           sessionID,
		AverageBandwidth:    cachedMetrics.BandwidthMbps,
		AverageLatency:      float64(cachedMetrics.LatencyMs),
		AveragePacketLoss:   cachedMetrics.PacketLossPercent,
		StabilityScore:      nis.calculateStabilityScore(cachedMetrics),
		RecommendedQuality:  cachedMetrics.RecommendedQuality,
		OptimalBufferSize:   nis.calculateBufferTarget(cachedMetrics),
		ConnectionStable:    cachedMetrics.QualityScore > 7,
		LastAnalyzedAt:      time.Now(),
	}

	return pattern, nil
}

// NetworkPattern represents analyzed network pattern data
type NetworkPattern struct {
	SessionID           string    `json:"session_id"`
	AverageBandwidth    float64   `json:"average_bandwidth"`
	AverageLatency      float64   `json:"average_latency"`
	AveragePacketLoss   float64   `json:"average_packet_loss"`
	StabilityScore      int       `json:"stability_score"`
	RecommendedQuality  string    `json:"recommended_quality"`
	OptimalBufferSize   int       `json:"optimal_buffer_size"`
	ConnectionStable    bool      `json:"connection_stable"`
	LastAnalyzedAt      time.Time `json:"last_analyzed_at"`
}

// calculateStabilityScore calculates network stability from 1-10
func (nis *NetworkIntelligenceService) calculateStabilityScore(metrics *model.NetworkMetrics) int {
	score := 10

	// Reduce score for high latency variations (simulated)
	if metrics.LatencyMs > 100 {
		score -= 2
	}

	// Reduce score for packet loss
	if metrics.PacketLossPercent > 0.01 {
		score -= 3
	}

	// Reduce score for low bandwidth
	if metrics.BandwidthMbps < 3.0 {
		score -= 2
	}

	if score < 1 {
		score = 1
	}

	return score
}