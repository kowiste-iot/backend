package repository

import (
	"backend/internal/features/dashboard/domain"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type WidgetDB struct {
	ID          string           `gorm:"primaryKey;column:id"`
	DashboardID string           `gorm:"column:dashboard_id;not null"`
	TypeWidget  byte             `gorm:"column:type_widget;not null"`
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

func (WidgetDB) TableName() string {
	return "widgets"
}

// Mapper domain model to db model
func widgetMapper(data *domain.Widget) (dbWidget *WidgetDB) {
	lMap := make([]WidgetLinkData, 0)
	for _, link := range data.Link() {
		lMap = append(lMap, WidgetLinkData{
			WidgetID: data.ID(),
			Measure:  link.MeasureID(),
			Tag:      link.Tag(),
			Legend:   link.Legend(),
		})
	}
	dbWidget = &WidgetDB{
		ID:          data.ID(),
		DashboardID: data.DashboardID(),
		TypeWidget:  data.TypeWidget(),
		X:           data.X(),
		Y:           data.Y(),
		W:           data.W(),
		H:           data.H(),
		Label:       data.Label(),
		ShowLabel:   data.ShowLabel(),
		ShowEmotion: data.ShowEmotion(),
		TrueEmotion: data.TrueEmotion(),
		Link:        lMap,
	}
	return
}
