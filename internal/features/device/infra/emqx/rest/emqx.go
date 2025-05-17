package emqx

import (
	"backend/pkg/config"
	"context"
	"fmt"
	
	"github.com/go-resty/resty/v2"
)

type DeviceBroker struct {
	cfg *brokerConfig
	client *resty.Client
}

type brokerConfig struct {
	URL    string
	Key    string
	Secret string
}

func NewDeviceBroker(config *config.IngestConfig) *DeviceBroker {
	client := resty.New()
	
	// Set basic auth credentials
	client.SetBasicAuth(config.ManageKey, config.ManageSecret)
	
	return &DeviceBroker{
		cfg: &brokerConfig{
			URL:    config.ManageURL,
			Key:    config.ManageKey,
			Secret: config.ManageSecret,
		},
		client: client,
	}
}

func (d *DeviceBroker) KickOut(ctx context.Context, clientID string) error {
	url := fmt.Sprintf("%s/api/v5/clients/%s", d.cfg.URL, clientID)
	
	resp, err := d.client.R().
		SetContext(ctx).
		Delete(url)
	
	if err != nil {
		return fmt.Errorf("failed to make request to EMQX: %w", err)
	}
	
	if resp.StatusCode() >= 400 {
		return fmt.Errorf("failed to kick out client, status code: %d, response: %s", 
			resp.StatusCode(), string(resp.Body()))
	}
	
	return nil
}
