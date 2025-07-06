package validator

import (
	"mira/app/dto"
	"mira/common/xerrors"
)

// CreatePostValidator validates the request to create a post.
func CreatePostValidator(param dto.CreatePostRequest) error {
	if param.PostCode == "" {
		return xerrors.ErrPostCodeEmpty
	}

	if param.PostName == "" {
		return xerrors.ErrPostNameEmpty
	}

	return nil
}

// UpdatePostValidator validates the request to update a post.
func UpdatePostValidator(param dto.UpdatePostRequest) error {
	if param.PostId <= 0 {
		return xerrors.ErrParam
	}

	if param.PostCode == "" {
		return xerrors.ErrPostCodeEmpty
	}

	if param.PostName == "" {
		return xerrors.ErrPostNameEmpty
	}

	return nil
}
