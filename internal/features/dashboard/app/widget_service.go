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
	ListWidgets(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Widget, error)
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

	// Convert link data to domain WidgetLinkData
	linkData := make([]domain.WidgetLinkData, len(input.Link))
	for i, link := range input.Link {
		linkData[i] = domain.NewWidgetLinkData(link.Measure, link.Tag, link.Legend)
	}

	widget, err = domain.NewWidget(
		input.TenantDomain,
		input.BranchName,
		input.DashboardID,
		input.TypeWidget,
		input.I,
		input.X,
		input.Y,
		input.W,
		input.H,
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

	widget, err = s.repo.FindByID(ctx, &input.BaseInput, input.DashboardID, input.WidgetID)
	if err != nil {
		return nil, fmt.Errorf("failed to get widget: %w", err)
	}

	return widget, nil
}

func (s *widgetService) ListWidgets(ctx context.Context, input *baseCmd.BaseInput) (widgets []*domain.Widget, err error) {
	err = s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: *input,
		Resource:  resourceDomain.Widget,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}

	widgets, err = s.repo.FindAll(ctx, input)
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

	widget, err = s.repo.FindByID(ctx, &input.BaseInput, input.DashboardID, input.ID)
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
		input.I,
		input.X,
		input.Y,
		input.W,
		input.H,
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

	widget, err = s.repo.FindByID(ctx, &input.BaseInput, input.DashboardID, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get widget: %w", err)
	}

	err = widget.UpdatePosition(
		input.I,
		input.X,
		input.Y,
		input.W,
		input.H,
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

	err = s.repo.Remove(ctx, &input.BaseInput, input.DashboardID, input.WidgetID)
	if err != nil {
		return err
	}

	return nil
}
