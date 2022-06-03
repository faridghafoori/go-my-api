package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Address struct {
	Id        primitive.ObjectID `json:"_id,omitempty"`
	Title     string             `json:"title,omitempty" validate:"required"`
	Street    string             `json:"street,omitempty" validate:"required"`
	City      string             `json:"city,omitempty" validate:"required"`
	State     string             `json:"state,omitempty" validate:"required"`
	Zip       string             `json:"zip,omitempty" validate:"required"`
	Country   string             `json:"country,omitempty" validate:"required"`
	CreatedAt time.Time          `json:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty"`
}
