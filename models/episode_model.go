package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Episode struct {
	Id          primitive.ObjectID
	Title       string `validate:"required"`
	Description string
	Image       string
	Background  string
	Rate        float64
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
