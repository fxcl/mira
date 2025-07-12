package app

import (
	monitorcontroller "mira/app/controller/monitor"
	systemcontroller "mira/app/controller/system"
	"mira/app/middleware"
	"mira/app/security"
	"mira/app/service"

	"github.com/gin-gonic/gin"
)

// AppContainer holds all instances of services, controllers, and middlewares.
type AppContainer struct {
	// Services
	LogininforService *service.LogininforService
	OperLogService    *service.OperLogService
	UserService       *service.UserService
	DeptService       *service.DeptService
	RoleService       *service.RoleService
	PostService       *service.PostService
	MenuService       *service.MenuService
	ConfigService     *service.ConfigService
	DictTypeService   *service.DictTypeService
	DictDataService   *service.DictDataService

	// Security
	Security *security.Security

	// Controllers
	LogininforController *monitorcontroller.LogininforController
	OperlogController    *monitorcontroller.OperlogController
	UserController       *systemcontroller.UserController
	RoleController       *systemcontroller.RoleController
	MenuController       *systemcontroller.MenuController
	DeptController       *systemcontroller.DeptController
	PostController       *systemcontroller.PostController
	DictTypeController   *systemcontroller.DictTypeController
	DictDataController   *systemcontroller.DictDataController
	ConfigController     *systemcontroller.ConfigController
}

// NewAppContainer creates and initializes a new AppContainer.
func NewAppContainer() *AppContainer {
	// Instantiate services
	logininforService := &service.LogininforService{}
	operLogService := &service.OperLogService{}
	userService := &service.UserService{}
	deptService := &service.DeptService{}
	roleService := &service.RoleService{}
	postService := &service.PostService{}
	menuService := &service.MenuService{}
	configService := &service.ConfigService{}
	dictTypeService := &service.DictTypeService{}
	dictDataService := &service.DictDataService{}

	// Instantiate security
	sec := security.NewSecurity(userService)

	// Instantiate controllers with dependencies
	logininforController := monitorcontroller.NewLogininforController(logininforService)
	operlogController := monitorcontroller.NewOperlogController(operLogService)
	userController := systemcontroller.NewUserController(userService, deptService, roleService, postService, configService)
	roleController := systemcontroller.NewRoleController(roleService, deptService, userService)
	menuController := systemcontroller.NewMenuController(menuService)
	deptController := systemcontroller.NewDeptController(deptService, userService)
	postController := systemcontroller.NewPostController(postService)
	dictTypeController := systemcontroller.NewDictTypeController(dictTypeService)
	dictDataController := systemcontroller.NewDictDataController(dictDataService)
	configController := systemcontroller.NewConfigController(configService)

	return &AppContainer{
		LogininforService:    logininforService,
		OperLogService:       operLogService,
		UserService:          userService,
		DeptService:          deptService,
		RoleService:          roleService,
		PostService:          postService,
		MenuService:          menuService,
		ConfigService:        configService,
		DictTypeService:      dictTypeService,
		DictDataService:      dictDataService,
		Security:             sec,
		LogininforController: logininforController,
		OperlogController:    operlogController,
		UserController:       userController,
		RoleController:       roleController,
		MenuController:       menuController,
		DeptController:       deptController,
		PostController:       postController,
		DictTypeController:   dictTypeController,
		DictDataController:   dictDataController,
		ConfigController:     configController,
	}
}

// LogininforMiddleware returns the login info middleware with its dependencies.
func (ac *AppContainer) LogininforMiddleware() gin.HandlerFunc {
	return middleware.LogininforMiddleware(ac.LogininforService)
}

// OperLogMiddleware returns the operation log middleware with its dependencies.
func (ac *AppContainer) OperLogMiddleware(title string, businessType int) gin.HandlerFunc {
	return middleware.OperLogMiddleware(ac.OperLogService, title, businessType, security.GetAuthUser)
}

// HasPerm returns the permission check middleware.
func (ac *AppContainer) HasPerm(perm string) gin.HandlerFunc {
	return middleware.HasPerm(ac.Security, perm)
}
