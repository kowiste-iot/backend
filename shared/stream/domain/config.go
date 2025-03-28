package domain

import "time"

type StreamConfig struct {
    URL            string        // Connection URL (e.g., NATS server URL)
    MaxReconnects  int          // Maximum reconnection attempts
    ReconnectWait  time.Duration // Time to wait between reconnection attempts
    QoS            int          // Quality of Service level
    ConnectTimeout time.Duration // Connection timeout
    WriteTimeout   time.Duration // Write timeout for publishing messages
    PersistMessage bool         // Whether to persist messages before sending
}