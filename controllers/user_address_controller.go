package controllers

import (
	"context"
	"gin-mongo-api/models"
	"gin-mongo-api/responses"
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
				Data:    user.Addresses,
			},
		)
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

		var address models.Address
		//validate the request body
		if err := c.BindJSON(&address); err != nil {
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
		if result.MatchedCount == 1 {
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

		var selectedUserAddressIndex int
		for i, address := range user.Addresses {
			if address.Id == addressObjId {
				selectedUserAddressIndex = i
			}
		}

		var inputAddress models.Address
		//validate the request body
		if err := c.BindJSON(&inputAddress); err != nil {
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
		if result.MatchedCount == 1 {
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
		if result.MatchedCount == 1 {
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
}
