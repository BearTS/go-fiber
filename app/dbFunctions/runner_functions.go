package dbFunctions

import (
	"context"

	"github.com/bearts/go-fiber/app/database"
	"github.com/bearts/go-fiber/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var runnerCollection *mongo.Collection = database.GetCollection(database.DB, "runner")

func GetAllRunner() ([]models.Runner, error) {
	var runners []models.Runner
	cursor, err := runnerCollection.Find(context.Background(), nil)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var runner models.Runner
		cursor.Decode(&runner)
		runners = append(runners, runner)
	}
	return runners, nil
}

func GetRunnerById(id primitive.ObjectID) (*models.Runner, error) {
	var runner models.Runner
	err := runnerCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&runner)
	if err != nil {
		return nil, err
	}
	return &runner, nil
}

func GetRunnerByEmail(email string) (*models.Runner, error) {
	var runner models.Runner
	err := runnerCollection.FindOne(context.Background(), bson.M{"email": email}).Decode(&runner)
	if err != nil {
		return nil, err
	}
	return &runner, nil
}

func CreateRunner(runner models.Runner) (*mongo.InsertOneResult, error) {
	return runnerCollection.InsertOne(context.Background(), runner)
}
