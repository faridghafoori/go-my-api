package routes

import (
	"gin-mongo-api/controllers"
	"gin-mongo-api/middleware"

	"github.com/gin-gonic/gin"
)

func AuthenticationRoutes(router *gin.Engine) {
	authController := new(controllers.AuthenticateController)

	router.POST("/authenticate", authController.Authenticate())
	router.POST("/register", authController.Register())
	router.POST("/token/refresh", authController.Refresh())
	router.POST("/totp", middleware.TokenAuthMiddleware("totp"), authController.TOTPGenerator())
	router.POST("/verify", middleware.TokenAuthMiddleware("totp"), authController.VerifyTOTP())
	router.POST("/logout", middleware.TokenAuthMiddleware(), authController.Logout())
}
