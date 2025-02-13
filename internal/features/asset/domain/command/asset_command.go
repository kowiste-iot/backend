package command

import "backend/shared/base/command"

type CreateAssetInput struct {
	command.BaseInput
	ID          string `validate:"omitempty,uuidv7"`
	Name        string `validate:"required,min=3,max=255"`
	Parent      string `validate:"omitempty,uuidv7"`
	Description string `validate:"omitempty,min=3,max=512"`
}

type UpdateAssetInput struct {
	command.BaseInput
	ID          string `validate:"required,uuidv7"`
	Name        string `validate:"required,min=3,max=255"`
	Parent      string `validate:"omitempty,uuidv7"`
	Description string `validate:"omitempty,min=3,max=512"`
}

type AssetIDInput struct {
	command.BaseInput
	AssetID string `validate:"required,uuidv7"`
}
