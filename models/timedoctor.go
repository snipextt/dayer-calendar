package models

import (
	"encoding/json"
	"time"

	"github.com/snipextt/dayer/storage"
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
	Images      []TimeDoctorImageData `json:"images" bson:"images"`
	Activities  []TimeDoctorActivity  `json:"activities" bson:"activities"`
	Tasks       []string              `json:"tasks" bson:"tasks"`
	CreatedAt   time.Time             `json:"createdAt" bson:"createdAt"`
	MemberId    primitive.ObjectID    `json:"memberId" bson:"memberId"`
	WorkspaceId primitive.ObjectID    `json:"workspaceId" bson:"workspaceId"`
}

func (t *TimeDoctorReportForAnalysis) ToBytes() ([]byte, error) {
	return json.Marshal(t)
}

type TimeDoctorReport struct {
	Id                primitive.ObjectID `json:"_id" bson:"_id"`
	ProductiveTime    float64            `json:"productiveTime" bson:"productiveTime"`
	UnproductiveTime  float64            `json:"unproductiveTime" bson:"unproductiveTime"`
	UncategorizedTime float64            `json:"uncategorizedTime" bson:"uncategorizedTime"`
	ProductiveApps    []string           `json:"productiveApps" bson:"productiveApps"`
	Images            []struct {
		Date    string `json:"date" bson:"date"`
		Summary string `json:"summary" bson:"summary"`
	}
}

func (t *TimeDoctorReport) collection() *mongo.Collection {
	return storage.Primary().Collection("timedoctorReports")
}
