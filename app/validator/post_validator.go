package validator

import (
	"mira/app/dto"
	"mira/common/xerrors"
)

// CreatePostValidator validates the request to create a post.
func CreatePostValidator(param dto.CreatePostRequest) error {
	switch {
	case param.PostCode == "":
		return xerrors.ErrPostCodeEmpty
	case param.PostName == "":
		return xerrors.ErrPostNameEmpty
	default:
		return nil
	}
}

// UpdatePostValidator validates the request to update a post.
func UpdatePostValidator(param dto.UpdatePostRequest) error {
	switch {
	case param.PostId <= 0:
		return xerrors.ErrParam
	case param.PostCode == "":
		return xerrors.ErrPostCodeEmpty
	case param.PostName == "":
		return xerrors.ErrPostNameEmpty
	default:
		return nil
	}
}
