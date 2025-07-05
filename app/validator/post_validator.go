package validator

import (
	"errors"
	"mira/app/dto"
)

// CreatePostValidator validates the request to create a post.
func CreatePostValidator(param dto.CreatePostRequest) error {
	if param.PostCode == "" {
		return errors.New("please enter the post code")
	}

	if param.PostName == "" {
		return errors.New("please enter the post name")
	}

	return nil
}

// UpdatePostValidator validates the request to update a post.
func UpdatePostValidator(param dto.UpdatePostRequest) error {
	if param.PostId <= 0 {
		return errors.New("parameter error")
	}

	if param.PostCode == "" {
		return errors.New("please enter the post code")
	}

	if param.PostName == "" {
		return errors.New("please enter the post name")
	}

	return nil
}
