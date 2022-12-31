package dbFunctions

import (
	"context"

	"github.com/bearts/go-fiber/app/database"
	"github.com/bearts/go-fiber/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var priceCollection *mongo.Collection = database.GetCollection(database.DB, "prices")

func GetPriceFromTo(from, to string) (int, error) {
	var price models.Price
	_from, err := GetPlaceByCode(from)
	if err != nil {
		return 0, err
	}
	_to, err := GetPlaceByCode(to)
	if err != nil {
		return 0, err
	}
	if err := priceCollection.FindOne(context.Background(), bson.M{"from": _from.Id, "to": _to.Id}).Decode(&price); err != nil {
		return 0, err
	}
	return price.Price, nil
}

func GetPriceFromToById(from, to primitive.ObjectID) (int, error) {
	var price models.Price
	if err := priceCollection.FindOne(context.Background(), bson.M{"from": from, "to": to}).Decode(&price); err != nil {
		return 0, err
	}
	return price.Price, nil
}
