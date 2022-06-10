package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	Id         primitive.ObjectID
	Name       string `validate:"required"`
	Descriptor string `validate:"required"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}
