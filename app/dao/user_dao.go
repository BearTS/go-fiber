package dao

import (
	"context"

	"github.com/bearts/go-fiber/app/database"
	"github.com/bearts/go-fiber/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "User")

func CreateUser(user models.User) (*mongo.InsertOneResult, error) {
	return userCollection.InsertOne(context.Background(), user)
}

func GetUserById(id primitive.ObjectID) (*models.User, error) {
	var user models.User
	err := userCollection.FindOne(context.Background(), models.User{Id: id}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := userCollection.FindOne(context.Background(), models.User{Email: email}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func UpdateUser(user models.User) (*mongo.UpdateResult, error) {
	return userCollection.UpdateOne(context.Background(), models.User{Email: user.Email}, bson.M{"$set": user})
}
