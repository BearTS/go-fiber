package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User model.
type User struct {
	Id                 primitive.ObjectID  `bson:"_id,omitempty" json:"id,omitempty"`
	Email              string              `bson:"email,omitempty" json:"email,omitempty" validate:"required,email"`
	Name               string              `bson:"name,omitempty" json:"name,omitempty" default:""`
	Phone              *string             `bson:"phone,omitempty" json:"phone,omitempty" default:""`
	RegistrationNumber *string             `bson:"registration_number,omitempty" json:"registration_number,omitempty" default:""`
	DefaultAddress     *primitive.ObjectID `bson:"default_address,omitempty" json:"default_address,omitempty"`
}
