package validator

import (
	"mira/app/dto"
	"mira/common/xerrors"
)

// CreateConfigValidator validates the request to create a configuration.
func CreateConfigValidator(param dto.CreateConfigRequest) error {
	if param.ConfigName == "" {
		return xerrors.ErrConfigNameEmpty
	}

	if param.ConfigKey == "" {
		return xerrors.ErrConfigKeyEmpty
	}

	if param.ConfigValue == "" {
		return xerrors.ErrConfigValueEmpty
	}

	return nil
}

// UpdateConfigValidator validates the request to update a configuration.
func UpdateConfigValidator(param dto.UpdateConfigRequest) error {
	if param.ConfigId <= 0 {
		return xerrors.ErrParam
	}

	if param.ConfigName == "" {
		return xerrors.ErrConfigNameEmpty
	}

	if param.ConfigKey == "" {
		return xerrors.ErrConfigKeyEmpty
	}

	if param.ConfigValue == "" {
		return xerrors.ErrConfigValueEmpty
	}

	return nil
}
