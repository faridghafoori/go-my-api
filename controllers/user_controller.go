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

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//validate the request body
		var user models.User
		err := c.BindJSON(&user)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		//use the validator library to validate required fields
		utils.ValidateStruct(&user)

		//check if the username is already taken
		if checkDuplicateUser(user.Username) {
			utils.GenerateErrorOutput(http.StatusBadRequest, errors.New("taken"), c, map[string]interface{}{
				"data": "Username already taken",
			})
		}

		newUser := models.User{
			Id:        primitive.NewObjectID(),
			Name:      user.Name,
			Username:  user.Username,
			Password:  utils.GetMD5Hash(user.Password),
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		_, err = userCollection.InsertOne(ctx, newUser)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		utils.GenerateSuccessOutput(newUser, c, http.StatusCreated)
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

func EditUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)
		//validate the request body
		var user models.User
		err := c.BindJSON(&user)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		//use the validator library to validate required fields
		utils.ValidateStruct(&user)

		//check if the username is already taken
		if checkDuplicateUser(user.Username) {
			utils.GenerateErrorOutput(http.StatusBadRequest, errors.New("taken"), c, map[string]interface{}{
				"data": "Username already taken",
			})
		}

		update := bson.M{
			"name":       user.Name,
			"username":   user.Username,
			"password":   utils.GetMD5Hash(user.Password),
			"updated_at": time.Now(),
		}

		result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
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
			utils.GenerateErrorOutput(http.StatusNotFound, err, c, map[string]interface{}{
				"data": "User with specified ID not found!",
			})
		}

		utils.GenerateSuccessOutput("User successfully deleted!", c)
	}
}

func GetAllUsers() gin.HandlerFunc {
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

func checkDuplicateUser(username string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	return err == nil
}
