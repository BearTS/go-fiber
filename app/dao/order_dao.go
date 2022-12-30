package dao

import (
	"context"
	"fmt"

	"github.com/bearts/go-fiber/app/database"
	"github.com/bearts/go-fiber/app/models"
	"github.com/bearts/go-fiber/app/structs"
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

func GetOrderById(id primitive.ObjectID, userid primitive.ObjectID) (*models.Order, error) {
	var order models.Order
	if userid != primitive.NilObjectID {
		if err := orderCollection.FindOne(context.Background(), bson.M{"_id": id, "user": userid}).Decode(&order); err != nil {
			return nil, err
		}
	} else {
		if err := orderCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&order); err != nil {
			return nil, err
		}
	}
	return &order, nil
}

func AssignOrderById(id primitive.ObjectID, runnerid primitive.ObjectID) (*mongo.UpdateResult, error) {
	var order models.Order
	if err := orderCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&order); err != nil {
		return nil, err
	}
	if order.Status != "pending" {
		return nil, fmt.Errorf("order is not pending")
	}
	updateResult, err := orderCollection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"status": "assigned", "runner": runnerid}})
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

func UpdateOrderStatus(id primitive.ObjectID, runnerid primitive.ObjectID, status string) (*mongo.UpdateResult, error) {
	var order models.Order
	if err := orderCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&order); err != nil {
		return nil, err
	}
	if *order.Runner != runnerid {
		return nil, fmt.Errorf("you are not assigned to this order")
	}
	if order.Status == status {
		return nil, fmt.Errorf("order is already %s", status)
	}
	updateResult, err := orderCollection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"status": status}})
	if err != nil {
		return nil, err
	}
	return updateResult, nil
}

func GetAllUnassignedOrders() ([]structs.GetAllUnassignedOrders, error) {
	var orders []structs.GetAllUnassignedOrders
	cursor, err := orderCollection.Find(context.Background(), bson.M{"status": "pending"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var order models.Order
		if err := cursor.Decode(&order); err != nil {
			return nil, err
		}
		order.Delivery_app.DeliveryPhone = 0
		order.RunnerOtp = 0
		order.Delivery_app.Otp = 0
		// fill location
		location, err := GetPlaceById(order.Location)
		if err != nil {
			return nil, err
		}
		// add order and location into 1
		var data structs.GetAllUnassignedOrders
		data.Order = order
		data.Location = *location

		orders = append(orders, data)
	}
	return orders, nil
}

func GetAllPreviousOrders(runnerId primitive.ObjectID) ([]models.Order, error) {
	var orders []models.Order
	cursor, err := orderCollection.Find(context.Background(), bson.M{"runner": runnerId, "status": bson.M{"$in": []string{"delivered", "cancelled"}}})
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
	return orders, nil
}

func GetAllCurrentOrders(runnerId primitive.ObjectID) ([]models.Order, error) {
	var orders []models.Order
	cursor, err := orderCollection.Find(context.Background(), bson.M{"runner": runnerId, "status": bson.M{"$in": []string{"assigned", "pickedup"}}})
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
	return orders, nil
}
