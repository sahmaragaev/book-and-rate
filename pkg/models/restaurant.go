package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Restaurant struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Name     string             `bson:"name"`
	Address  string             `bson:"address"`
	Phone    string             `bson:"phone"`
	Password string             `bson:"password"`
}
