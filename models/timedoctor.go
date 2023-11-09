package models

import (
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

type TimeDoctorReport struct {
	Id         primitive.ObjectID    `json:"_id" bson:"_id,omitempty"`
	Images     []TimeDoctorImageData `json:"images" bson:"images"`
	Activities []TimeDoctorActivity  `json:"activities" bson:"activities"`
	Tasks      []string              `json:"tasks" bson:"tasks"`
	CreatedAt  time.Time             `json:"createdAt" bson:"createdAt"`
}

func (*TimeDoctorReport) Collection() *mongo.Collection {
	return storage.GetMongoInstance().Collection("timedoctorReports")
}

func (t *TimeDoctorReport) Save(update ...any) (err error) {
	if t.Id.IsZero() {
		err = t.Create()
	} else {
		err = t.Update(update...)
	}
	return
}

func (t *TimeDoctorReport) Create() (err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	res, err := t.Collection().InsertOne(ctx, t)
	if err != nil {
		t.Id = res.InsertedID.(primitive.ObjectID)
	}
	return
}

func (t *TimeDoctorReport) Update(update ...any) (err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	_, err = t.Collection().UpdateByID(ctx, t.Id, update)
	return
}
