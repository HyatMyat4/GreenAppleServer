package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Database_connection() *mongo.Client {

	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file could't be loaded")
	}

	Mongo_URI := os.Getenv("MONGODB_URI")

	client, err := mongo.NewClient(options.Client().ApplyURI(Mongo_URI))
	if err != nil {
		log.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	defer cancel()

	err = client.Connect(ctx)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("MongoDb Connected Successfully")

	return client
}

var MongoDb *mongo.Client = Database_connection()

func OpenCollection(mongodb *mongo.Client, collectionName string) *mongo.Collection {
	var collection *mongo.Collection = mongodb.Database("test").Collection(collectionName)

	return collection
}
