package user

import (
	"github.com/snipextt/dayer/storage"
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func New(cid string) (user Model) {
	user = Model{
		ClerkId:    cid,
		Active:     true,
		Workspaces: []interface{}{},
	}
	return
}

func collection() *mongo.Collection {
	return storage.GetMongoInstance().Collection("users")
}

func FindById(id string) (user Model, err error) {
	ctx, cancel := utils.GetContext()
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(id)
	err = collection().FindOne(ctx, bson.M{"_id": oid}).Decode(&user)
	return
}
