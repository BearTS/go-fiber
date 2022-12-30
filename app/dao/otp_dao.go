package dao

import (
	"context"

	"github.com/bearts/go-fiber/app/database"
	"github.com/bearts/go-fiber/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var otpCollection *mongo.Collection = database.GetCollection(database.DB, "otp")

func CreateOtp(otp models.Otp) (*mongo.InsertOneResult, error) {
	return otpCollection.InsertOne(context.Background(), otp)
}

func FindOtpByUser(user models.User) (*models.Otp, error) {
	var obj models.Otp
	err := otpCollection.FindOne(context.Background(), bson.M{"user": user.Id}).Decode(&obj)
	if err != nil {
		return nil, err
	}
	return &obj, nil
}

func DeleteOtpByUser(user models.User) (*mongo.DeleteResult, error) {
	return otpCollection.DeleteOne(context.Background(), bson.M{"user": user.Id})
}
