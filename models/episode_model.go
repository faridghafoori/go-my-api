package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Episode struct {
	Id          primitive.ObjectID `json:"id"`
	Title       string             `json:"title,omitempty" validate:"required"`
	Description string             `json:"description,omitempty"`
	Image       string             `json:"image,omitempty"`
	Background  string             `json:"background,omitempty"`
	Rate        float64            `json:"rate,omitempty"`
	CreatedAt   time.Time          `json:"created_at,omitempty"`
	UpdatedAt   time.Time          `json:"updated_at,omitempty"`
}
