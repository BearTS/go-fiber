package dao

import (
	"context"
	"fmt"

	"github.com/bearts/go-fiber/app/database"
	"github.com/bearts/go-fiber/app/models"
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

func CreateOrder(id string, NameOfApp string, NameOfRes string, EstimatedTime int, deliveryPhone int, location string, otp int) (*mongo.InsertOneResult, error) {
	price, err := GetPriceFromTo("main_gate", location)
	if err != nil {
		return nil, err
	}
	locationObj, err := GetPlaceByCode(location)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	// convert id to object id
	user, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	fmt.Println(user)
	// create order
	var order models.Order
	order.Delivery_app.NameOfApp = NameOfApp
	order.Delivery_app.NameOfRes = NameOfRes
	order.Delivery_app.EstimatedTime = EstimatedTime
	order.Delivery_app.DeliveryPhone = deliveryPhone
	order.Location = locationObj.Id
	if otp != 0 {
		order.Delivery_app.Otp = otp
	}
	order.Price = price
	order.Status = "pending"
	order.User = user
	orderObj, err := orderCollection.InsertOne(context.Background(), order)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return orderObj, nil
}
