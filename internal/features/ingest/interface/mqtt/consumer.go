package mqtt

import (
	"backend/internal/features/ingest/domain"
	"encoding/json"
	"fmt"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type Config struct {
	BrokerURL string
	ClientID  string
	Username  string
	Password  string
	Topics    []string
}

type Consumer struct {
	client    mqtt.Client
	handler   func(msg *domain.Message) error
	config    *Config
	opts      *mqtt.ClientOptions
	connected bool
}

func NewConsumer(config *Config) (*Consumer, error) {
	opts := mqtt.NewClientOptions().
		AddBroker(config.BrokerURL).
		SetClientID(config.ClientID).
		SetUsername(config.Username).
		SetPassword(config.Password).
		SetAutoReconnect(true).
		SetConnectionLostHandler(func(client mqtt.Client, err error) {
			fmt.Printf("Connection lost: %v\n", err)
		})

	return &Consumer{
		config: config,
		opts:   opts,
	}, nil
}

func (c *Consumer) Start() error {
	c.client = mqtt.NewClient(c.opts)

	if token := c.client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect: %w", token.Error())
	}

	for _, topic := range c.config.Topics {
		if token := c.client.Subscribe(topic, 0, c.handleMessage); token.Wait() && token.Error() != nil {
			return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
		}
	}

	c.connected = true
	return nil
}

func (c *Consumer) Stop() error {
	if c.connected && c.client != nil {
		for _, topic := range c.config.Topics {
			if token := c.client.Unsubscribe(topic); token.Wait() && token.Error() != nil {
				fmt.Printf("Error unsubscribing from %s: %v\n", topic, token.Error())
			}
		}

		c.client.Disconnect(250)
		c.connected = false
	}
	return nil
}

func (c *Consumer) Subscribe(handler func(msg *domain.Message) error) error {
	c.handler = handler
	return nil
}

func (c *Consumer) handleMessage(_ mqtt.Client, msg mqtt.Message) {
	var message domain.Message

	if err := json.Unmarshal(msg.Payload(), &message); err != nil {
		fmt.Printf("Error unmarshaling message: %v\n", err)
		return
	}
	if message.Validate() != nil {
		return 
	}

	if c.handler != nil {
		if err := c.handler(&message); err != nil {
			fmt.Printf("Error handling message: %v\n", err)
		}
	}
}
