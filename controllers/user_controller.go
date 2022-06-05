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

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func GetUsers() gin.HandlerFunc {
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

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)
		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		utils.GenerateSuccessOutput(user, c)
	}
}

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//validate the request body
		var inputUser models.UserWithRoleId
		err := c.BindJSON(&inputUser)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		//use the validator library to validate required fields
		validationErr := validate.Struct(&inputUser)
		utils.GenerateErrorOutput(http.StatusBadRequest, validationErr, c)

		//check if the username is already taken
		if checkDuplicateUser(inputUser.Username) {
			utils.GenerateErrorOutput(http.StatusBadRequest, errors.New("taken"), c, map[string]interface{}{
				"data": "Username already taken",
			})
		}

		newUser := models.User{
			Id:        primitive.NewObjectID(),
			Name:      inputUser.Name,
			Username:  inputUser.Username,
			Password:  utils.GetMD5Hash(inputUser.Password),
			Roles:     services.GenerateUserRoles(inputUser.RoleIds),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err = userCollection.InsertOne(ctx, newUser)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		utils.GenerateSuccessOutput(newUser, c, http.StatusCreated)
	}
}

func EditUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		//validate the request body
		var inputUser models.UserWithRoleId
		err = c.BindJSON(&inputUser)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		//use the validator library to validate required fields
		validationErr := validate.Struct(&inputUser)
		utils.GenerateErrorOutput(http.StatusBadRequest, validationErr, c)

		//check if the username is already taken
		if checkDuplicateUser(inputUser.Username) {
			utils.GenerateErrorOutput(
				http.StatusBadRequest,
				errors.New("taken"),
				c,
				map[string]interface{}{
					"data": "Username already taken",
				},
			)
		}
		user.Roles = services.GenerateUserRoles(inputUser.RoleIds)
		user.UpdatedAt = time.Now()

		result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": user})
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		//get updated user details
		var updatedUser models.User
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
			utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)
		}

		utils.GenerateSuccessOutput(updatedUser, c)
	}
}

func DeleteUser() gin.HandlerFunc {
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

func checkDuplicateUser(username string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	return err == nil
}
