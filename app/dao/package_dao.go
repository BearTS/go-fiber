package dao

import (
	"context"

	"github.com/bearts/go-fiber/app/database"
	"github.com/bearts/go-fiber/app/models"
	"go.mongodb.org/mongo-driver/mongo"
)

var packageCollection *mongo.Collection = database.GetCollection(database.DB, "packages")

func CreatePackage(packageobj models.Package) (*mongo.InsertOneResult, error) {
	return packageCollection.InsertOne(context.Background(), packageobj)
}
