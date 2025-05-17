// internal/core/services/container.go
package services

import (
	appAction "backend/internal/features/action/app"
	appAlert "backend/internal/features/alert/app"
	appAsset "backend/internal/features/asset/app"
	appDashboard "backend/internal/features/dashboard/app"
	appDataStore "backend/internal/features/datastore/app"
	appDevice "backend/internal/features/device/app"
	appIngest "backend/internal/features/ingest/app"
	appMeasure "backend/internal/features/measure/app"
	appPermission "backend/internal/features/permission/app"
	appResource "backend/internal/features/resource/app"
	appRole "backend/internal/features/user/app"
	appScope "backend/internal/features/scope/app"
	appTenant "backend/internal/features/tenant/app"
	appUser "backend/internal/features/user/app"

	// domainStream "backend/shared/stream/domain"
	appToken "backend/shared/token/app"
	appWS "backend/shared/websocket/app"

	"backend/pkg/config"
	kcCore "backend/shared/keycloak"

	"backend/shared/base"
	"fmt"
)

type Services struct {
	// StreamService     domainStream.StreamClient
	AssetDepService   appAsset.AssetDependencyService
	AssetService      appAsset.AssetService
	MeasureService    appMeasure.MeasureService
	DashboardService  appDashboard.DashboardService
	WidgetService     appDashboard.WidgetService
	DeviceService     appDevice.DeviceService
	ActionService     appAction.ActionService
	AlertService      appAlert.AlertService
	UserService       appUser.UserService
	ScopeService      appScope.ScopeService
	RoleService       appRole.RoleService
	PermissionService appPermission.PermissionService
	ResourceService   appResource.ResourceService
	BranchService     appTenant.BranchService
	TenantService     appTenant.TenantService
	IngestService     appIngest.IngestService
	DataStoreService  appDataStore.DataStoreService
	TokenService      *appToken.TokenService
	WebSocketHub      *appWS.Hub
	WebSocketService  *appWS.WebSocketStreamService
}

type Container struct {
	base         *base.BaseService
	auth         *kcCore.Keycloak
	config       *config.Config
	tenantConfig *config.TenantConfiguration
}

func NewContainer(
	base *base.BaseService,
	auth *kcCore.Keycloak,
	tenantConfig *config.TenantConfiguration,
	config *config.Config) *Container {
	return &Container{
		base:         base,
		auth:         auth,
		config:       config,
		tenantConfig: tenantConfig,
	}
}

func (c *Container) Initialize() (*Services, error) {
	services := &Services{}

	// if err := c.initializeStreamService(services); err != nil {
	// 	return nil, fmt.Errorf("failed to initialize strem services: %w", err)
	// }
	if err := c.initializeAssetServices(services); err != nil {
		return nil, fmt.Errorf("failed to initialize asset services: %w", err)
	}
	if err := c.initializeMeasureService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize measure service: %w", err)
	}
	if err := c.initializeDashboardService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize dashboard service: %w", err)
	}
	if err := c.initializeDeviceService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize device service: %w", err)
	}
	if err := c.initializeActionService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize action service: %w", err)
	}
	if err := c.initializeAlertService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize alert service: %w", err)
	}
	if err := c.initializeUserService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize user service: %w", err)
	}

	if err := c.initializeRoleService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize role service: %w", err)
	}

	if err := c.initializeScopeService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize scope service: %w", err)
	}

	if err := c.initializePermissionService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize permission service: %w", err)
	}

	if err := c.initializeResourceService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize resource service: %w", err)
	}

	if err := c.initializeBranchService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize branch service: %w", err)
	}

	if err := c.initializeTenantService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize tenant service: %w", err)
	}
	if err := c.initializeIngestService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize ingest service: %w", err)
	}
	if err := c.initializeStoreService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize store service: %w", err)
	}
	if err := c.initializeTokenService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize token service: %w", err)
	}
	if err := c.initializeWebsocketService(services); err != nil {
		return nil, fmt.Errorf("failed to initialize websocket service: %w", err)
	}
	return services, nil
}
