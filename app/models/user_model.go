package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User model.
type User struct {
	Id       primitive.ObjectID `json:"id,omitempty"`
	Email    string             `json:"email,omitempty" validate:"required,email"`
	Password string             `json:"password,omitempty" validate:"required,min=6"`
}
