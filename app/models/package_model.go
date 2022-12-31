package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Package struct {
	Id      primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Package struct {
		NameOfApp     string             `bson:"nameOfApp,omitempty" json:"nameOfApp,omitempty" validate:"required"` // name of the delivery app
		Location      primitive.ObjectID `bson:"location,omitempty" json:"location,omitempty"`
		TrackingId    string             `bson:"trackingId,omitempty" json:"trackingId,omitempty" validate:"required"`
		Otp           *int               `bson:"otp,omitempty" json:"otp,omitempty"`
		DeliveryPhone *int               `bson:"delivery_Phone,omitempty" json:"delivery_Phone,omitempty"`
		Eta           *int               `bson:"eta,omitempty" json:"eta,omitempty"`
		Status        string             `bson:"status,omitempty" json:"status,omitempty"`
	} `bson:"package,omitempty" json:"package,omitempty"`
	DeliveryLocation string              `bson:"delivery_location,omitempty" json:"delivery_location,omitempty" validate:"required"`
	Price            int                 `bson:"price,omitempty" json:"price,omitempty"`
	Status           string              `bson:"status,omitempty" json:"status,omitempty"`
	Runner           *primitive.ObjectID `bson:"runner,omitempty" json:"runner,omitempty"`
	RunnerOtp        int                 `bson:"runner_otp,omitempty" json:"runner_otp,omitempty"`
}
