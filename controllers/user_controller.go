package controllers

import (
	"context"
	"errors"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/services"
	"gin-mongo-api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserController struct{}

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func (u UserController) GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users []models.User
		defer cancel()

		results, err := userCollection.Find(ctx, bson.M{})
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		//reading from the db in an optimal way
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.User
			err = results.Decode(&singleUser)
			utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)
			users = append(users, singleUser)
		}

		utils.GenerateSuccessOutput(users, c)
	}
}

func (u UserController) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("userId")

		user, err := services.FetchUser(userId)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		utils.GenerateSuccessOutput(user, c)
	}
}

func (u UserController) CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//validate the request body
		var inputUser models.UserInputBody
		err := c.BindJSON(&inputUser)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		//use the validator library to validate required fields
		validationErr := validate.Struct(&inputUser)
		utils.GenerateErrorOutput(http.StatusBadRequest, validationErr, c)

		//check if the username is already taken
		if services.CheckDuplicateUser(inputUser.Username) {
			utils.GenerateErrorOutput(http.StatusBadRequest, errors.New("taken"), c, map[string]interface{}{
				"data": "Username already taken",
			})
		}

		//create the user
		newUser := services.FillNewUserWithDefaultValues(inputUser)

		_, err = userCollection.InsertOne(ctx, newUser)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		utils.GenerateSuccessOutput(newUser, c, http.StatusCreated)
	}
}

func (u UserController) EditUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		userObjId, _ := primitive.ObjectIDFromHex(userId)

		findedUser, err := services.FetchUser(userId)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		// validate the request body
		var inputUser models.UserInputBody
		err = c.BindJSON(&inputUser)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		//use the validator library to validate required fields
		validationErr := validate.Struct(&inputUser)
		utils.GenerateErrorOutput(http.StatusBadRequest, validationErr, c)

		//check if the username is already taken
		if findedUser.Username != inputUser.Username && services.CheckDuplicateUser(inputUser.Username) {
			utils.GenerateErrorOutput(
				http.StatusBadRequest,
				errors.New("taken"),
				c,
				map[string]interface{}{
					"data": "Username already taken",
				},
			)
		}

		//update the user
		updatedUser := models.User{
			Id:         userObjId,
			Name:       inputUser.Name,
			Username:   inputUser.Username,
			Password:   utils.GetSHA256Hash(inputUser.Password),
			Roles:      services.GenerateUserRoles(services.HandleExistedUserRoles(findedUser, inputUser)),
			TotpActive: inputUser.TotpActive,
			Addresses:  findedUser.Addresses,
			TotpKey:    findedUser.TotpKey,
			CreatedAt:  findedUser.CreatedAt,
			UpdatedAt:  time.Now(),
		}

		_, err = userCollection.UpdateOne(ctx, bson.M{"id": userObjId}, bson.M{"$set": updatedUser})
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		utils.GenerateSuccessOutput(updatedUser, c)
	}
}

func (u UserController) DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		if result.DeletedCount < 1 {
			utils.GenerateErrorOutput(
				http.StatusNotFound,
				errors.New("user not found"),
				c,
				map[string]interface{}{
					"data": "User with specified ID not found!",
				},
			)
		}

		utils.GenerateSuccessOutput("User successfully deleted!", c)
	}
}
