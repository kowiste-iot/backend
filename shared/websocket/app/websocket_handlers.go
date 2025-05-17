package app

import (
	indestDomain "backend/internal/features/ingest/domain"
	"fmt"

	"backend/shared/stream/domain"
	wsDomain "backend/shared/websocket/domain"
	"context"
	"encoding/json"
	"time"
)

// handleStreamMessage processes messages from the stream service
func (s *WebSocketStreamService) handleStreamMessage(ctx context.Context, msg *domain.WireMessage) error {

	// Look for a handler for this message type
	s.mu.RLock()
	handler, exists := s.messageHandlers[msg.Topic]
	s.mu.RUnlock()

	if exists {
		return handler(ctx, msg)
	}

	// Default behavior for unknown message types
	s.base.Logger.Debug(ctx, fmt.Sprintf("No handler registered for message type: %s", msg.Topic),
		map[string]interface{}{"messageType": msg.Topic})

	return nil
}

// handleMeasureUpdate processes measure update messages
func (s *WebSocketStreamService) handleMeasureUpdate(ctx context.Context, msg *domain.WireMessage) error {
	var ingestMessage indestDomain.Message

	// Convert the wire message data to an ingest message
	err := msg.DataToModel(&ingestMessage)
	if err != nil {
		return err
	}

	// If no measure ID, cannot determine subscribers
	if ingestMessage.ID == "" {
		s.base.Logger.Debug(ctx, "Measure update doesn't contain measureID", nil)
		return nil
	}

	// Get all users subscribed to this measure
	subscribers := s.hub.GetSubscribedUsers(ingestMessage.TenantID, ingestMessage.BranchID, ingestMessage.ID)
	if len(subscribers) == 0 {
		// No subscribed users, nothing to do
		return nil
	}

	// Create a WebSocket message to send to subscribers
	wsMessage := &wsDomain.Message{
		TenantID:  ingestMessage.TenantID,
		Type:      "measure_update",
		Content:   msg.Data,
		CreatedAt: time.Now(),
	}

	// Convert to JSON once (optimization)
	jsonData, err := json.Marshal(wsMessage)
	if err != nil {
		return err
	}

	// Send to all subscribed users
	for _, userID := range subscribers {
		s.hub.SendToUser(ingestMessage.TenantID, ingestMessage.BranchID, userID, jsonData)
	}

	return nil
}

// handleDirectMessage processes direct messages to specific users
func (s *WebSocketStreamService) handleDirectMessage(ctx context.Context, msg *domain.WireMessage) error {
	var ingestMessage indestDomain.Message

	// Convert the wire message data to an ingest message
	err := msg.DataToModel(&ingestMessage)
	if err != nil {
		return err
	}
	// Extract the target user ID
	var userID string
	if userIDVal, ok := ingestMessage.Data["userID"]; ok {
		if userIDStr, ok := userIDVal.(string); ok {
			userID = userIDStr
		}
	}

	// If no user ID, cannot send the message
	if userID == "" {
		s.base.Logger.Debug(ctx, "Direct message doesn't contain userID", nil)
		return nil
	}

	// Create a WebSocket message
	wsMessage := &wsDomain.Message{
		TenantID:  ingestMessage.TenantID,
		UserID:    userID,
		Type:      "direct_message",
		Content:   msg.Data,
		CreatedAt: time.Now(),
	}

	// Convert to JSON
	jsonData, err := json.Marshal(wsMessage)
	if err != nil {
		return err
	}

	// Send to the specific user
	s.hub.SendToUser(ingestMessage.TenantID, ingestMessage.BranchID, userID, jsonData)
	return nil
}

// handleBroadcast processes broadcast messages to all users in a tenant
func (s *WebSocketStreamService) handleBroadcast(ctx context.Context, msg *domain.WireMessage) error {
	var ingestMessage indestDomain.Message

	// Convert the wire message data to an ingest message
	err := msg.DataToModel(&ingestMessage)
	if err != nil {
		return err
	}
	// Create a WebSocket message
	wsMessage := &wsDomain.Message{
		TenantID:  ingestMessage.TenantID,
		Type:      "broadcast",
		Content:   msg.Data,
		CreatedAt: time.Now(),
	}

	// Convert to JSON
	jsonData, err := json.Marshal(wsMessage)
	if err != nil {
		return err
	}

	// Get all connected users for this tenant
	users := s.hub.GetConnectedUsers(ingestMessage.TenantID)

	// Send to all users in the tenant
	for _, userID := range users {
		s.hub.SendToUser(ingestMessage.TenantID, ingestMessage.BranchID, userID, jsonData)
	}

	return nil
}
