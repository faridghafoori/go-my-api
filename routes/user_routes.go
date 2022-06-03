package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	// All routes related to users comes here
	router.GET("/users", controllers.TokenAuthMiddleware(), controllers.GetAllUsers())
	router.POST("/users", controllers.CreateUser())
	router.GET("/users/:userId", controllers.GetUser())
	router.PUT("/users/:userId", controllers.EditUser())
	router.DELETE("/users/:userId", controllers.DeleteUser())

	// All routes related to user addresses comes here
	router.GET("/users/:userId/addresses", controllers.GetUserAddresses())
	router.POST("/users/:userId/addresses", controllers.AddNewAddressToUser())
	router.PUT("/users/:userId/addresses/:addressId", controllers.EditAddressOfUser())
	router.DELETE("/users/:userId/addresses/:addressId", controllers.DeleteAddressOfUser())
}
