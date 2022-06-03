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

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//validate the request body
		var user models.User
		if err := c.BindJSON(&user); err != nil {
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
		utils.ValidateStruct(&user)

		//check if the username is already taken
		if checkDuplicateUser(user.Username) {
			c.JSON(
				http.StatusBadRequest,
				responses.GeneralResponse{
					Status:  http.StatusBadRequest,
					Message: utils.ErrorMessage,
					Data:    "username already taken",
				},
			)
			return
		}

		newUser := models.User{
			Id:        primitive.NewObjectID(),
			Name:      user.Name,
			Username:  user.Username,
			Password:  utils.GetMD5Hash(user.Password),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err := userCollection.InsertOne(ctx, newUser)
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
				Data:    newUser,
			},
		)
	}
}

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
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
			http.StatusOK,
			responses.GeneralResponse{
				Status:  http.StatusOK,
				Message: utils.SuccessMessage,
				Data:    user,
			},
		)
	}
}

func EditUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)
		//validate the request body
		var user models.User
		if err := c.BindJSON(&user); err != nil {
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
		utils.ValidateStruct(&user)

		//check if the username is already taken
		if checkDuplicateUser(user.Username) {
			c.JSON(
				http.StatusBadRequest,
				responses.GeneralResponse{
					Status:  http.StatusBadRequest,
					Message: utils.ErrorMessage,
					Data:    "Username already taken !",
				},
			)
			return
		}

		update := bson.M{
			"name":       user.Name,
			"username":   user.Username,
			"password":   utils.GetMD5Hash(user.Password),
			"updated_at": time.Now(),
		}

		result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
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

		//get updated user details
		var updatedUser models.User
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
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
		}

		c.JSON(
			http.StatusOK,
			responses.GeneralResponse{
				Status:  http.StatusOK,
				Message: utils.SuccessMessage,
				Data:    updatedUser,
			},
		)
	}
}

func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
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

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.GeneralResponse{
					Status:  http.StatusNotFound,
					Message: utils.ErrorMessage,
					Data:    "User with specified ID not found!",
				},
			)
			return
		}

		c.JSON(http.StatusOK,
			responses.GeneralResponse{
				Status:  http.StatusOK,
				Message: utils.SuccessMessage,
				Data:    "User successfully deleted!",
			},
		)
	}
}

func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users []models.User
		defer cancel()

		_, err := ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(
				http.StatusUnauthorized,
				responses.GeneralResponse{
					Status:  http.StatusUnauthorized,
					Message: utils.UnauthorizedMessage,
					Data:    err.Error(),
				},
			)
			return
		}

		results, err := userCollection.Find(ctx, bson.M{})

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

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.User
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(
					http.StatusInternalServerError,
					responses.GeneralResponse{
						Status:  http.StatusInternalServerError,
						Message: utils.ErrorMessage,
						Data:    err.Error(),
					},
				)
			}

			users = append(users, singleUser)
		}

		c.JSON(http.StatusOK,
			responses.GeneralResponse{
				Status:  http.StatusOK,
				Message: utils.SuccessMessage,
				Data:    users,
			},
		)
	}
}

func checkDuplicateUser(username string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	return err == nil
}
