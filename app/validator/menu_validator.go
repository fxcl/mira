package validator

import (
	"strings"

	"mira/app/dto"
	"mira/common/types/constant"
	"mira/common/utils"
	"mira/common/xerrors"
)

// CreateMenuValidator validates the request to create a menu.
func CreateMenuValidator(param dto.CreateMenuRequest) error {
	if param.MenuName == "" {
		return xerrors.ErrMenuNameEmpty
	}

	if utils.Contains([]string{constant.MENU_TYPE_DIRECTORY, constant.MENU_TYPE_MENU}, param.Path) && param.Path == "" {
		return xerrors.ErrMenuPathEmpty
	}

	if param.IsFrame == constant.MENU_YES_FRAME && !strings.HasPrefix(param.Path, "http") {
		return xerrors.ErrMenuPathHttpPrefix
	}

	return nil
}

// UpdateMenuValidator validates the request to update a menu.
func UpdateMenuValidator(param dto.UpdateMenuRequest) error {
	if param.MenuId <= 0 {
		return xerrors.ErrParam
	}

	if param.MenuName == "" {
		return xerrors.ErrMenuNameEmpty
	}

	if utils.Contains([]string{constant.MENU_TYPE_DIRECTORY, constant.MENU_TYPE_MENU}, param.Path) && param.Path == "" {
		return xerrors.ErrMenuPathEmpty
	}

	if param.IsFrame == constant.MENU_YES_FRAME && !strings.HasPrefix(param.Path, "http") {
		return xerrors.ErrMenuPathHttpPrefix
	}

	if param.MenuId == param.ParentId {
		return xerrors.ErrMenuParentSelf
	}

	return nil
}
