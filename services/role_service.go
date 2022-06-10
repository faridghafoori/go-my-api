package services

import (
	"context"
	"gin-mongo-api/configs"
	"gin-mongo-api/models"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var roleCollection *mongo.Collection = configs.GetCollection(configs.DB, "roles")

func GenerateUserRoles(roleIds []string) []models.Role {
	var roles []models.Role
	ch := make(chan models.Role)

	if roleIds != nil {
		go FetchRole(roleIds, ch)
	} else {
		// normal user
		go FetchRole([]string{"62a3363d5f766861278f7a0c"}, ch)
	}

	for l := range ch {
		roles = append(roles, l)
	}

	return roles
}

func FetchRole(roleIds []string, ch chan models.Role) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for _, roleId := range roleIds {
		roleObjId, _ := primitive.ObjectIDFromHex(roleId)
		var role models.Role
		err := roleCollection.FindOne(ctx, bson.M{"id": roleObjId}).Decode(&role)
		if err != nil {
			return
		}
		ch <- role
	}

	close(ch)
}
