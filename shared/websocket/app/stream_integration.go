package app

import (
	ingestDomain "backend/internal/features/ingest/domain"
	"backend/shared/base"
	"backend/shared/stream/domain"

	wsDomain "backend/shared/websocket/domain"
	"context"
	"encoding/json"
	"errors"
	"sync"
	"time"
)

// WebSocketStreamService integrates the existing WebSocket hub with the streaming service
type WebSocketStreamService struct {
	base         *base.BaseService
	hub          *Hub
	streamClient domain.StreamClient
	ctx          context.Context
	cancel       context.CancelFunc
	mu           sync.RWMutex
	isRunning    bool
	// Map of message types to handler functions
	messageHandlers map[string]wsDomain.MessageHandler
}

// NewWebSocketStreamService creates a new WebSocket stream service
func NewWebSocketStreamService(
	base *base.BaseService,
	hub *Hub,
	streamClient domain.StreamClient,
) *WebSocketStreamService {
	ctx, cancel := context.WithCancel(context.Background())

	service := &WebSocketStreamService{
		base:            base,
		hub:             hub,
		streamClient:    streamClient,
		ctx:             ctx,
		cancel:          cancel,
		messageHandlers: make(map[string]wsDomain.MessageHandler),
	}

	// Register default handlers
	service.registerDefaultHandlers()

	return service
}

// registerDefaultHandlers sets up the default message type handlers
func (s *WebSocketStreamService) registerDefaultHandlers() {
	// Handler for measure updates
	s.messageHandlers[ingestDomain.TopicIngest] = s.handleMeasureUpdate
	// Handler for direct messages
	s.messageHandlers[ingestDomain.TopicMessageDirect] = s.handleDirectMessage
	// Handler for broadcast messages
	s.messageHandlers[ingestDomain.TopicMessageBroadcast] = s.handleBroadcast
}

// Start initializes the WebSocket stream service
func (s *WebSocketStreamService) Start(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.isRunning {
		return errors.New("WebSocket stream service already running")
	}

	// Start the hub if it's not already running
	go s.hub.Run(s.ctx)

	// Subscribe to WebSocket messages from the stream
	err := s.streamClient.Subscribe(ctx, ingestDomain.TopicIngest, s.handleStreamMessage)
	if err != nil {
		return err
	}

	s.isRunning = true
	return nil
}

// Stop shuts down the WebSocket stream service
func (s *WebSocketStreamService) Stop() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.isRunning {
		return nil
	}

	// Unsubscribe from the stream
	err := s.streamClient.Unsubscribe(s.ctx, "websocket.messages")
	if err != nil {
		s.base.Logger.Error(s.ctx, err, "Failed to unsubscribe from WebSocket messages", nil)
	}

	// Cancel the context
	s.cancel()

	s.isRunning = false
	return nil
}

// RegisterHandler allows external code to register a custom handler for a message type
func (s *WebSocketStreamService) RegisterHandler(messageType string, handler wsDomain.MessageHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.messageHandlers[messageType] = handler
}

// SendMessage sends a message through the stream service
func (s *WebSocketStreamService) SendMessage(msg *wsDomain.Message) error {
	// Create a MessageData implementation
	wsMessageData := &WebSocketMessageData{
		Message: msg,
	}

	// Create a stream message
	streamMsg := &domain.Message{
		TenantID:  msg.TenantID,
		UserID:    msg.UserID,
		Topic:     "websocket.messages",
		Data:      wsMessageData,
		Timestamp: time.Now(),
		Event:     "websocket.message",
	}

	// Publish to the stream
	return s.streamClient.Publish(s.ctx, streamMsg)
}

// WebSocketMessageData implements the MessageData interface for WebSocket messages
type WebSocketMessageData struct {
	Message *wsDomain.Message
}

func (d *WebSocketMessageData) Validate() error {
	if d.Message == nil {
		return errors.New("message cannot be nil")
	}
	if d.Message.TenantID == "" {
		return errors.New("tenant ID cannot be empty")
	}
	if d.Message.UserID == "" {
		return errors.New("user ID cannot be empty")
	}
	return nil
}

func (d *WebSocketMessageData) ToBytes() ([]byte, error) {
	return json.Marshal(d.Message)
}
