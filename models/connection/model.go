package connection

import (
	"github.com/snipextt/dayer/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Model struct {
	Id          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	Extension   string             `json:"extension" bson:"extension"`
	WorkspaceId primitive.ObjectID `json:"-" bson:"workspaceId"`
	VendorID    string             `json:"email" bson:"email"`
	Provider    string             `json:"provider" bson:"provider"`
	Token       string             `bson:"token" json:"-"`
	ExpiresAt   string             `bson:"expiresAt" json:"-"`
}

func collection() *mongo.Collection {
	return storage.GetMongoInstance().Collection("connections")
}

func (c *Model) Save() error {
	if c.Id.IsZero() {
		return c.Create()
	}
	_, err := collection().UpdateOne(nil, primitive.M{"_id": c.Id}, primitive.M{"$set": c})
	return err
}

func (c *Model) Create() (err error) {
	conn, err := collection().InsertOne(nil, c)
	c.Id = conn.InsertedID.(primitive.ObjectID)
	return
}
