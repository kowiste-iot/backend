package app

import (
	"context"
	"ddd/shared/auth/domain/command"
	"ddd/shared/auth/domain/resource"
	baseCmd "ddd/shared/base/command"

	"fmt"
)

func (s *Service) GetResources(ctx context.Context, input *baseCmd.BaseInput) (resources []resource.Resource, err error) {

	client, err := s.clientProvider.GetClientByClientID(ctx, input.TenantDomain, command.ClientName(input.BranchName))
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}
	var temp resource.Resources
	temp, err = s.resourceProvider.ListResources(ctx, input.TenantDomain, *client.ID)
	if err != nil {
		return
	}

	resources = temp.FilterResource()
	return
}
