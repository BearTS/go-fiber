package dao

import (
	"context"
	"fmt"

	"github.com/bearts/go-fiber/app/database"
	"github.com/bearts/go-fiber/app/models"
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
