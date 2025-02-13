package app

import (
	"backend/internal/features/device/domain"
	"backend/internal/features/device/domain/command"
	resourceDomain "backend/internal/features/resource/domain"
	scopeDomain "backend/internal/features/scope/domain"
	"backend/shared/base"
	baseCmd "backend/shared/base/command"
	"backend/shared/validator"
	"context"
	"fmt"
)

type DeviceService interface {
	CreateDevice(ctx context.Context, input *command.CreateDeviceInput) (*domain.Device, error)
	GetDevice(ctx context.Context, input *command.DeviceIDInput) (*domain.Device, error)
	ListDevices(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Device, error)
	UpdateDevice(ctx context.Context, input *command.UpdateDeviceInput) (*domain.Device, error)
	DeleteDevice(ctx context.Context, input *command.DeviceIDInput) error
}
type deviceService struct {
	repo domain.DeviceRepository
	*base.BaseService
}

func NewService(base *base.BaseService, repo domain.DeviceRepository) DeviceService {
	return &deviceService{
		repo:        repo,
		BaseService: base,
	}
}
func (s *deviceService) CreateDevice(ctx context.Context, input *command.CreateDeviceInput) (*domain.Device, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Device,
		Scope:     scopeDomain.Create,
	})
	if err != nil {
		return nil, err
	}
	err = validator.Validate(input)
	if err != nil {
		return nil, fmt.Errorf("validation error %s", err.Error())
	}
	device, err := domain.New(input.TenantDomain, input.BranchName, input.Name, input.Parent, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	err = s.repo.Create(ctx, device)
	if err != nil {
		return nil, fmt.Errorf("failed to create device: %w", err)
	}

	return device, nil
}

func (s *deviceService) GetDevice(ctx context.Context, input *command.DeviceIDInput) (*domain.Device, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Device,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}
	device, err := s.repo.FindByID(ctx, &input.BaseInput, input.DeviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	return device, nil
}

func (s *deviceService) ListDevices(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Device, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: *input,
		Resource:  resourceDomain.Device,
		Scope:     scopeDomain.View,
	})
	if err != nil {
		return nil, err
	}
	devices, err := s.repo.FindAll(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}
	return devices, nil
}

func (s *deviceService) UpdateDevice(ctx context.Context, input *command.UpdateDeviceInput) (*domain.Device, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Device,
		Scope:     scopeDomain.Update,
	})
	if err != nil {
		return nil, err
	}
	device, err := s.repo.FindByID(ctx, &input.BaseInput, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	err = device.Update(input.Name, input.Parent, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	if err := s.repo.Update(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	return device, nil
}
func (s *deviceService) DeleteDevice(ctx context.Context, input *command.DeviceIDInput) error {

	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Device,
		Scope:     scopeDomain.Delete,
	})
	if err != nil {
		return err
	}
	err = s.repo.Remove(ctx, &input.BaseInput, input.DeviceID)
	if err != nil {
		return err
	}
	return nil
}
