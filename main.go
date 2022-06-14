package main

import (
	"gin-mongo-api/configs"
	"gin-mongo-api/server"
)

func main() {
	//run database
	configs.ConnectDB()

	//run redis
	configs.InitRedis()

	//run minio server
	configs.InitMinio()

	//run server
	server.Init()
}
