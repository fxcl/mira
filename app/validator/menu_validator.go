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
	switch {
	case param.MenuName == "":
		return xerrors.ErrMenuNameEmpty
	case utils.Contains([]string{constant.MENU_TYPE_DIRECTORY, constant.MENU_TYPE_MENU}, param.MenuType) && param.Path == "":
		return xerrors.ErrMenuPathEmpty
	case param.MenuType == constant.MENU_TYPE_MENU && param.IsFrame == constant.MENU_YES_FRAME && !strings.HasPrefix(param.Path, "http"):
		return xerrors.ErrMenuPathHttpPrefix
	default:
		return nil
	}
}

// UpdateMenuValidator validates the request to update a menu.
func UpdateMenuValidator(param dto.UpdateMenuRequest) error {
	switch {
	case param.MenuId <= 0:
		return xerrors.ErrParam
	case param.MenuName == "":
		return xerrors.ErrMenuNameEmpty
	case utils.Contains([]string{constant.MENU_TYPE_DIRECTORY, constant.MENU_TYPE_MENU}, param.MenuType) && param.Path == "":
		return xerrors.ErrMenuPathEmpty
	case param.MenuType == constant.MENU_TYPE_MENU && param.IsFrame == constant.MENU_YES_FRAME && !strings.HasPrefix(param.Path, "http"):
		return xerrors.ErrMenuPathHttpPrefix
	case param.MenuId == param.ParentId:
		return xerrors.ErrMenuParentSelf
	default:
		return nil
	}
}
