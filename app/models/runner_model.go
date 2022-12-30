package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Runner model.

type Review struct {
	Id     primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	User   primitive.ObjectID `bson:"user,omitempty" json:"user,omitempty"`
	Text   string             `bson:"text,omitempty" json:"text,omitempty"`
	Rate   int                `bson:"rate,omitempty" json:"rate,omitempty"`
	Runner primitive.ObjectID `bson:"runner,omitempty" json:"runner,omitempty"`
}

type Runner struct {
	Id       primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name     string             `bson:"name,omitempty" json:"name,omitempty" default:""`
	Phone    string             `bson:"phone,omitempty" json:"phone,omitempty" default:"" validate:"required"`
	Email    string             `bson:"email,omitempty" json:"email,omitempty" validate:"required,email"`
	Password string             `bson:"password,omitempty" json:"password,omitempty" validate:"required"`
	Photo    string             `bson:"photo,omitempty" json:"photo,omitempty" default:"http://s3.amazonaws.com/37assets/svn/765-default-avatar.png"`
	IsAval   bool               `bson:"is_aval,omitempty" json:"isAvailable,omitempty" default:"true"`
}
