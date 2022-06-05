package controllers

import (
	"context"
	"gin-mongo-api/models"
	"gin-mongo-api/utils"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetUserAddresses() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		utils.GenerateSuccessOutput(user.Addresses, c)
	}
}

func AddNewAddressToUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		userObjId, _ := primitive.ObjectIDFromHex(userId)

		var user models.User
		err := userCollection.FindOne(ctx, bson.M{"id": userObjId}).Decode(&user)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		//validate the request body
		var address models.Address
		err = c.BindJSON(&address)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		//use the validator library to validate required fields
		utils.ValidateStruct(&address)

		newAddress := models.Address{
			Id:        primitive.NewObjectID(),
			Title:     address.Title,
			Street:    address.Street,
			City:      address.City,
			State:     address.State,
			Zip:       address.Zip,
			Country:   address.Country,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		user.Addresses = append(user.Addresses, newAddress)

		update := bson.M{
			"name":       user.Name,
			"addresses":  user.Addresses,
			"username":   user.Username,
			"password":   utils.GetMD5Hash(user.Password),
			"updated_at": time.Now(),
		}

		result, err := userCollection.UpdateOne(ctx, bson.M{"id": userObjId}, bson.M{"$set": update})
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		if result.MatchedCount == 1 {
			utils.GenerateSuccessOutput(user, c)
		}
	}
}

func EditAddressOfUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		addressId := c.Param("addressId")
		defer cancel()

		userObjId, _ := primitive.ObjectIDFromHex(userId)
		addressObjId, _ := primitive.ObjectIDFromHex(addressId)
		// FIXME: this query must be fixed
		filter := bson.M{
			"id":           userObjId,
			"addresses.id": addressObjId,
		}
		var user models.User
		err := userCollection.FindOne(ctx, filter).Decode(&user)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		var selectedUserAddressIndex int
		for i, address := range user.Addresses {
			if address.Id == addressObjId {
				selectedUserAddressIndex = i
			}
		}

		var inputAddress models.Address
		//validate the request body
		err = c.BindJSON(&inputAddress)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		//use the validator library to validate required fields
		utils.ValidateStruct(&inputAddress)

		user.Addresses[selectedUserAddressIndex].City = inputAddress.City
		user.Addresses[selectedUserAddressIndex].Country = inputAddress.Country
		user.Addresses[selectedUserAddressIndex].State = inputAddress.State
		user.Addresses[selectedUserAddressIndex].Street = inputAddress.Street
		user.Addresses[selectedUserAddressIndex].Title = inputAddress.Title
		user.Addresses[selectedUserAddressIndex].Zip = inputAddress.Zip
		user.Addresses[selectedUserAddressIndex].UpdatedAt = time.Now()

		update := bson.M{
			"name":       user.Name,
			"addresses":  user.Addresses,
			"username":   user.Username,
			"password":   utils.GetMD5Hash(user.Password),
			"updated_at": time.Now(),
		}

		result, err := userCollection.UpdateOne(ctx, bson.M{"id": userObjId}, bson.M{"$set": update})
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		if result.MatchedCount == 1 {
			utils.GenerateSuccessOutput(user, c)
		}
	}
}

func DeleteAddressOfUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		addressId := c.Param("addressId")
		defer cancel()

		userObjId, _ := primitive.ObjectIDFromHex(userId)
		addressObjId, _ := primitive.ObjectIDFromHex(addressId)
		// FIXME: this query must be fixed
		filter := bson.M{
			"id":           userObjId,
			"addresses.id": addressObjId,
		}
		var user models.User
		err := userCollection.FindOne(ctx, filter).Decode(&user)
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		var selectedUserAddressIndex int
		for i, address := range user.Addresses {
			if address.Id == addressObjId {
				selectedUserAddressIndex = i
			}
		}

		user.Addresses = append(user.Addresses[:selectedUserAddressIndex], user.Addresses[selectedUserAddressIndex+1:]...)

		update := bson.M{
			"name":       user.Name,
			"addresses":  user.Addresses,
			"username":   user.Username,
			"password":   utils.GetMD5Hash(user.Password),
			"updated_at": time.Now(),
		}

		result, err := userCollection.UpdateOne(ctx, bson.M{"id": userObjId}, bson.M{"$set": update})
		utils.GenerateErrorOutput(http.StatusInternalServerError, err, c)

		if result.MatchedCount == 1 {
			utils.GenerateSuccessOutput(user, c)
		}
	}
}
