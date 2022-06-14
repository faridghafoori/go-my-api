package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func RoleRoutes(router *gin.Engine) {
	router.GET("/roles", controllers.GetRoles())
	router.GET("/roles/:roleId", controllers.GetRole())
	router.POST("/roles", controllers.CreateRole())
	router.PUT("/roles/:roleId", controllers.EditRole())
	router.DELETE("/roles/:roleId", controllers.DeleteRole())
}
