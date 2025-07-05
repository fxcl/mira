package validator

import (
	"errors"
	"mira/app/dto"
)

// CreateDictTypeValidator validates the request to create a dictionary type.
func CreateDictTypeValidator(param dto.CreateDictTypeRequest) error {
	if param.DictName == "" {
		return errors.New("please enter the dictionary name")
	}

	if param.DictType == "" {
		return errors.New("please enter the dictionary type")
	}

	return nil
}

// UpdateDictTypeValidator validates the request to update a dictionary type.
func UpdateDictTypeValidator(param dto.UpdateDictTypeRequest) error {
	if param.DictId <= 0 {
		return errors.New("parameter error")
	}

	if param.DictName == "" {
		return errors.New("please enter the dictionary name")
	}

	if param.DictType == "" {
		return errors.New("please enter the dictionary type")
	}

	return nil
}

// CreateDictDataValidator validates the request to create dictionary data.
func CreateDictDataValidator(param dto.CreateDictDataRequest) error {
	if param.DictLabel == "" {
		return errors.New("please enter the data label")
	}

	if param.DictValue == "" {
		return errors.New("please enter the data key value")
	}

	return nil
}

// UpdateDictDataValidator validates the request to update dictionary data.
func UpdateDictDataValidator(param dto.UpdateDictDataRequest) error {
	if param.DictCode <= 0 {
		return errors.New("parameter error")
	}

	if param.DictLabel == "" {
		return errors.New("please enter the data label")
	}

	if param.DictValue == "" {
		return errors.New("please enter the data key value")
	}

	return nil
}
