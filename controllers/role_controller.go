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
	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var roleCollection *mongo.Collection = configs.GetCollection(configs.DB, "roles")
var validate = validator.New()

func GetRoles() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		results, err := roleCollection.Find(ctx, bson.M{})
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		var roles []models.Role
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleRole models.Role
			err = results.Decode(&singleRole)
			utils.GenerateErrorOutput(http.StatusBadRequest, err, c)
			roles = append(roles, singleRole)
		}

		utils.GenerateSuccessOutput(roles, c)
	}
}

func GetRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		roleId := c.Param("roleId")
		defer cancel()

		roleObjId, _ := primitive.ObjectIDFromHex(roleId)
		var role models.Role
		err := roleCollection.FindOne(ctx, bson.M{"id": roleObjId}).Decode(&role)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		utils.GenerateSuccessOutput(role, c)
	}
}

func CreateRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		//validate the request body
		var role models.Role
		err := c.BindJSON(&role)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		//use the validator library to validate required fields
		validationErr := validate.Struct(&role)
		utils.GenerateErrorOutput(http.StatusBadRequest, validationErr, c)

		newRole := models.Role{
			Id:         primitive.NewObjectID(),
			Name:       role.Name,
			Descriptor: role.Descriptor,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		}

		_, err = roleCollection.InsertOne(ctx, newRole)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		utils.GenerateSuccessOutput(newRole, c)
	}
}

func EditRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		roleId := c.Param("roleId")
		defer cancel()

		roleObjId, _ := primitive.ObjectIDFromHex(roleId)

		var role models.Role
		err := roleCollection.FindOne(ctx, bson.M{"id": roleObjId}).Decode(&role)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		err = c.BindJSON(&role)
		utils.GenerateErrorOutput(http.StatusBadRequest, err, c)

		validationErr := validate.Struct(&role)
		utils.GenerateErrorOutput(http.StatusBadRequest, validationErr, c)

		role.UpdatedAt = time.Now()

		_, err = roleCollection.UpdateOne(ctx, bson.M{"id": roleObjId}, bson.M{"$set": role})
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		utils.GenerateSuccessOutput(role, c)
	}
}

func DeleteRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		roleId := c.Param("roleId")
		defer cancel()

		roleObjId, _ := primitive.ObjectIDFromHex(roleId)
		result, err := roleCollection.DeleteOne(ctx, bson.M{"id": roleObjId})
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		if result.DeletedCount < 1 {
			utils.GenerateErrorOutput(
				http.StatusNotFound,
				errors.New("role not found"),
				c,
				map[string]interface{}{
					"data": "Role with specified ID not found!",
				},
			)
		}

		utils.GenerateSuccessOutput("Role successfully deleted!", c)
	}
}
