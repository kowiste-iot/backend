package repository

import (
    "backend/shared/stream/domain"
    "encoding/json"
    "gorm.io/gorm"
    "time"
)

type MessageModel struct {
    ID        string `gorm:"primarykey"`
    Topic     string
    Data      []byte 
    Event     string
    Status    string
    Timestamp time.Time
    UpdatedAt time.Time
}

type MessageRepository struct {
    db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
    return &MessageRepository{
        db: db,
    }
}

func (r *MessageRepository) Save(message *domain.Message) error {
    data, err := json.Marshal(message.Data)
    if err != nil {
        return err
    }

    model := &MessageModel{
        ID:        message.ID,
        Topic:     message.Topic,
        Data:      data,
        Event:     message.Event,
        Status:    string(message.Status),
        Timestamp: message.Timestamp,
    }

    return r.db.Create(model).Error
}

func (r *MessageRepository) UpdateStatus(id string, status domain.MessageStatus) error {
    return r.db.Model(&MessageModel{}).
        Where("id = ?", id).
        Update("status", string(status)).
        Error
}

func (r *MessageRepository) GetPendingMessages() ([]*domain.StoredMessage, error) {
    var models []MessageModel
    
    err := r.db.Where("status = ?", string(domain.MessageStatusPending)).Find(&models).Error
    if err != nil {
        return nil, err
    }

    messages := make([]*domain.StoredMessage, len(models))
    for i, model := range models {
        messages[i] = &domain.StoredMessage{
            ID:        model.ID,
            Topic:     model.Topic,
            Data:      model.Data,
            Event:     model.Event,
            Status:    domain.MessageStatus(model.Status),
            Timestamp: model.Timestamp,
        }
    }

    return messages, nil
}