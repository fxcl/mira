package router

import (
	"mira/app"

	"github.com/gin-gonic/gin"
)

// Register routes
func Register(server *gin.Engine) {
	api := server.Group("/api")

	// Create a new app container
	container := app.NewAppContainer()

	RegisterAdminGroupApi(api, container)
}
