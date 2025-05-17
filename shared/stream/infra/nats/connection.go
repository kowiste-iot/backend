package nats

import (
    "backend/shared/stream/domain"
    "fmt"
    "github.com/nats-io/nats.go"
    "sync"
)

type Connection struct {
    conn      *nats.Conn
    config    *domain.StreamConfig
    mu        sync.RWMutex
    connected bool
}

func NewConnection(config *domain.StreamConfig) *Connection {
    return &Connection{
        config: config,
    }
}

func (c *Connection) Connect() error {
    c.mu.Lock()
    defer c.mu.Unlock()

    options := []nats.Option{
        nats.MaxReconnects(c.config.MaxReconnects),
        nats.ReconnectWait(c.config.ReconnectWait),
        nats.Timeout(c.config.ConnectTimeout),
    }

    conn, err := nats.Connect(c.config.URL, options...)
    if err != nil {
        return fmt.Errorf("failed to connect to NATS: %w", err)
    }

    c.conn = conn
    c.connected = true
    return nil
}

func (c *Connection) Close() error {
    c.mu.Lock()
    defer c.mu.Unlock()

    if c.conn != nil {
        c.conn.Close()
        c.connected = false
    }
    return nil
}

func (c *Connection) IsConnected() bool {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.connected
}

func (c *Connection) GetConn() *nats.Conn {
    c.mu.RLock()
    defer c.mu.RUnlock()
    return c.conn
}