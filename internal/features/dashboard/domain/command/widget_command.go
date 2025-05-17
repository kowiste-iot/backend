package command

import baseCmd "backend/shared/base/command"


// CreateWidgetInput represents input data for creating a widget
type CreateWidgetInput struct {
	baseCmd.BaseInput
	DashboardID string `json:"dashboardID" validate:"required"`
	TypeWidget  byte   `json:"type"`
	// Widget data fields
	Label       string       `json:"label"`
	ShowLabel   bool         `json:"showLabel"`
	ShowEmotion bool         `json:"showEmotion"`
	TrueEmotion bool         `json:"trueEmotion"`
	Link        []WidgetLink `json:"link"`
	Options     any          `json:"options"`
}

// WidgetLink represents a link in widget data
type WidgetLink struct {
	Measure string `json:"measure"`
	Tag     string `json:"tag"`
	Legend  string `json:"legend"`
}

// UpdateWidgetInput represents input data for updating a widget
type UpdateWidgetInput struct {
	baseCmd.BaseInput
	ID          string `json:"id" validate:"required"`
	DashboardID string `json:"dashboardID" validate:"required"`
	TypeWidget  byte   `json:"type"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	W           int    `json:"w"`
	H           int    `json:"h"`

	// Widget data fields
	Label       string       `json:"label"`
	ShowLabel   bool         `json:"showLabel"`
	ShowEmotion bool         `json:"showEmotion"`
	TrueEmotion bool         `json:"trueEmotion"`
	Link        []WidgetLink `json:"link"`
	Options     any          `json:"options"`
}

// UpdateWidgetPositionInput represents input for updating just widget position
type UpdateWidgetPositionInput struct {
	baseCmd.BaseInput
	ID          string `json:"id" validate:"required"`
	DashboardID string `json:"dashboardID" validate:"required"`
	X           int    `json:"x"`
	Y           int    `json:"y"`
	W           int    `json:"w"`
	H           int    `json:"h"`
}

// WidgetIDInput represents input for getting a widget by ID
type WidgetIDInput struct {
	baseCmd.BaseInput
	DashboardID string `json:"dashboardID" validate:"required"`
	WidgetID    string `json:"widgetID" validate:"required"`
}
