package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id        primitive.ObjectID `json:"id,omitempty"`
	Name      string             `json:"name,omitempty" validate:"required"`
	Username  string             `json:"username,omitempty" validate:"required"`
	Password  string             `json:"password,omitempty" validate:"required"`
	Roles     []Role             `json:"roles,omitempty"`
	Addresses []Address          `json:"addresses,omitempty"`
	CreatedAt time.Time          `json:"created_at,omitempty"`
	UpdatedAt time.Time          `json:"updated_at,omitempty"`
}

type UserWithRoleId struct {
	Name      string    `json:"name,omitempty" validate:"required"`
	Username  string    `json:"username,omitempty" validate:"required"`
	Password  string    `json:"password,omitempty" validate:"required"`
	RoleIds   []string  `json:"role_ids,omitempty"`
	Addresses []Address `json:"addresses,omitempty"`
}
