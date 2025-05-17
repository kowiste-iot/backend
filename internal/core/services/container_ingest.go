package services

import (
	"backend/internal/features/ingest/app"
	"backend/internal/features/ingest/domain"
	"backend/internal/features/ingest/interface/mqtt"

	"fmt"
)

func (c *Container) initializeIngestService(s *Services) (err error) {

	// Initialize MQTT consumer config
	mqttConfig := &mqtt.Config{
		BrokerURL: c.config.Ingest.URL, // Default EMQX port
		ClientID:  c.config.Ingest.ClientID,
		Username:  c.config.Ingest.User,
		Password:  c.config.Ingest.Password,
		Topics:    []string{"devices/#"},
	}

	// Create MQTT consumer
	consumer, err := mqtt.NewConsumer(mqttConfig)
	if err != nil {
		return fmt.Errorf("failed to create MQTT consumer: %w", err)
	}

	// Create ingest service
	serviceConfig := &app.ServiceConfig{
		Topic:           domain.TopicIngest,
		PersistMessages: true,
	}
	stream, err := c.initializeStreamService(s)
	if err != nil {
		return
	}
	s.IngestService = app.NewService(stream, serviceConfig)
	s.IngestService.AddConsumer(consumer)
	err = s.IngestService.Start()
	if err != nil {
		return fmt.Errorf("failed to start MQTT consumer: %w", err)
	}
	return nil
}
