package dashboardhandler

import (
	"backend/internal/features/dashboard/domain"
)

// CreateWidgetRequest maps to the frontend IWidget structure
type CreateWidgetRequest struct {
	DashboardID string            `json:"dashboardID" binding:"required"`
	Type        byte              `json:"type" binding:"required" example:"1"`
	Data        WidgetDataRequest `json:"data"`
}

// WidgetDataRequest maps to the frontend IWidgetData structure
type WidgetDataRequest struct {
	Label       string              `json:"label" example:"Temperature Readings"`
	ShowLabel   bool                `json:"showLabel" example:"true"`
	ShowEmotion bool                `json:"showEmotion" example:"false"`
	TrueEmotion bool                `json:"trueEmotion" example:"false"`
	Link        []WidgetLinkRequest `json:"link"`
	Options     any                 `json:"options"`
}

// WidgetLinkRequest maps to the frontend IWidgetLinkData structure
type WidgetLinkRequest struct {
	Measure string `json:"measure" example:"sensor1"`
	Tag     string `json:"tag" example:"temperature"`
	Legend  string `json:"legend" example:"Room Temp"`
}

// UpdateWidgetRequest for updating widgets
type UpdateWidgetRequest struct {
	Type byte              `json:"type" example:"1"`
	I    int               `json:"i" example:"0"`
	X    int               `json:"x" example:"0"`
	Y    int               `json:"y" example:"0"`
	W    int               `json:"w" example:"6"`
	H    int               `json:"h" example:"4"`
	Data WidgetDataRequest `json:"data"`
}

// UpdateWidgetPositionRequest for updating just widget position
type UpdateWidgetPositionRequest struct {
	X int `json:"x" example:"1"`
	Y int `json:"y" example:"2"`
	W int `json:"w" example:"6"`
	H int `json:"h" example:"4"`
}

// WidgetResponse maps to the frontend IWidget structure
type WidgetResponse struct {
	ID          string             `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	DashboardID string             `json:"dashboardID" example:"550e8400-e29b-41d4-a716-446655440001"`
	Type        byte               `json:"type" example:"1"`
	I           int                `json:"i" example:"0"`
	X           int                `json:"x" example:"0"`
	Y           int                `json:"y" example:"0"`
	W           int                `json:"w" example:"6"`
	H           int                `json:"h" example:"4"`
	Data        WidgetDataResponse `json:"data"`
	UpdatedAt   int64              `json:"updatedAt" example:"1615482896"`
}

// WidgetDataResponse maps to the frontend IWidgetData structure
type WidgetDataResponse struct {
	Label       string               `json:"label" example:"Temperature Readings"`
	ShowLabel   bool                 `json:"showLabel" example:"true"`
	ShowEmotion bool                 `json:"showEmotion" example:"false"`
	TrueEmotion bool                 `json:"trueEmotion" example:"false"`
	Link        []WidgetLinkResponse `json:"link"`
	Options     any                  `json:"options"`
}

// WidgetLinkResponse maps to the frontend IWidgetLinkData structure
type WidgetLinkResponse struct {
	Measure string `json:"measure" example:"sensor1"`
	Tag     string `json:"tag" example:"temperature"`
	Legend  string `json:"legend" example:"Room Temp"`
}

// ToWidgetResponse converts domain Widget to handler WidgetResponse
func ToWidgetResponse(w *domain.Widget, num int) WidgetResponse {
	linkResponses := make([]WidgetLinkResponse, len(w.Link()))
	for i, link := range w.Link() {
		linkResponses[i] = WidgetLinkResponse{
			Measure: link.MeasureID(),
			Tag:     link.Tag(),
			Legend:  link.Legend(),
		}
	}

	return WidgetResponse{
		ID:          w.ID(),
		DashboardID: w.DashboardID(),
		Type:        w.TypeWidget(),
		I:           num, //addin a index for the web
		X:           w.X(),
		Y:           w.Y(),
		W:           w.W(),
		H:           w.H(),
		Data: WidgetDataResponse{
			Label:       w.Label(),
			ShowLabel:   w.ShowLabel(),
			ShowEmotion: w.ShowEmotion(),
			TrueEmotion: w.TrueEmotion(),
			Link:        linkResponses,
			Options:     w.Options(),
		},
		UpdatedAt: w.UpdatedAt().Unix(),
	}
}

// ToWidgetResponses converts domain Widget slice to handler WidgetResponse slice
func ToWidgetResponses(widgets []*domain.Widget) []WidgetResponse {
	responses := make([]WidgetResponse, len(widgets))
	for i, w := range widgets {
		responses[i] = ToWidgetResponse(w, i)
	}
	return responses
}
