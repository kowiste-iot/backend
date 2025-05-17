package app

import (
	"backend/internal/features/device/domain"
	"backend/internal/features/device/domain/command"

	appAsset "backend/internal/features/asset/app"
	assetCmd "backend/internal/features/asset/domain/command"

	resourceDomain "backend/internal/features/resource/domain"
	scopeDomain "backend/internal/features/scope/domain"

	"backend/shared/base"
	baseCmd "backend/shared/base/command"
	"backend/shared/validator"
	"context"
	"fmt"
)

type DeviceService interface {
	CreateDevice(ctx context.Context, input *command.CreateDeviceInput) (*domain.Device, string, error)
	GetDevice(ctx context.Context, input *command.DeviceIDInput) (*domain.Device, error)
	ListDevices(ctx context.Context, input *baseCmd.BaseInput) ([]*domain.Device, error)
	UpdateDevice(ctx context.Context, input *command.UpdateDeviceInput) (*domain.Device, error)
	DeleteDevice(ctx context.Context, input *command.DeviceIDInput) error
}
type deviceService struct {
	repo     domain.DeviceRepository
	broker   domain.BrokerProvider
	assetDep appAsset.AssetDependencyService
	*base.BaseService
}

const (
	featureName string = "device" //TODO:move to domian
)

type ServiceDependencies struct {
	Repo     domain.DeviceRepository
	Broker   domain.BrokerProvider
	AssetDep appAsset.AssetDependencyService
}

func NewService(base *base.BaseService, dep *ServiceDependencies) DeviceService {
	return &deviceService{
		repo:        dep.Repo,
		assetDep:    dep.AssetDep,
		broker:      dep.Broker,
		BaseService: base,
	}
}

func (s *deviceService) CreateDevice(ctx context.Context, input *command.CreateDeviceInput) (*domain.Device, string, error) {
	err := s.CheckPermission(ctx, &baseCmd.CheckPermissionInput{
		BaseInput: input.BaseInput,
		Resource:  resourceDomain.Device,
		Scope:     scopeDomain.Create,
	})
	if err != nil {
		return nil, "", err
	}
	err = validator.Validate(input)
	if err != nil {
		return nil, "", fmt.Errorf("validation error %s", err.Error())
	}
	device, err := domain.New(input.TenantDomain, input.BranchName, input.Name, input.Parent, input.Description)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create device: %w", err)
	}

	password, err := s.repo.Create(ctx, device)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create device: %w", err)
	}

	//Update asset parent dependecy
	err = s.assetDep.UpdateDependency(ctx, &assetCmd.DependencyChangeInput{
		BaseInput:  input.BaseInput,
		Feature:    featureName,
		Action:     assetCmd.DependencyActionCreate,
		FeatureID:  device.ID(),
		NewAssetID: device.Parent(),
	})
	if err != nil {
		return nil, "", err
	}

	return device, password, nil
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
	device, err := s.repo.FindByID(ctx, input.DeviceID)
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
	devices, err := s.repo.FindAll(ctx)
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
	device, err := s.repo.FindByID(ctx, input.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	oldParent := device.Parent()

	err = device.Update(input.Name, input.Parent, input.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}
	if err := s.repo.Update(ctx, device); err != nil {
		return nil, fmt.Errorf("failed to update device: %w", err)
	}

	if oldParent != device.Parent() {
		//Update asset parent dependecy
		err = s.assetDep.UpdateDependency(ctx, &assetCmd.DependencyChangeInput{
			BaseInput:       input.BaseInput,
			PreviousAssetID: oldParent,
			Feature:         featureName,
			Action:          assetCmd.DependencyActionUpdate,
			FeatureID:       device.ID(),
			NewAssetID:      device.Parent(),
		})
		if err != nil {
			return nil, err
		}
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
	err = s.repo.Remove(ctx, input.DeviceID)
	if err != nil {
		return err
	}

	err = s.assetDep.UpdateDependency(ctx, &assetCmd.DependencyChangeInput{
		BaseInput: input.BaseInput,
		Feature:   featureName,
		Action:    assetCmd.DependencyActionDelete,
		FeatureID: input.DeviceID,
	})
	if err != nil {
		return err
	}
	err = s.broker.KickOut(ctx, input.DeviceID)
	if err != nil {
		return err
	}
	return nil
}
