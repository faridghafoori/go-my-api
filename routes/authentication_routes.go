package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthenticationRoutes(router *gin.Engine) {
	router.POST("/authenticate", controllers.Authenticate())
	router.POST("/logout", controllers.Logout())
	router.POST("/token/refresh", controllers.Refresh())
}
