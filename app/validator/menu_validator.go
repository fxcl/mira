package validator

import (
	"errors"
	"mira/app/dto"
	"mira/common/types/constant"
	"mira/common/utils"
	"strings"
)

// CreateMenuValidator validates the request to create a menu.
func CreateMenuValidator(param dto.CreateMenuRequest) error {
	if param.MenuName == "" {
		return errors.New("please enter the menu name")
	}

	if utils.Contains([]string{constant.MENU_TYPE_DIRECTORY, constant.MENU_TYPE_MENU}, param.Path) && param.Path == "" {
		return errors.New("please enter the route address")
	}

	if param.IsFrame == constant.MENU_YES_FRAME && !strings.HasPrefix(param.Path, "http") {
		return errors.New("failed to add menu " + param.MenuName + ", the address must start with http(s)://")
	}

	return nil
}

// UpdateMenuValidator validates the request to update a menu.
func UpdateMenuValidator(param dto.UpdateMenuRequest) error {
	if param.MenuId <= 0 {
		return errors.New("parameter error")
	}

	if param.MenuName == "" {
		return errors.New("please enter the menu name")
	}

	if utils.Contains([]string{constant.MENU_TYPE_DIRECTORY, constant.MENU_TYPE_MENU}, param.Path) && param.Path == "" {
		return errors.New("please enter the route address")
	}

	if param.IsFrame == constant.MENU_YES_FRAME && !strings.HasPrefix(param.Path, "http") {
		return errors.New("failed to modify menu " + param.MenuName + ", the address must start with http(s)://")
	}

	if param.MenuId == param.ParentId {
		return errors.New("failed to modify menu " + param.MenuName + ", the parent menu cannot be itself")
	}

	return nil
}
