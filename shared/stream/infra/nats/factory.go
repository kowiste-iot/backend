package nats

import "backend/shared/stream/domain"

type NatsClientFactory struct {
    repository domain.MessageRepository
}

func NewNatsClientFactory(repo domain.MessageRepository) *NatsClientFactory {
    return &NatsClientFactory{
        repository: repo,
    }
}

func (f *NatsClientFactory) CreateClient(config *domain.StreamConfig) (domain.StreamClient, error) {
    conn := NewConnection(config)
    client := NewClient(conn, config, f.repository)
    
    if err := client.Connect(); err != nil {
        return nil, err
    }
    
    return client, nil
}