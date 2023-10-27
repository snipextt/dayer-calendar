package extension

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Extension struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Name        string             `json:"name" bson:"name"`
	Description string             `json:"description" bson:"description"`
	IconLight   string             `json:"iconLight" bson:"iconLight"`
	IconDark    string             `json:"iconDark" bson:"iconDark"`
}

type Extensions []Extension
