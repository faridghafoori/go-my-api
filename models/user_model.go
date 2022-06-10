package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id         primitive.ObjectID
	Name       string
	Password   string `validate:"required"` // json:"-"
	Username   string `validate:"required"`
	TotpActive bool   `validate:"required"`
	TotpKey    string
	Roles      []Role
	Addresses  []Address
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type UserInputBody struct {
	Name       string
	Username   string `validate:"required"`
	Password   string `validate:"required"`
	TotpActive bool
	RoleIds    []string
	Addresses  []Address
	TotpKey    string
}
