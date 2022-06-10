package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	Id        primitive.ObjectID
	Title     string `validate:"required"`
	Street    string `validate:"required"`
	City      string `validate:"required"`
	State     string `validate:"required"`
	Zip       string `validate:"required"`
	Country   string `validate:"required"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
