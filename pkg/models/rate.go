package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Rate struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	UserID       primitive.ObjectID `bson:"userId"`
	RestaurantID primitive.ObjectID `bson:"restaurantId"`
	Rating       byte               `bson:"rating"`
	Comment      string             `bson:"comment"`
	Date         time.Time          `bson:"date"`
}
