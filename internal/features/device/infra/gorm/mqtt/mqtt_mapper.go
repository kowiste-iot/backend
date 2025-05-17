// internal/features/device/infra/gorm/mqtt/models.go
package mqtt

import (
	"fmt"
	"time"
)

type MqttUserDB struct {
	ID          int `gorm:"primaryKey;autoIncrement"`
	Username    string
	ClientID    string
	Password    string
	Salt        string
	IsSuperuser bool
	Created     time.Time
}

func (MqttUserDB) TableName() string {
	return "mqtt.mqtt_users"
}
func newMQTTUser(deviceID, password string) *MqttUserDB {
	return &MqttUserDB{
		Username:    fmt.Sprintf("device_%s", deviceID),
		Password:    password,
		ClientID:    deviceID,
		IsSuperuser: false,
		Created:     time.Now(),
	}
}

type MqttAclDB struct {
	ID       int `gorm:"primaryKey;autoIncrement"`
	Allow    int `gorm:"default:1"`
	Username string
	ClientID string
	Access   int `gorm:"default:1"`
	Topic    string
}

func (MqttAclDB) TableName() string {
	return "mqtt.mqtt_acls"
}
func newMQTTAclPublish(userName, deviceID string) *MqttAclDB {
	return &MqttAclDB{
		Allow:    1,
		Username: userName,
		Access:   2, // 2 is for publish
		ClientID: userName,
		Topic:    fmt.Sprintf("devices/%s/#", deviceID),
	}
}
func newMQTTAclSubscribe(userName, deviceID string) *MqttAclDB {
	return &MqttAclDB{
		Allow:    1,
		Username: userName,
		Access:   1, // 1 is for subscribe
		ClientID: userName,
		Topic:    fmt.Sprintf("devices/%s/#", deviceID),
	}
}
