package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func RoleRoutes(router *gin.Engine) {
	router.GET("/roles", controllers.TokenAuthMiddleware(), controllers.GetRoles())
	router.GET("/roles/:roleId", controllers.TokenAuthMiddleware(), controllers.GetRole())
	router.POST("/roles", controllers.TokenAuthMiddleware(), controllers.CreateRole())
	router.PUT("/roles/:roleId", controllers.TokenAuthMiddleware(), controllers.EditRole())
	router.DELETE("/roles/:roleId", controllers.TokenAuthMiddleware(), controllers.DeleteRole())
}
