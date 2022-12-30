package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Order model.
type Order struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Delivery_app struct {
		NameOfApp     string `bson:"nameOfApp,omitempty" json:"nameOfApp,omitempty" validate:"required"`               // name of the delivery app
		NameOfRes     string `bson:"nameOfRestaurant,omitempty" json:"nameOfRestaurant,omitempty" validate:"required"` // name of the restaurant
		EstimatedTime int    `bson:"estimated_time,omitempty" json:"estimated_time,omitempty"`
		Otp           int    `bson:"otp,omitempty" json:"otp,omitempty"`
		DeliveryPhone int    `bson:"delivery_Phone,omitempty" json:"delivery_Phone,omitempty" validate:"required"`
	} `bson:"delivery_app,omitempty" json:"delivery_app,omitempty"`
	User      primitive.ObjectID  `bson:"user,omitempty" json:"user,omitempty"`
	Location  primitive.ObjectID  `bson:"location,omitempty" json:"location,omitempty"`
	Status    string              `bson:"status,omitempty" json:"status,omitempty"`
	Price     int                 `bson:"price,omitempty" json:"price,omitempty"`
	Runner    *primitive.ObjectID `bson:"runner,omitempty" json:"runner,omitempty"`
	CreatedAt time.Time           `bson:"created_at,omitempty" json:"created_at,omitempty"`
	UpdatedAt time.Time           `bson:"updated_at,omitempty" json:"updated_at,omitempty"`
}

// Order model.
