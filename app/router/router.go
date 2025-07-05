package router

import (
	"github.com/gin-gonic/gin"
)

// Register routes
func Register(server *gin.Engine) {
	api := server.Group("/api")

	RegisterAdminGroupApi(api)
}
