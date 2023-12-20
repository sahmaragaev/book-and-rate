package db

import "go.mongodb.org/mongo-driver/mongo"

var (
	UserCollection       *mongo.Collection
	RestaurantCollection *mongo.Collection
	BookingCollection    *mongo.Collection
	RateCollection       *mongo.Collection
)

var db = Client.Database("bookandrate")

func InitializeCollections() {
	UserCollection = db.Collection("users")
	RestaurantCollection = db.Collection("restaurants")
	BookingCollection = db.Collection("bookings")
	RateCollection = db.Collection("rates")
}
