package routes

import (
	"gin-mongo-api/controllers"
	"net/http"

	"github.com/gin-gonic/gin"
)

func IORoutes(route *gin.Engine) {
	// All routes related to io comes here
	route.POST("/upload", controllers.TokenAuthMiddleware(), controllers.Upload())
	route.StaticFS("/file", http.Dir("public"))
}