package dbFunctions

import (
	"context"

	"github.com/bearts/go-fiber/app/database"
	"github.com/bearts/go-fiber/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var placeCollection *mongo.Collection = database.GetCollection(database.DB, "places")

func GetPlaceByCode(code string) (*models.Place, error) {
	var place models.Place
	err := placeCollection.FindOne(context.Background(), bson.M{"code": code}).Decode(&place)
	if err != nil {
		return nil, err
	}
	return &place, nil
}

func GetPlaceById(id primitive.ObjectID) (*models.Place, error) {
	var place models.Place
	err := placeCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&place)
	if err != nil {
		return nil, err
	}
	return &place, nil
}
