package services

import (
	"context"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"gin-mongo-api/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

func FetchUser(userId string) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	objId, _ := primitive.ObjectIDFromHex(userId)
	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
	return user, err
}

func CheckDuplicateUser(username string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := userCollection.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	return err == nil
}

func HandleExistedUserRoles(user models.User, inputUser models.UserInputBody) []string {
	if inputUser.RoleIds != nil {
		return inputUser.RoleIds
	}
	var userRolesIds []string
	for _, roles := range user.Roles {
		userRolesIds = append(userRolesIds, roles.Id.Hex())
	}
	return userRolesIds
}

func FillNewUserWithDefaultValues(inputUser models.UserInputBody) models.User {
	return models.User{
		Id:         primitive.NewObjectID(),
		Name:       inputUser.Name,
		Username:   inputUser.Username,
		TotpActive: inputUser.TotpActive,
		TotpKey:    inputUser.TotpKey,
		Password:   utils.GetSHA256Hash(inputUser.Password),
		Roles:      GenerateUserRoles(inputUser.RoleIds),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}
