package validator

import (
	"mira/app/dto"
	"mira/common/xerrors"
)

// CreateDictTypeValidator validates the request to create a dictionary type.
func CreateDictTypeValidator(param dto.CreateDictTypeRequest) error {
	switch {
	case param.DictName == "":
		return xerrors.ErrDictNameEmpty
	case param.DictType == "":
		return xerrors.ErrDictTypeEmpty
	default:
		return nil
	}
}

// UpdateDictTypeValidator validates the request to update a dictionary type.
func UpdateDictTypeValidator(param dto.UpdateDictTypeRequest) error {
	switch {
	case param.DictId <= 0:
		return xerrors.ErrParam
	case param.DictName == "":
		return xerrors.ErrDictNameEmpty
	case param.DictType == "":
		return xerrors.ErrDictTypeEmpty
	default:
		return nil
	}
}

// CreateDictDataValidator validates the request to create dictionary data.
func CreateDictDataValidator(param dto.CreateDictDataRequest) error {
	switch {
	case param.DictLabel == "":
		return xerrors.ErrDictLabelEmpty
	case param.DictValue == "":
		return xerrors.ErrDictValueEmpty
	default:
		return nil
	}
}

// UpdateDictDataValidator validates the request to update dictionary data.
func UpdateDictDataValidator(param dto.UpdateDictDataRequest) error {
	switch {
	case param.DictCode <= 0:
		return xerrors.ErrParam
	case param.DictLabel == "":
		return xerrors.ErrDictLabelEmpty
	case param.DictValue == "":
		return xerrors.ErrDictValueEmpty
	default:
		return nil
	}
}
