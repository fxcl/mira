package validator

import (
	"errors"
	"mira/app/dto"
)

// CreateConfigValidator validates the request to create a configuration.
func CreateConfigValidator(param dto.CreateConfigRequest) error {
	if param.ConfigName == "" {
		return errors.New("please enter the parameter name")
	}

	if param.ConfigKey == "" {
		return errors.New("please enter the parameter key")
	}

	if param.ConfigValue == "" {
		return errors.New("please enter the parameter value")
	}

	return nil
}

// UpdateConfigValidator validates the request to update a configuration.
func UpdateConfigValidator(param dto.UpdateConfigRequest) error {
	if param.ConfigId <= 0 {
		return errors.New("parameter error")
	}

	if param.ConfigName == "" {
		return errors.New("please enter the parameter name")
	}

	if param.ConfigKey == "" {
		return errors.New("please enter the parameter key")
	}

	if param.ConfigValue == "" {
		return errors.New("please enter the parameter value")
	}

	return nil
}
