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

	//run minio server
	configs.InitMinio()

	//routes
	routes.AuthenticationRoutes(router)
	routes.UserRoutes(router)
	routes.RoleRoutes(router)
	routes.EpisodeRoutes(router)
	routes.IORoutes(router)
	router.Run(configs.ENV_LAUNCH_PROJECT_URI())
}
