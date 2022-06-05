package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role struct {
	Id         primitive.ObjectID `json:"id,omitempty"`
	Name       string             `json:"name,omitempty" validate:"required"`
	Descriptor string             `json:"descriptor,omitempty" validate:"required"`
	CreatedAt  time.Time          `json:"created_at,omitempty"`
	UpdatedAt  time.Time          `json:"updated_at,omitempty"`
}
