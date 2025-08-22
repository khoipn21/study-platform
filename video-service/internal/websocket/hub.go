package websocket

import (
	"context"
	"encoding/json"
	"log"
	"sync"
	"time"

	"video-service/internal/model"
	"video-service/internal/queue"
	"video-service/internal/service"
)

// Hub maintains the set of active clients and broadcasts messages
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Session-based client mapping
	sessionClients map[string]*Client

	// Redis client for message queue
	redisClient *queue.RedisClient

	// Network intelligence service
	networkService *service.NetworkIntelligenceService

	// Mutex for thread-safe operations
	mutex sync.RWMutex

	// Context for graceful shutdown
	ctx    context.Context
	cancel context.CancelFunc
}

func NewHub(redisClient *queue.RedisClient, networkService *service.NetworkIntelligenceService) *Hub {
	ctx, cancel := context.WithCancel(context.Background())

	return &Hub{
		broadcast:      make(chan []byte),
		register:       make(chan *Client),
		unregister:     make(chan *Client),
		clients:        make(map[*Client]bool),
		sessionClients: make(map[string]*Client),
		redisClient:    redisClient,
		networkService: networkService,
		ctx:            ctx,
		cancel:         cancel,
	}
}

func (h *Hub) Run() {
	defer h.cancel()

	// Start Redis subscribers in goroutines
	go h.startNetworkStatusSubscriber()
	go h.startQualityChangeSubscriber()
	go h.startAnalyticsSubscriber()
	go h.startHeartbeatProcessor()

	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			h.broadcastMessage(message)

		case <-h.ctx.Done():
			log.Println("WebSocket hub shutting down...")
			return
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	h.clients[client] = true
	h.sessionClients[client.SessionID] = client

	log.Printf("Client registered for session: %s", client.SessionID)

	// Add session to Redis active sessions
	if h.redisClient != nil {
		err := h.redisClient.AddActiveSession(h.ctx, client.SessionID, client.UserID, 2*time.Hour)
		if err != nil {
			log.Printf("Failed to add active session to Redis: %v", err)
		}
	}

	// Send initial connection confirmation
	welcome := model.WSMessage{
		Type: "connection_established",
		Data: map[string]interface{}{
			"session_id": client.SessionID,
			"timestamp":  time.Now().Format(time.RFC3339),
		},
	}

	client.SendMessage(welcome)
}

func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		delete(h.sessionClients, client.SessionID)
		close(client.send)

		log.Printf("Client unregistered for session: %s", client.SessionID)

		// Remove session from Redis active sessions
		if h.redisClient != nil {
			err := h.redisClient.RemoveActiveSession(h.ctx, client.SessionID, client.UserID)
			if err != nil {
				log.Printf("Failed to remove active session from Redis: %v", err)
			}
		}
	}
}

func (h *Hub) broadcastMessage(message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		select {
		case client.send <- message:
		default:
			close(client.send)
			delete(h.clients, client)
			delete(h.sessionClients, client.SessionID)
		}
	}
}

// SendToSession sends a message to a specific session
func (h *Hub) SendToSession(sessionID string, message model.WSMessage) {
	h.mutex.RLock()
	client, ok := h.sessionClients[sessionID]
	h.mutex.RUnlock()

	if ok && client != nil {
		client.SendMessage(message)
	}
}

// ProcessClientMessage processes incoming messages from clients
func (h *Hub) ProcessClientMessage(client *Client, messageType int, data []byte) {
	var wsMessage model.WSMessage
	if err := json.Unmarshal(data, &wsMessage); err != nil {
		log.Printf("Failed to unmarshal WebSocket message: %v", err)
		return
	}

	switch wsMessage.Type {
	case "network_status":
		h.processNetworkStatus(client, wsMessage.Data)
	case "progress_update":
		h.processProgressUpdate(client, wsMessage.Data)
	case "quality_change":
		h.processQualityChange(client, wsMessage.Data)
	case "heartbeat":
		h.processHeartbeat(client, wsMessage.Data)
	default:
		log.Printf("Unknown message type: %s", wsMessage.Type)
	}
}

func (h *Hub) processNetworkStatus(client *Client, data interface{}) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		log.Printf("Invalid network status data format")
		return
	}

	// Extract network status data
	networkUpdate := &model.NetworkStatusUpdate{}
	
	if bw, ok := dataMap["bandwidth_mbps"].(float64); ok {
		networkUpdate.BandwidthMbps = bw
	}
	if latency, ok := dataMap["latency_ms"].(float64); ok {
		networkUpdate.LatencyMs = int(latency)
	}
	if packetLoss, ok := dataMap["packet_loss"].(float64); ok {
		networkUpdate.PacketLoss = packetLoss
	}
	if connType, ok := dataMap["connection_type"].(string); ok {
		networkUpdate.ConnectionType = connType
	}
	if bufferHealth, ok := dataMap["buffer_health"].(float64); ok {
		networkUpdate.BufferHealth = int(bufferHealth)
	}
	if currentTime, ok := dataMap["current_time"].(float64); ok {
		networkUpdate.CurrentTime = int(currentTime)
	}
	if quality, ok := dataMap["current_quality"].(string); ok {
		networkUpdate.CurrentQuality = quality
	}

	// Process through network intelligence service
	if h.networkService != nil {
		response, err := h.networkService.ProcessNetworkUpdate(h.ctx, client.SessionID, networkUpdate)
		if err != nil {
			log.Printf("Failed to process network update: %v", err)
			return
		}

		// Send quality recommendation back to client
		recommendation := model.WSMessage{
			Type: "quality_recommendation",
			Data: model.WSQualityRecommendation{
				RecommendedQuality: response.RecommendedQuality,
				Reason:            "network_optimization",
				Confidence:        0.85,
			},
		}

		client.SendMessage(recommendation)

		// Send preload instruction if appropriate
		if response.ShouldPreload {
			segments := h.networkService.GeneratePreloadSegments(
				networkUpdate.CurrentTime,
				networkUpdate.BufferHealth,
				response.RecommendedQuality,
			)

			preload := model.WSMessage{
				Type: "preload",
				Data: model.WSPreloadInstruction{
					Segments: segments,
					Priority: "high",
				},
			}

			client.SendMessage(preload)
		}
	}
}

func (h *Hub) processProgressUpdate(client *Client, data interface{}) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		log.Printf("Invalid progress update data format")
		return
	}

	// Extract progress data
	var currentTime int
	var quality string
	var paused bool

	if ct, ok := dataMap["current_time"].(float64); ok {
		currentTime = int(ct)
	}
	if q, ok := dataMap["quality"].(string); ok {
		quality = q
	}
	if p, ok := dataMap["paused"].(bool); ok {
		paused = p
	}

	// Publish analytics event
	if h.redisClient != nil {
		payload := map[string]interface{}{
			"session_id":     client.SessionID,
			"video_id":       client.VideoID,
			"user_id":        client.UserID,
			"current_time":   currentTime,
			"quality":        quality,
			"paused":         paused,
			"timestamp":      time.Now().Format(time.RFC3339),
		}

		// Store in Redis for batch processing
		err := h.redisClient.PublishHeartbeat(h.ctx, client.SessionID, payload)
		if err != nil {
			log.Printf("Failed to publish progress update: %v", err)
		}
	}
}

func (h *Hub) processQualityChange(client *Client, data interface{}) {
	dataMap, ok := data.(map[string]interface{})
	if !ok {
		log.Printf("Invalid quality change data format")
		return
	}

	var requestedQuality, reason string
	if rq, ok := dataMap["requested_quality"].(string); ok {
		requestedQuality = rq
	}
	if r, ok := dataMap["reason"].(string); ok {
		reason = r
	}

	// Create quality change event
	qualityChange := &model.WSQualityChange{
		SessionID:   client.SessionID,
		VideoID:     client.VideoID,
		ToQuality:   requestedQuality,
		Reason:      reason,
		Timestamp:   time.Now().Format(time.RFC3339),
	}

	// Publish quality change event
	if h.redisClient != nil {
		err := h.redisClient.PublishQualityChange(h.ctx, client.SessionID, qualityChange)
		if err != nil {
			log.Printf("Failed to publish quality change: %v", err)
		}
	}

	// Send analytics event
	analyticsEvent := &model.WSAnalyticsEvent{
		Event: "quality_change",
		To:    requestedQuality,
	}

	if h.redisClient != nil {
		err := h.redisClient.PublishAnalyticsEvent(h.ctx, analyticsEvent)
		if err != nil {
			log.Printf("Failed to publish analytics event: %v", err)
		}
	}
}

func (h *Hub) processHeartbeat(client *Client, data interface{}) {
	// Update client's last activity
	client.UpdateActivity()

	// Publish heartbeat to Redis
	if h.redisClient != nil {
		heartbeatData := map[string]interface{}{
			"user_id":    client.UserID,
			"video_id":   client.VideoID,
			"timestamp":  time.Now().Format(time.RFC3339),
		}

		if data != nil {
			if dataMap, ok := data.(map[string]interface{}); ok {
				for k, v := range dataMap {
					heartbeatData[k] = v
				}
			}
		}

		err := h.redisClient.PublishHeartbeat(h.ctx, client.SessionID, heartbeatData)
		if err != nil {
			log.Printf("Failed to publish heartbeat: %v", err)
		}
	}
}

// Redis subscribers
func (h *Hub) startNetworkStatusSubscriber() {
	// Subscribe to all network status channels using pattern
	// This is a simplified version - in practice, you'd want more targeted subscriptions
	log.Println("Starting network status subscriber...")
}

func (h *Hub) startQualityChangeSubscriber() {
	// Subscribe to quality change events
	log.Println("Starting quality change subscriber...")
}

func (h *Hub) startAnalyticsSubscriber() {
	// Subscribe to analytics events
	err := h.redisClient.SubscribeAnalytics(h.ctx, func(event *model.WSAnalyticsEvent) {
		// Process analytics event
		log.Printf("Received analytics event: %s", event.Event)
	})
	if err != nil {
		log.Printf("Analytics subscriber error: %v", err)
	}
}

func (h *Hub) startHeartbeatProcessor() {
	// Process heartbeat events
	err := h.redisClient.SubscribeHeartbeat(h.ctx, func(sessionID string, data map[string]interface{}) {
		log.Printf("Received heartbeat for session: %s", sessionID)
		
		// Update session activity in database if needed
		// This is where you'd implement session timeout logic
	})
	if err != nil {
		log.Printf("Heartbeat processor error: %v", err)
	}
}

// GetActiveClients returns the count of active clients
func (h *Hub) GetActiveClients() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.clients)
}

// GetSessionClient returns the client for a specific session
func (h *Hub) GetSessionClient(sessionID string) *Client {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return h.sessionClients[sessionID]
}

// BroadcastToAll sends a message to all connected clients
func (h *Hub) BroadcastToAll(message model.WSMessage) {
	data, err := json.Marshal(message)
	if err != nil {
		log.Printf("Failed to marshal broadcast message: %v", err)
		return
	}

	select {
	case h.broadcast <- data:
	default:
		log.Printf("Broadcast channel is full, dropping message")
	}
}

// Shutdown gracefully shuts down the hub
func (h *Hub) Shutdown() {
	log.Println("Shutting down WebSocket hub...")
	h.cancel()

	h.mutex.Lock()
	defer h.mutex.Unlock()

	// Close all client connections
	for client := range h.clients {
		close(client.send)
	}

	// Clear maps
	h.clients = make(map[*Client]bool)
	h.sessionClients = make(map[string]*Client)
}