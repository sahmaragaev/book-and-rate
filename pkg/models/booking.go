package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Booking struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"userId"`
	RestaurantID primitive.ObjectID `bson:"restaurantId"`
	Date         time.Time          `bson:"date"`
	Cancelled    bool               `bson:"cancelled"`
}
