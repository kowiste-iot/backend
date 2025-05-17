package app

import (
	"backend/internal/features/dashboard/domain"
	"backend/internal/features/dashboard/domain/command"
	resourceDomain "backend/internal/features/resource/domain"
	scopeDomain "backend/internal/features/scope/domain"
	"backend/shared/base"
	baseCmd "backend/shared/base/command"
	"backend/shared/validator"
	"context"
	"fmt"
)

type WidgetService interface {
	CreateWidget(ctx context.Context, input *command.CreateWidgetInput) (*domain.Widget, error)
	GetWidget(ctx context.Context, input *command.WidgetIDInput) (*domain.Widget, error)
	ListWidgets(ctx context.Context, input *command.DashboardIDInput) ([]*domain.Widget, error)
	UpdateWidget(ctx context.Context, input *command.UpdateWidgetInput) (*domain.Widget, error)
	UpdateWidgetPosition(ctx context.Context, input *command.UpdateWidgetPositionInput) (*domain.Widget, error)
	DeleteWidget(ctx context.Context, input *command.WidgetIDInput) error
}

type widgetService struct {
	repo domain.WidgetRepository
	*base.BaseService
}

func NewWidgetService(base *base.BaseService, repo domain.WidgetRepository) WidgetService {
	return &widgetService{
		repo:        repo,
		BaseService: base,
	}
}

func (s *widgetService) CreateWidget(ctx context.Context, input *command.CreateWidgetInput) (widget *domain.Widget, err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Widget,
		Scope:     scopeDomain.Create,
	})
	if err != nil {
		return nil, err
	}

	err = validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}

	requiredWidth := 4
	requiredHeight := 3

	// Find an empty space for the new widget
	x, y, err := s.findEmptySpace(ctx, input.DashboardID, input.TenantDomain, input.BranchName, requiredWidth, requiredHeight)
	if err != nil {
		return nil, fmt.Errorf("failed to find position for widget: %w", err)
	}

	// Convert link data to domain WidgetLinkData
	linkData := make([]domain.WidgetLinkData, len(input.Link))
	for i, link := range input.Link {
		linkData[i] = domain.NewWidgetLinkData(link.Measure, link.Tag, link.Legend)
	}

	widget, err = domain.NewWidget(
		input.TenantDomain, input.BranchName,
		input.DashboardID, input.TypeWidget,
		x, y,
		requiredWidth, requiredHeight,
		input.Label,
		input.ShowLabel,
		input.ShowEmotion,
		input.TrueEmotion,
		linkData,
		input.Options,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create widget: %w", err)
	}

	err = s.repo.Create(ctx, widget)
	if err != nil {
		return nil, fmt.Errorf("failed to create widget: %w", err)
	}

	return widget, nil
}

func (s *widgetService) GetWidget(ctx context.Context, input *command.WidgetIDInput) (widget *domain.Widget, err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Widget,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}

	widget, err = s.repo.FindByID(ctx, input.DashboardID, input.WidgetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get widget: %w", err)
	}

	return widget, nil
}

func (s *widgetService) ListWidgets(ctx context.Context, input *command.DashboardIDInput) (widgets []*domain.Widget, err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Widget,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}

	widgets, err = s.repo.FindAll(ctx, input.DashboardID)
	if err != nil {
		return nil, fmt.Errorf("failed to list widgets: %w", err)
	}

	return widgets, nil
}

func (s *widgetService) UpdateWidget(ctx context.Context, input *command.UpdateWidgetInput) (widget *domain.Widget, err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Widget,
		Scope:     scopeDomain.Update,
	})
	if err != nil {
		return nil, err
	}

	widget, err = s.repo.FindByID(ctx, input.DashboardID, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get widget: %w", err)
	}

	// Convert link data to domain WidgetLinkData
	linkData := make([]domain.WidgetLinkData, len(input.Link))
	for i, link := range input.Link {
		linkData[i] = domain.NewWidgetLinkData(link.Measure, link.Tag, link.Legend)
	}

	// Create widget data
	wData := domain.NewWidgetData(input.Label, input.ShowLabel, input.ShowEmotion, input.TrueEmotion)
	wData.SetOptions(input.Options)
	lTemp := make([]domain.WidgetLinkData, 0)
	for i := range input.Link {
		lTemp = append(lTemp, domain.NewWidgetLinkData(input.Link[i].Measure, input.Link[i].Tag, input.Link[i].Legend))
	}
	wData.SetLink(lTemp)
	err = widget.Update(
		input.TypeWidget,
		input.X, input.Y,
		input.W, input.H,
		wData,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update widget: %w", err)
	}

	if err := s.repo.Update(ctx, widget); err != nil {
		return nil, fmt.Errorf("failed to update widget: %w", err)
	}

	return widget, nil
}

func (s *widgetService) UpdateWidgetPosition(ctx context.Context, input *command.UpdateWidgetPositionInput) (widget *domain.Widget, err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Widget,
		Scope:     scopeDomain.Update,
	})
	if err != nil {
		return nil, err
	}

	widget, err = s.repo.FindByID(ctx, input.DashboardID, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get widget: %w", err)
	}

	err = widget.UpdatePosition(
		input.X, input.Y,
		input.W, input.H,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update widget position: %w", err)
	}

	if err := s.repo.Update(ctx, widget); err != nil {
		return nil, fmt.Errorf("failed to update widget: %w", err)
	}

	return widget, nil
}

func (s *widgetService) DeleteWidget(ctx context.Context, input *command.WidgetIDInput) (err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Widget,
		Scope:     scopeDomain.Delete,
	})
	if err != nil {
		return err
	}

	err = s.repo.Remove(ctx, input.DashboardID, input.WidgetID)
	if err != nil {
		return err
	}

	return nil
}

func (s *widgetService) findEmptySpace(ctx context.Context, dashboardID, tenantID, branchName string, requiredWidth, requiredHeight int) (int, int, error) {

	widgets, err := s.repo.FindAll(ctx, dashboardID)
	if err != nil {
		return 0, 0, fmt.Errorf("error finding widgets: %w", err)
	}

	// Maximum columns in the grid
	const maxColumns = 24

	// If no widgets exist, place at the top
	if len(widgets) == 0 {
		return 0, 0, nil
	}

	// Create a grid representation of the dashboard
	// gridHeight will be determined by the maximum y + height of existing widgets
	gridHeight := 0

	for _, w := range widgets {
		bottomEdge := w.Y() + w.H()
		if bottomEdge > gridHeight {
			gridHeight = bottomEdge
		}
	}

	// Add some buffer for future widgets
	gridHeight += 10

	// Create a grid to track occupied spaces
	grid := make([][]bool, gridHeight)
	for i := range grid {
		grid[i] = make([]bool, maxColumns)
	}

	// Mark occupied spaces in the grid
	for _, w := range widgets {
		for i := w.Y(); i < w.Y()+w.H(); i++ {
			for j := w.X(); j < w.X()+w.W(); j++ {
				if i < gridHeight && j < maxColumns {
					grid[i][j] = true
				}
			}
		}
	}

	// Try to find a valid position
	// First check the top row
	for x := 0; x <= maxColumns-requiredWidth; x++ {
		valid := true
		for j := x; j < x+requiredWidth; j++ {
			if grid[0][j] {
				valid = false
				break
			}
		}

		if valid {
			// Found a valid position at the top
			return x, 0, nil
		}
	}

	// If no position at the top, look for positions touching other widgets
	for y := 1; y < gridHeight; y++ {
		for x := 0; x <= maxColumns-requiredWidth; x++ {
			// Check if this position is empty
			valid := true
			for i := y; i < y+requiredHeight && valid; i++ {
				for j := x; j < x+requiredWidth && valid; j++ {
					if i >= gridHeight || j >= maxColumns || grid[i][j] {
						valid = false
					}
				}
			}

			if !valid {
				continue
			}

			// Check if it's touching another widget from above
			touching := false
			for j := x; j < x+requiredWidth; j++ {
				if y > 0 && grid[y-1][j] {
					touching = true
					break
				}
			}

			// Check if it's touching from the left
			for i := y; i < y+requiredHeight && !touching; i++ {
				if x > 0 && grid[i][x-1] {
					touching = true
					break
				}
			}

			// Check if it's touching from the right
			for i := y; i < y+requiredHeight && !touching; i++ {
				if x+requiredWidth < maxColumns && grid[i][x+requiredWidth] {
					touching = true
					break
				}
			}

			if touching {
				return x, y, nil
			}
		}
	}

	// If no valid position found, place it at the bottom
	// Find the first available position at the bottom row
	bottomY := 0
	for _, w := range widgets {
		if w.Y()+w.H() > bottomY {
			bottomY = w.Y() + w.H()
		}
	}

	// Check if there's enough space at the bottom
	for x := 0; x <= maxColumns-requiredWidth; x++ {
		valid := true
		for j := x; j < x+requiredWidth; j++ {
			if bottomY < gridHeight && grid[bottomY][j] {
				valid = false
				break
			}
		}

		if valid {
			return x, bottomY, nil
		}
	}

	// If still no valid position, return error
	return 0, 0, fmt.Errorf("no valid position found for new widget")
}
