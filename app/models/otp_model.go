package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Otp struct {
	Id   primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Otp  int                `bson:"otp,omitempty" json:"otp,omitempty"`
	User primitive.ObjectID `bson:"user,omitempty" json:"user,omitempty"`
}
