package connections

import (
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindConnectionsForUid(uid string) (connections []Model, error error) {
	ctx, cancel := utils.GetContext()
	defer cancel()
	res, err := collection().Find(ctx, bson.M{"uid": uid})
	if err != nil {
		return
	}
	err = res.All(ctx, &connections)
	return
}

func FindById(id string) (conn Model, err error) {
	ctx, cancel := utils.GetContext()
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	err = collection().FindOne(ctx, bson.M{"_id": oid}).Decode(&conn)
	return
}

func NewConnection(wid string, vid string, provider string, token string) (connection Model) {
	connection = Model{
		WorkspaceId: wid,
		VendorID:    vid,
		Provider:    provider,
		Token:       token,
	}
	return connection
}
