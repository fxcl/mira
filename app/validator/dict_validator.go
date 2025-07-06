package validator

import (
	"mira/app/dto"
	"mira/common/xerrors"
)

// CreateDictTypeValidator validates the request to create a dictionary type.
func CreateDictTypeValidator(param dto.CreateDictTypeRequest) error {
	if param.DictName == "" {
		return xerrors.ErrDictNameEmpty
	}

	if param.DictType == "" {
		return xerrors.ErrDictTypeEmpty
	}

	return nil
}

// UpdateDictTypeValidator validates the request to update a dictionary type.
func UpdateDictTypeValidator(param dto.UpdateDictTypeRequest) error {
	if param.DictId <= 0 {
		return xerrors.ErrParam
	}

	if param.DictName == "" {
		return xerrors.ErrDictNameEmpty
	}

	if param.DictType == "" {
		return xerrors.ErrDictTypeEmpty
	}

	return nil
}

// CreateDictDataValidator validates the request to create dictionary data.
func CreateDictDataValidator(param dto.CreateDictDataRequest) error {
	if param.DictLabel == "" {
		return xerrors.ErrDictLabelEmpty
	}

	if param.DictValue == "" {
		return xerrors.ErrDictValueEmpty
	}

	return nil
}

// UpdateDictDataValidator validates the request to update dictionary data.
func UpdateDictDataValidator(param dto.UpdateDictDataRequest) error {
	if param.DictCode <= 0 {
		return xerrors.ErrParam
	}

	if param.DictLabel == "" {
		return xerrors.ErrDictLabelEmpty
	}

	if param.DictValue == "" {
		return xerrors.ErrDictValueEmpty
	}

	return nil
}
