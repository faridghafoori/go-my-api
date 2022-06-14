package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func UserRoutes(router *gin.Engine) {
	// All routes related to users comes here
	userController := new(controllers.UserController)

	router.GET("/users", userController.GetUsers())
	router.POST("/users", userController.CreateUser())
	router.GET("/users/:userId", userController.GetUser())
	router.PUT("/users/:userId", userController.EditUser())
	router.DELETE("/users/:userId", userController.DeleteUser())

	// All routes related to user addresses comes here
	addressController := new(controllers.AddressController)
	router.GET("/users/:userId/addresses", addressController.GetUserAddresses())
	router.POST("/users/:userId/addresses", addressController.AddNewAddressToUser())
	router.PUT("/users/:userId/addresses/:addressId", addressController.EditAddressOfUser())
	router.DELETE("/users/:userId/addresses/:addressId", addressController.DeleteAddressOfUser())
}
