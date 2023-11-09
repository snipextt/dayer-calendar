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

func FindById(id string) (conn Model, err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}
	err = collection().FindOne(ctx, bson.M{"_id": oid}).Decode(&conn)
	return
}

func FindByWorkspaceId(id, provider string) (conn Model, err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return
	}

	err = collection().FindOne(ctx, bson.M{"workspaceId": oid, "provider": provider}).Decode(&conn)
	return
}

func NewCalendarConnection(wid string, vid string, provider string, token string) (connection Model) {
	oid, _ := primitive.ObjectIDFromHex(wid)
	connection = Model{
		WorkspaceId: oid,
		VendorID:    vid,
		Provider:    provider,
		Token:       token,
	}
	return connection
}

func NewTimeDoctorConnection(wid string, vid string, token string, expiresAt string) (connection Model) {
	oid, _ := primitive.ObjectIDFromHex(wid)
	connection = Model{
		Extension:   "timedoctor",
		WorkspaceId: oid,
		VendorID:    vid,
		Token:       token,
		ExpiresAt:   expiresAt,
		Provider:    "timedoctor",
	}
	return connection
}
