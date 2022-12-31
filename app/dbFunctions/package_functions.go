package dbFunctions

import (
	"context"
	"fmt"

	"github.com/bearts/go-fiber/app/database"
	"github.com/bearts/go-fiber/app/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var packageCollection *mongo.Collection = database.GetCollection(database.DB, "packages")

func CreatePackage(packageobj models.Package) (*mongo.InsertOneResult, error) {
	return packageCollection.InsertOne(context.Background(), packageobj)
}

func GetAllPackagesOfUser(uid primitive.ObjectID, status string) ([]models.Package, error) {
	var packages []models.Package
	if status != "" {
		cursor, err := packageCollection.Find(context.Background(), bson.M{"user": uid, "status": status})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(context.Background())
		for cursor.Next(context.Background()) {
			var packageobj models.Package
			if err := cursor.Decode(&packageobj); err != nil {
				return nil, err
			}
			packages = append(packages, packageobj)
		}
	} else {
		cursor, err := packageCollection.Find(context.Background(), bson.M{"user": uid})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(context.Background())
		for cursor.Next(context.Background()) {
			var packageobj models.Package
			if err := cursor.Decode(&packageobj); err != nil {
				return nil, err
			}
			packages = append(packages, packageobj)
		}
	}
	return packages, nil
}

func GetPackageById(id primitive.ObjectID) (models.Package, error) {
	var packageobj models.Package
	err := packageCollection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&packageobj)
	if err != nil {
		return models.Package{}, err
	}
	return packageobj, nil
}

func UpdatePackageDeliveryStatus(id primitive.ObjectID, status string, uid *primitive.ObjectID) (*mongo.UpdateResult, error) {
	fmt.Println("uid", uid)
	if *uid != primitive.NilObjectID {
		// update Package.Package.status
		return packageCollection.UpdateOne(context.Background(), bson.M{"_id": id, "user": uid}, bson.M{"$set": bson.M{"package.status": status}})
	}
	return packageCollection.UpdateOne(context.Background(), bson.M{"_id": id}, bson.M{"$set": bson.M{"package.status": status}})
}

func GetAllUnAssignedPackages() ([]models.Package, error) {
	var packages []models.Package
	cursor, err := packageCollection.Find(context.Background(), bson.M{"status": "pending"})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var packageobj models.Package
		if err := cursor.Decode(&packageobj); err != nil {
			return nil, err
		}
		packageobj.RunnerOtp = 0
		*packageobj.Package.Otp = 0

		packages = append(packages, packageobj)
	}
	return packages, nil
}

func GetAllPackageByRunner(rid primitive.ObjectID) ([]models.Package, error) {
	var packages []models.Package
	cursor, err := packageCollection.Find(context.Background(), bson.M{"runner": rid})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var packageobj models.Package
		if err := cursor.Decode(&packageobj); err != nil {
			return nil, err
		}
		packageobj.RunnerOtp = 0
		packages = append(packages, packageobj)
	}
	return packages, nil
}

func GetAllPackageByRunnerByStatus(rid primitive.ObjectID, status string) ([]models.Package, error) {
	var packages []models.Package
	cursor, err := packageCollection.Find(context.Background(), bson.M{"runner": rid, "status": status})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var packageobj models.Package
		if err := cursor.Decode(&packageobj); err != nil {
			return nil, err
		}
		packageobj.RunnerOtp = 0
		packages = append(packages, packageobj)
	}
	return packages, nil
}

func AssignPackage(pid primitive.ObjectID, rid primitive.ObjectID) (*mongo.UpdateResult, error) {
	return packageCollection.UpdateOne(context.Background(), bson.M{"_id": pid}, bson.M{"$set": bson.M{"runner": rid, "status": "assigned"}})
}

func UpdatePackageStatus(pid primitive.ObjectID, status string) (*mongo.UpdateResult, error) {
	return packageCollection.UpdateOne(context.Background(), bson.M{"_id": pid}, bson.M{"$set": bson.M{"status": status}})
}
