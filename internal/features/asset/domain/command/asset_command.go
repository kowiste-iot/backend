package command

import "ddd/shared/base/command"

type CreateAssetInput struct {
	command.BaseInput
	ID          string
	Name        string
	Parent      string
	Description string
}

type UpdateAssetInput struct {
	command.BaseInput
	ID          string
	Name        string
	Parent      string
	Description string
}


type AssetIDInput struct {
	command.BaseInput
	AssetID string
}
