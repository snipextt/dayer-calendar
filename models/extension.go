package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Extension struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	Icon        string             `json:"icon" bson:"icon"`
}
