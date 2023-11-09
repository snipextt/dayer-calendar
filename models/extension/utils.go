package extension

import (
	"github.com/snipextt/dayer/storage"
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func collection() *mongo.Collection {
	return storage.GetMongoInstance().Collection("extensions")
}

func PaginatedExtensions(page int64, limit int64) (extensions Extensions, err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	cursor, err := collection().Find(ctx, bson.M{}, options.Find().SetSkip((page-1)*limit).SetLimit(limit))
	if err != nil {
		return
	}
	err = cursor.All(ctx, &extensions)

	return
}
