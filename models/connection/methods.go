package connection

import (
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func FindConnectionsForUid(uid string) (connections []Model, error error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	res, err := collection().Find(ctx, bson.M{"uid": uid})
	if err != nil {
		return
	}
	err = res.All(ctx, &connections)
	return
}

func ById(id string) (conn Model, err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	err = collection().FindOne(ctx, bson.M{"_id": oid}).Decode(&conn)
	return
}

func ForWorkspace(id primitive.ObjectID) (conn []Model, err error) {
  ctx, cancel := utils.NewContext()
  defer cancel()
  if err != nil {
    return
  }
  cur, err := collection().Find(ctx, bson.M{"workspace": id})
  if err != nil  {
    return
  }
  err = cur.All(ctx, &conn)
  return
}

func ForWorkspaceByProvider(id primitive.ObjectID, provider string) (conn Model, err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	if err != nil {
		return
	}
	err = collection().FindOne(ctx, bson.M{"workspace": id, "provider": provider}).Decode(&conn)
	return
}
