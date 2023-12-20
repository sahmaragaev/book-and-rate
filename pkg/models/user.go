package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID          primitive.ObjectID `bson:"_id,omitempty"`
	FirstName   string             `bson:"firstName"`
	LastName    string             `bson:"lastName"`
	PhoneNumber string             `bson:"phoneNumber"`
	Password    string             `bson:"password"`
}
