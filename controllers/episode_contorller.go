package controllers

import (
	"context"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var episodeCollection *mongo.Collection = configs.GetCollection(configs.DB, "episodes")

func GetEpisodes() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		results, err := episodeCollection.Find(ctx, bson.M{})
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		var episodes []models.Episode
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleEpisode models.Episode
			err = results.Decode(&singleEpisode)
			utils.GenerateErrorOutput(http.StatusBadRequest, err, c)
			episodes = append(episodes, singleEpisode)
		}

		utils.GenerateSuccessOutput(episodes, c)
	}
}

func CreateEpisode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//validate the request body
		var episode models.Episode
		err := c.BindJSON(&episode)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		//use the validator library to validate required fields
		utils.ValidateStruct(&episode)

		newEpisode := models.Episode{
			Id:          primitive.NewObjectID(),
			Title:       episode.Title,
			Description: episode.Description,
			Image:       episode.Image,
			Background:  episode.Background,
			Rate:        episode.Rate,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		_, err = episodeCollection.InsertOne(ctx, newEpisode)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		utils.GenerateSuccessOutput(newEpisode, c)
	}
}

func GetEpisode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		episodeId := c.Param("episodeId")
		defer cancel()

		episodeObjId, _ := primitive.ObjectIDFromHex(episodeId)
		var episode models.Episode
		err := episodeCollection.FindOne(ctx, bson.M{"id": episodeObjId}).Decode(&episode)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		utils.GenerateSuccessOutput(episode, c)
	}
}
