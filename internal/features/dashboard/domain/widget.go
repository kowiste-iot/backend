package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type Widget struct {
	id          string
	tenantID    string
	branchName  string
	dashboardID string
	typeWidget  byte
	x           int
	y           int
	w           int
	h           int
	WidgetData
	updatedAt time.Time
	deletedAt *time.Time
}

type WidgetData struct {
	label       string
	showLabel   bool
	showEmotion bool
	trueEmotion bool
	link        []WidgetLinkData
	options     any
}

type WidgetLinkData struct {
	measure string
	tag     string
	legend  string
}

func NewWidget(tenantID, branchName, dashboardID string, typeWidget byte,  x, y, w, h int,
	label string, showLabel, showEmotion, trueEmotion bool, link []WidgetLinkData, options any) (widget *Widget, err error) {

	id, err := uuid.NewV7()
	if err != nil {
		return
	}

	widget = &Widget{
		id:          id.String(),
		branchName:  branchName,
		tenantID:    tenantID,
		dashboardID: dashboardID,
		typeWidget:  typeWidget,
		x:           x,
		y:           y,
		w:           w,
		h:           h,
		WidgetData: WidgetData{
			label:       label,
			showLabel:   showLabel,
			showEmotion: showEmotion,
			trueEmotion: trueEmotion,
			link:        link,
			options:     options,
		},
		updatedAt: time.Now(),
	}
	return
}

func NewWidgetFromRepository(id, tenantID, branchName, dashboardID string, typeWidget byte,
	 x, y, w, h int, widgetData WidgetData, updatedAt time.Time, deletedAt *time.Time) *Widget {
	return &Widget{
		id:          id,
		tenantID:    tenantID,
		branchName:  branchName,
		dashboardID: dashboardID,
		typeWidget:  typeWidget,
		x:           x,
		y:           y,
		w:           w,
		h:           h,
		WidgetData:  widgetData,
		updatedAt:   updatedAt,
		deletedAt:   deletedAt,
	}
}

func (a *Widget) Update( typeWidget byte,  x, y, w, h int, widgetData WidgetData) error {
	a.typeWidget = typeWidget
	a.x = x
	a.y = y
	a.w = w
	a.h = h
	a.WidgetData = widgetData
	a.updatedAt = time.Now()
	return nil
}

func (a *Widget) UpdatePosition( x, y, w, h int) error {
	a.x = x
	a.y = y
	a.w = w
	a.h = h
	a.updatedAt = time.Now()
	return nil
}

func (a *Widget) Delete() {
	now := time.Now()
	a.deletedAt = &now
}

func (a *Widget) IsDeleted() bool {
	return a.deletedAt != nil
}

// Add this to your domain/widget.go file
func NewWidgetData(label string, showLabel, showEmotion, trueEmotion bool) WidgetData {
	return WidgetData{
		label:       label,
		showLabel:   showLabel,
		showEmotion: showEmotion,
		trueEmotion: trueEmotion,
		link:        []WidgetLinkData{},
		options:     nil,
	}
}
func (wd *WidgetData) SetLink(links []WidgetLinkData) *WidgetData {
	wd.link = links
	return wd
}
func (wd *WidgetData) SetOptions(options any) *WidgetData {
	wd.options = options
	return wd
}
func NewWidgetLinkData(measure, tag, legend string) WidgetLinkData {
	return WidgetLinkData{
		measure: measure,
		tag:     tag,
		legend:  legend,
	}
}

// Getters for all fields
func (a *Widget) ID() string            { return a.id }
func (a *Widget) TenantID() string      { return a.tenantID }
func (a *Widget) BranchName() string    { return a.branchName }
func (a *Widget) DashboardID() string   { return a.dashboardID }
func (a *Widget) TypeWidget() byte      { return a.typeWidget }
func (a *Widget) X() int                { return a.x }
func (a *Widget) Y() int                { return a.y }
func (a *Widget) W() int                { return a.w }
func (a *Widget) H() int                { return a.h }
func (a *Widget) UpdatedAt() time.Time  { return a.updatedAt }
func (a *Widget) DeletedAt() *time.Time { return a.deletedAt }

// Getters for WidgetData
func (a *Widget) Label() string          { return a.label }
func (a *Widget) ShowLabel() bool        { return a.showLabel }
func (a *Widget) ShowEmotion() bool      { return a.showEmotion }
func (a *Widget) TrueEmotion() bool      { return a.trueEmotion }
func (a *Widget) Link() []WidgetLinkData { return a.link }
func (a *Widget) Options() any           { return a.options }

// Getters for WidgetLinkData
func (a *WidgetLinkData) MeasureID() string { return a.measure }
func (a *WidgetLinkData) Legend() string    { return a.legend }
func (a *WidgetLinkData) Tag() string       { return a.tag }

// Errors
var (
	ErrWidgetNotFound = errors.New("widget not found")
)
