package controllers

import (
	"context"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/responses"
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
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.GeneralResponse{
					Status:  http.StatusInternalServerError,
					Message: utils.ErrorMessage,
					Data:    err.Error(),
				},
			)
			return
		}

		var episodes []models.Episode
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleEpisode models.Episode
			if err = results.Decode(&singleEpisode); err != nil {
				c.JSON(
					http.StatusInternalServerError,
					responses.GeneralResponse{
						Status:  http.StatusInternalServerError,
						Message: utils.ErrorMessage,
						Data:    err.Error(),
					},
				)
				return
			}
			episodes = append(episodes, singleEpisode)
		}

		c.JSON(
			http.StatusOK,
			responses.GeneralResponse{
				Status:  http.StatusOK,
				Message: utils.SuccessMessage,
				Data:    episodes,
			},
		)
	}
}

func CreateEpisode() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//validate the request body
		var episode models.Episode
		if err := c.BindJSON(&episode); err != nil {
			c.JSON(
				http.StatusBadRequest,
				responses.GeneralResponse{
					Status:  http.StatusBadRequest,
					Message: utils.ErrorMessage,
					Data:    err.Error(),
				},
			)
			return
		}

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

		_, err := episodeCollection.InsertOne(ctx, newEpisode)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.GeneralResponse{
					Status:  http.StatusInternalServerError,
					Message: utils.ErrorMessage,
					Data:    err.Error(),
				},
			)
			return
		}

		c.JSON(
			http.StatusCreated,
			responses.GeneralResponse{
				Status:  http.StatusCreated,
				Message: utils.SuccessMessage,
				Data:    newEpisode,
			},
		)
	}
}
