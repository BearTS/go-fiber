package models

import "go.mongodb.org/mongo-driver/bson/primitive"

// Place model.
type Place struct {
	Id          primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	Name        string             `bson:"name,omitempty" json:"name,omitempty"`
	Description string             `bson:"description,omitempty" json:"description,omitempty"`
	Code        string             `bson:"code,omitempty" json:"code,omitempty"`
}

type Price struct {
	Id    primitive.ObjectID `bson:"_id,omitempty" json:"id,omitempty"`
	From  primitive.ObjectID `bson:"from,omitempty" json:"from,omitempty"`
	To    primitive.ObjectID `bson:"to,omitempty" json:"to,omitempty"`
	Price int                `bson:"price,omitempty" json:"price,omitempty"`
}
