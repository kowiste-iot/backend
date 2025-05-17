// internal/features/device/infra/gorm/mqtt/repository.go
package mqtt

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"gorm.io/gorm"
)

type MqttRepository interface {
	CreateMqttUser(ctx context.Context, deviceID, password string) error
	DeleteMqttUser(ctx context.Context, deviceID string) error
}

type mqttRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) MqttRepository {
	return &mqttRepository{db: db}
}

// hashPassword creates a SHA256 hash of the password
func hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (r *mqttRepository) CreateMqttUser(ctx context.Context, deviceID, password string) error {
	hashedPassword := hashPassword(password)

	user := newMQTTUser(deviceID, hashedPassword)

	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return fmt.Errorf("failed to create MQTT user: %w", err)
	}

	// Create ACL rules for this device
	// Rule for publishing to its own topic
	publishAcl := newMQTTAclPublish(user.Username, deviceID)

	// Rule for subscribing to its own topic
	subscribeAcl := newMQTTAclSubscribe(user.Username, deviceID)

	if err := r.db.WithContext(ctx).Create(&publishAcl).Error; err != nil {
		return fmt.Errorf("failed to create MQTT publish ACL: %w", err)
	}

	if err := r.db.WithContext(ctx).Create(&subscribeAcl).Error; err != nil {
		return fmt.Errorf("failed to create MQTT subscribe ACL: %w", err)
	}

	return nil
}

func (r *mqttRepository) DeleteMqttUser(ctx context.Context, deviceID string) error {
	username := fmt.Sprintf("device_%s", deviceID)

	// Delete ACL rules first
	if err := r.db.WithContext(ctx).Where("username = ?", username).Delete(&MqttAclDB{}).Error; err != nil {
		return fmt.Errorf("failed to delete MQTT ACLs: %w", err)
	}

	// Delete the user
	if err := r.db.WithContext(ctx).Where("username = ?", username).Delete(&MqttUserDB{}).Error; err != nil {
		return fmt.Errorf("failed to delete MQTT user: %w", err)
	}

	return nil
}
