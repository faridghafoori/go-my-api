package main

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	//run database
	configs.ConnectDB()

	//run redis
	configs.InitRedis()

	//routes
	routes.UserRoutes(router)
	routes.AuthenticationRoutes(router)

	router.Run("localhost:6000")
}
