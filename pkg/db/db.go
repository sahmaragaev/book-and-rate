package db

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var Client *mongo.Client

func Connect(connectionString string) {
	var err error

	Client, err = mongo.Connect(context.TODO(), options.Client().ApplyURI(connectionString))
	if err != nil {
		log.Fatal(err)
	}

	err = Client.Ping(context.TODO(), nil)
	log.Println("Connected to MongoDB!")
}

func GetDbCollection(collectionName string) *mongo.Collection {
	return Client.Database("bookandrate").Collection(collectionName)
}
