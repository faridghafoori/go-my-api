package controllers

import (
	"context"
	"errors"
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

func CreateEpisode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		// validate the request body
		var episode models.Episode
		err := c.BindJSON(&episode)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		//use the validator library to validate required fields
		validationErr := validate.Struct(&episode)
		utils.GenerateErrorOutput(http.StatusBadRequest, validationErr, c)

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

func EditEpisode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		episodeId := c.Param("episodeId")
		defer cancel()

		episodeObjId, _ := primitive.ObjectIDFromHex(episodeId)

		var episode models.Episode
		err := episodeCollection.FindOne(ctx, bson.M{"id": episodeObjId}).Decode(&episode)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		err = c.BindJSON(&episode)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		validationErr := validate.Struct(&episode)
		utils.GenerateErrorOutput(http.StatusBadRequest, validationErr, c)

		episode.UpdatedAt = time.Now()

		_, err = episodeCollection.UpdateOne(ctx, bson.M{"id": episodeObjId}, bson.M{"$set": episode})
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		utils.GenerateSuccessOutput(episode, c)
	}
}

func DeleteEpisode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		episodeId := c.Param("episodeId")
		defer cancel()

		episodeObjId, _ := primitive.ObjectIDFromHex(episodeId)
		result, err := episodeCollection.DeleteOne(ctx, bson.M{"id": episodeObjId})
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		if result.DeletedCount < 1 {
			utils.GenerateErrorOutput(
				http.StatusNotFound,
				errors.New("role not found"),
				c,
				map[string]interface{}{
					"data": "Episode with specified ID not found!",
				},
			)
		}

		utils.GenerateSuccessOutput("Episode successfully deleted!", c)
	}
}
