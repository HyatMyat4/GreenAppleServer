package database

import (
	"go.mongodb.org/mongo-driver/mongo"
)

var UsersCollection *mongo.Collection = OpenCollection(MongoDb, "users")

var CompanyCollection *mongo.Collection = OpenCollection(MongoDb, "company")
