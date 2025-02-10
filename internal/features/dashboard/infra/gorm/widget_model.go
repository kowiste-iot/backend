package repository

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type Widget struct {
	ID          string           `gorm:"primaryKey;column:id"`
	TenantID    string           `gorm:"column:tenant_id;not null"`
	BranchName  string           `gorm:"column:branch_name;not null"`
	DashboardID string           `gorm:"column:dashboard_id;not null"`
	Name        string           `gorm:"column:name;not null"`
	TypeWidget  byte             `gorm:"column:type_widget;not null"`
	I           int              `gorm:"column:i;not null"`
	X           int              `gorm:"column:x;not null"`
	Y           int              `gorm:"column:y;not null"`
	W           int              `gorm:"column:w;not null"`
	H           int              `gorm:"column:h;not null"`
	Label       string           `gorm:"column:label"`
	ShowLabel   bool             `gorm:"column:show_label"`
	ShowEmotion bool             `gorm:"column:show_emotion"`
	TrueEmotion bool             `gorm:"column:true_emotion"`
	Link        []WidgetLinkData `gorm:"foreignKey:WidgetID;references:ID"`
	Options     JSON             `gorm:"column:options;type:text"`
	UpdatedAt   time.Time        `gorm:"column:updated_at;not null"`
	DeletedAt   *time.Time       `gorm:"column:deleted_at"`
}


type WidgetLinkData struct {
	ID       uint   `gorm:"primaryKey"`
	WidgetID string `gorm:"column:widget_id;not null"`
	Measure  string `gorm:"column:measure;not null"`
	Tag      string `gorm:"column:tag;not null"`
	Legend   string `gorm:"column:legend;not null"`
}

// Custom JSON type to handle any JSON data
type JSON json.RawMessage

// Scan implements the sql.Scanner interface
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	s, ok := value.([]byte)
	if !ok {
		return errors.New("invalid scan source for JSON")
	}

	*j = append((*j)[0:0], s...)
	return nil
}

// Value implements the driver.Valuer interface
func (j JSON) Value() (driver.Value, error) {
	if j == nil {
		return nil, nil
	}
	return string(j), nil
}

func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

func (j *JSON) UnmarshalJSON(data []byte) error {
	*j = append((*j)[0:0], data...)
	return nil
}
