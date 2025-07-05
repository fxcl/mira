package dto

// Save Menu
type SaveMenu struct {
	MenuId    int    `json:"menuId"`
	MenuName  string `json:"menuName"`
	ParentId  int    `json:"parentId"`
	OrderNum  int    `json:"orderNum"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Query     string `json:"query"`
	RouteName string `json:"routeName"`
	IsFrame   *int   `json:"isFrame"`
	IsCache   *int   `json:"isCache"`
	MenuType  string `json:"menuType"`
	Visible   string `json:"visible"`
	Perms     string `json:"perms"`
	Icon      string `json:"icon"`
	Status    string `json:"status"`
	CreateBy  string `json:"createBy"`
	UpdateBy  string `json:"updateBy"`
	Remark    string `json:"remark"`
}

// Menu List
type MenuListRequest struct {
	MenuName string `query:"menuName" form:"menuName"`
	Status   string `query:"status" form:"status"`
}

// Create Menu
type CreateMenuRequest struct {
	MenuName  string `json:"menuName"`
	ParentId  int    `json:"parentId"`
	OrderNum  int    `json:"orderNum"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Query     string `json:"query"`
	RouteName string `json:"routeName"`
	IsFrame   int    `json:"isFrame,string"`
	IsCache   int    `json:"isCache,string"`
	MenuType  string `json:"menuType"`
	Visible   string `json:"visible"`
	Perms     string `json:"perms"`
	Icon      string `json:"icon"`
	Status    string `json:"status"`
}

// Update Menu
type UpdateMenuRequest struct {
	MenuId    int    `json:"menuId"`
	MenuName  string `json:"menuName"`
	ParentId  int    `json:"parentId"`
	OrderNum  int    `json:"orderNum"`
	Path      string `json:"path"`
	Component string `json:"component"`
	Query     string `json:"query"`
	RouteName string `json:"routeName"`
	IsFrame   int    `json:"isFrame"`
	IsCache   int    `json:"isCache"`
	MenuType  string `json:"menuType"`
	Visible   string `json:"visible"`
	Perms     string `json:"perms"`
	Icon      string `json:"icon"`
	Status    string `json:"status"`
}
