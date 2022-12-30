package dao

import (
	"context"
	"fmt"

	"github.com/bearts/go-fiber/app/database"
	"github.com/bearts/go-fiber/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// create a function that takes {
//  NameOfApp: "app",
//  NameOfRestaurants: "restaurants",
//  estimatedTime: 30,
// otp: 1234 or optional
// deliveryPhone: 1234567890
// location: 'code'
// }

var orderCollection *mongo.Collection = database.GetCollection(database.DB, "order")

func CreateOrder(order *models.Order) (*mongo.InsertOneResult, error) {

	orderObj, err := orderCollection.InsertOne(context.Background(), order)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return orderObj, nil
}

func GetAllOrdersOfUser(id primitive.ObjectID, status string) ([]models.Order, error) {
	var orders []models.Order
	// status is optional
	if status != "" {
		cursor, err := orderCollection.Find(context.Background(), bson.M{"user": id, "status": status})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(context.Background())
		for cursor.Next(context.Background()) {
			var order models.Order
			if err := cursor.Decode(&order); err != nil {
				return nil, err
			}
			orders = append(orders, order)
		}

	} else {
		cursor, err := orderCollection.Find(context.Background(), bson.M{"user": id})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(context.Background())
		for cursor.Next(context.Background()) {
			var order models.Order
			if err := cursor.Decode(&order); err != nil {
				return nil, err
			}
			orders = append(orders, order)
		}
	}
	return orders, nil
}
