package connection

import (
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func NewCalendarConnection(wid string, vid string, provider string, token string) (connection Model) {
	oid, _ := primitive.ObjectIDFromHex(wid)
	connection = Model{
		Workspace: oid,
		VendorID:  vid,
		Provider:  provider,
		Token:     token,
	}
	return connection
}

func NewTimeDoctorConnection(wid primitive.ObjectID, vid string, token string, expiresAt string) (connection Model) {
	connection = Model{
		Extension: "timedoctor",
		Workspace: wid,
		VendorID:  vid,
		Token:     token,
		ExpiresAt: expiresAt,
		Provider:  "timedoctor",
	}
	return connection
}

func FindConnectionsByProvider(provider string) (connections []Model, err error) {
  ctx, cancel := utils.NewContext()
  defer cancel()
  res, err := collection().Find(ctx, bson.M{"provider": provider})
  if err != nil {
    return
  }
  err = res.All(ctx, &connections)
  return
}

func GetTimedoctorConnections() (connections []Model, err error) {
  connections, err = FindConnectionsByProvider("timedoctor")
  return
}
