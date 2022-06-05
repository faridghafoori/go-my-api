package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func EpisodeRoutes(router *gin.Engine) {
	// All routes related to episodes comes here
	router.GET("/episodes", controllers.GetEpisodes())
	router.GET("/episodes/:episodeId", controllers.GetEpisode())
	router.POST("/episodes", controllers.TokenAuthMiddleware(), controllers.CreateEpisode())

}
