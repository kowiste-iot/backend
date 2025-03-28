package app

import (
	"backend/shared/websocket/domain"
	"encoding/json"
	"fmt"
)

// Handle subscribe command
func (c *Client) handleSubscribe(msg *domain.Message) error {
	msgContent := new(domain.SubjectContent)
	// Convert the Content field to SubjectContent
	contentData, err := json.Marshal(msg.Content)
	if err != nil {
		return fmt.Errorf("error marshaling content: %w", err)
	}

	if err := json.Unmarshal(contentData, msgContent); err != nil {
		return fmt.Errorf("error unmarshaling to SubjectContent: %w", err)
	}

	c.hub.SubscribeToMeasures(c.tenantID, c.userID, msgContent.Measures)

	// Send confirmation back to client
	resp := map[string]interface{}{
		"type":     "subscribeConfirmed",
		"success":  true,
		"measures": msgContent.Measures,
	}
	respData, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	c.send <- respData
	return nil
}

// Handle unsubscribe command
func (c *Client) handleUnsubscribe(clientMsg *domain.Message) error {
	msgContent := new(domain.SubjectContent)

	c.hub.UnsubscribeFromMeasures(c.tenantID, c.userID, msgContent.Measures)

	// Send confirmation back to client
	resp := map[string]interface{}{
		"type":    "unsubscribeConfirmed",
		"success": true,
	}
	respData, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	c.send <- respData
	return nil
}

// Handle get current value command
func (c *Client) handleGetCurrentValue(clientMsg *domain.Message) error {
	// Example implementation:
	resp := map[string]interface{}{
		"type":    "currentValueResponse",
		"success": false,
		"message": "Not implemented",
	}
	respData, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	c.send <- respData
	return nil
}
