package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func AuthenticationRoutes(router *gin.Engine) {
	router.POST("/authenticate", controllers.Authenticate())
	router.POST("/register", controllers.Register())
	router.POST("/token/refresh", controllers.Refresh())
	router.POST("/totp", controllers.TokenAuthMiddleware("totp"), controllers.TOTPGenerator())
	router.POST("/verify", controllers.TokenAuthMiddleware("totp"), controllers.VerifyTOTP())
	router.POST("/logout", controllers.TokenAuthMiddleware(), controllers.Logout())
}
