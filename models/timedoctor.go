package models

import (
	"encoding/json"
	"time"

	"github.com/snipextt/dayer/storage"
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TimeDoctorActivity struct {
	Time     float64 `json:"time"`
	Category struct {
		Entity string `json:"entity"`
		Id     string `json:"id"`
		Name   string `json:"name"`
	} `json:"category"`
	Type  string `json:"type"`
	Title string `json:"title"`
	Value string `json:"value"`
	Start string `json:"start"`
}

type TimeDoctorImageData struct {
	Date    string `json:"date"`
	Numbers []struct {
		Url string `json:"url"`
	} `json:"numbers"`
}

type TimeDoctorReportForAnalysis struct {
  Id          primitive.ObjectID    `json:"_id" bson:"_id,omitempty"`
	Images      []TimeDoctorImageData `json:"images" bson:"images"`
	Activities  []TimeDoctorActivity  `json:"activities" bson:"activities"`
	Tasks       []string              `json:"tasks" bson:"tasks"`
	CreatedAt   time.Time             `json:"createdAt" bson:"createdAt"`
	MemberId    primitive.ObjectID    `json:"memberId" bson:"memberId"`
	WorkspaceId primitive.ObjectID    `json:"workspaceId" bson:"workspaceId"`
  Teams       []primitive.ObjectID  `json:"teams" bson:"teams"`
}

func (t *TimeDoctorReportForAnalysis) ToBytes() ([]byte, error) {
	return json.Marshal(t)
}

func (t *TimeDoctorReportForAnalysis) collection() *mongo.Collection {
  return storage.Primary().Collection("timedoctorReportsForAnalysis")
}

func (t *TimeDoctorReportForAnalysis) Save() error {
  ctx, cancel := utils.NewContext()
  defer cancel()
  _, err := t.collection().InsertOne(ctx, t)
  return err
}

type TimeDoctorReport struct {
	Id                primitive.ObjectID `json:"_id" bson:"_id"`
	ProductiveTime    float64            `json:"productiveTime" bson:"productiveTime"`
	UnproductiveTime  float64            `json:"unproductiveTime" bson:"unproductiveTime"`
	UncategorizedTime float64            `json:"uncategorizedTime" bson:"uncategorizedTime"`
  Summary           string             `json:"summary" bson:"summary"`
	ProductiveApps    []string           `json:"productiveApps" bson:"productiveApps"`
	Images            []struct {
		Date    time.Time `json:"date" bson:"date"`
		Summary string `json:"summary" bson:"summary"`
	}
  CreatedAt         time.Time          `json:"createdAt" bson:"createdAt"`
  Member            primitive.ObjectID `json:"member" bson:"member"`
  Workspace         primitive.ObjectID `json:"workspace" bson:"workspace"`
  Teams             []primitive.ObjectID `json:"teams" bson:"teams"`
}

func (t *TimeDoctorReport) collection() *mongo.Collection {
	return storage.Primary().Collection("timedoctorReports")
}

