package routes

import (
	"gin-mongo-api/controllers"

	"github.com/gin-gonic/gin"
)

func EpisodeRoutes(router *gin.Engine) {
	episodeController := new(controllers.EpisodeController)

	router.GET("/episodes", episodeController.GetEpisodes())
	router.GET("/episodes/:episodeId", episodeController.GetEpisode())
}

func EpisodePrivateRoutes(router *gin.Engine) {
	episodeController := new(controllers.EpisodeController)

	router.POST("/episodes", episodeController.CreateEpisode())
	router.PUT("/episodes/:episodeId", episodeController.EditEpisode())
	router.DELETE("/episodes/:episodeId", episodeController.DeleteEpisode())
}
