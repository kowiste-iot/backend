package command

import "backend/shared/base/command"

type DependencyChangeInput struct {
	command.BaseInput
	PreviousAssetID string
	Feature         string           // type of feature (measure, device, etc)
	Action          DependencyAction // create, update, delete
	FeatureID       string
	NewAssetID      string
}

type DependencyAction string

const (
	DependencyActionCreate DependencyAction = "create"
	DependencyActionUpdate DependencyAction = "update"
	DependencyActionDelete DependencyAction = "delete"
)
