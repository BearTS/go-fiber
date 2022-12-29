package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// User model.
type User struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Email    string             `bson:"email,omitempty" json:"email,omitempty" validate:"required,email"` // this is not returned to the user
	Password string             `bson:"password,omitempty" json:"password,omitempty"`                     // this is not returned to the user
}
