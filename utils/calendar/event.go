package calendar

import (
	"github.com/snipextt/dayer/storage"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Event struct {
	Id                     string `json:"id" bson:"_id,omitempty"`
	Title                  string `json:"title" bson:"title"`
	Description            string `json:"description" bson:"description"`
	Start                  string `json:"start" bson:"start"`
	End                    string `json:"end" bson:"end"`
	Location               string `json:"location" bson:"location"`
	Link                   string `json:"link" bson:"link"`
	Vendor                 string `json:"vendor" bson:"vendor"`
	CalendarConnectionId   string `json:"calendarConnectionId" bson:"calendarConnectionId"`
	CalendarConnectionType string `json:"calendarConnectionType" bson:"calendarConnectionType"`
	IsRecurring            bool   `json:"isRecurring" bson:"isRecurring"`
	RecurringId            string `json:"recurringId" bson:"recurringId"`
	RecurrenceRule         `json:"recurrenceRule" bson:"recurrenceRule"`
}

type RecurrenceRule struct {
	Frequency string   `json:"frequency" bson:"frequency"`
	Days      []string `json:"days" bson:"days"`
	Until     string   `json:"until" bson:"until"`
	Count     int      `json:"count" bson:"count"`
}

func (e *Event) collection() *mongo.Collection {
	return storage.Primary().Collection("events")
}

func (e *Event) Save() (err error) {
	if e.Id == "" {
		err = e.Create()
	} else {
		err = e.Update()
	}
	return
}

func (e *Event) Update() (err error) {
	return
}

func (e *Event) Create() (err error) {
	r, err := e.collection().InsertOne(nil, e)
	if err != nil {
		return
	}
	e.Id = r.InsertedID.(primitive.ObjectID).Hex()
	return
}
