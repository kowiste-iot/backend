package app

import "backend/shared/stream/domain"

type WebSocketHandler struct {
	hub *Hub
}

func NewWebSocketHandler(hub *Hub) *WebSocketHandler {
	return &WebSocketHandler{hub: hub}
}

func (h *WebSocketHandler) Handle(msg domain.Message) (err error) {
	b, err := msg.Data.ToBytes()
	if err != nil {
		return
	}
	h.hub.SendToUser(msg.TenantID, msg.UserID, b)
	return
}
