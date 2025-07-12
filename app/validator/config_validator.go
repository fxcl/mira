package validator

import (
	"mira/app/dto"
	"mira/common/xerrors"
)

// CreateConfigValidator validates the request to create a configuration.
func CreateConfigValidator(param dto.CreateConfigRequest) error {
	switch {
	case param.ConfigName == "":
		return xerrors.ErrConfigNameEmpty
	case param.ConfigKey == "":
		return xerrors.ErrConfigKeyEmpty
	case param.ConfigValue == "":
		return xerrors.ErrConfigValueEmpty
	default:
		return nil
	}
}

// UpdateConfigValidator validates the request to update a configuration.
func UpdateConfigValidator(param dto.UpdateConfigRequest) error {
	switch {
	case param.ConfigId <= 0:
		return xerrors.ErrParam
	case param.ConfigName == "":
		return xerrors.ErrConfigNameEmpty
	case param.ConfigKey == "":
		return xerrors.ErrConfigKeyEmpty
	case param.ConfigValue == "":
		return xerrors.ErrConfigValueEmpty
	default:
		return nil
	}
}
