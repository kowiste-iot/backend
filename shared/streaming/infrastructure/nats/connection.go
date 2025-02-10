package nats

import (
	"backend/shared/streaming/domain"

	"github.com/nats-io/nats.go"
)

type Connection struct {
	nc     *nats.Conn
	config domain.ConnectionConfig
}

func NewConnection(config domain.ConnectionConfig) *Connection {
	return &Connection{
		config: config,
	}
}

func (c *Connection) Connect() error {
	opts := []nats.Option{
		nats.Name("Your Service Name"),
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(-1),
	}

	if c.config.Username != "" && c.config.Password != "" {
		opts = append(opts, nats.UserInfo(c.config.Username, c.config.Password))
	}

	nc, err := nats.Connect(c.config.URL, opts...)
	if err != nil {
		return err
	}

	c.nc = nc
	return nil
}

func (c *Connection) Close() error {
	if c.nc != nil {
		return c.nc.Drain()
	}
	return nil
}

func (c *Connection) IsConnected() bool {
	return c.nc != nil && c.nc.IsConnected()
}

// This method is specific to NATS implementation
func (c *Connection) GetConn() *nats.Conn {
	return c.nc
}
