package services

import (
	"backend/shared/websocket/app"
	"context"
)

func (c *Container) initializeWebsocketService(s *Services) (err error) {

	// Create WebSocket hub
	s.WebSocketHub = app.NewHub()
	stream, err := c.initializeStreamService(s)
	if err!=nil{
		return
	}
	// Create WebSocket stream service
	webSocketService := app.NewWebSocketStreamService(
		c.base,
		s.WebSocketHub,
		stream,
	)

	// Start the service
	err = webSocketService.Start(context.Background())
	if err != nil {
		return err
	}

	s.WebSocketService = webSocketService

	return
}
