package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type RefreshToken struct {
	Id        primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Token     string             `bson:"token,omitempty" json:"token,omitempty"`
	User      primitive.ObjectID `bson:"user,omitempty" json:"user,omitempty"`
	ExpiresAt int64              `bson:"expires_at,omitempty" json:"expires_at,omitempty"`
}
