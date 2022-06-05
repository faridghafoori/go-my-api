package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	// All routes related to users comes here
	router.GET("/users", controllers.TokenAuthMiddleware(), controllers.GetUsers())
	router.POST("/users", controllers.TokenAuthMiddleware(), controllers.CreateUser())
	router.GET("/users/:userId", controllers.TokenAuthMiddleware(), controllers.GetUser())
	router.PUT("/users/:userId", controllers.TokenAuthMiddleware(), controllers.EditUser())
	router.DELETE("/users/:userId", controllers.TokenAuthMiddleware(), controllers.DeleteUser())

	// All routes related to user addresses comes here
	router.GET("/users/:userId/addresses", controllers.TokenAuthMiddleware(), controllers.GetUserAddresses())
	router.POST("/users/:userId/addresses", controllers.TokenAuthMiddleware(), controllers.AddNewAddressToUser())
	router.PUT("/users/:userId/addresses/:addressId", controllers.TokenAuthMiddleware(), controllers.EditAddressOfUser())
	router.DELETE("/users/:userId/addresses/:addressId", controllers.TokenAuthMiddleware(), controllers.DeleteAddressOfUser())
}
