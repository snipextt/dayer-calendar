package timedoctor

import (
	"time"

	"github.com/snipextt/dayer/models"
	"github.com/snipextt/dayer/storage"
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func collection() *mongo.Collection {
  return storage.Primary().Collection("timedoctorReports")
}

func ReportForWorkspace(wid, team primitive.ObjectID, startDate, endDate time.Time) (reports []models.TimeDoctorReport, err error) {
  ctx, cancel := utils.NewContext()
  defer cancel()

  res, err := collection().Find(ctx, bson.M{"workspace": wid, "createdAt": bson.M{"$gte": primitive.NewDateTimeFromTime(startDate), "$lt": primitive.NewDateTimeFromTime(endDate)}, "teams": team })
  if err != nil {
    return
  }
  err = res.All(ctx, &reports)
  return
}
