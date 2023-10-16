package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventType int

const (
	Work EventType = iota
	Leisure
	Personal
)

type Event struct {
	Id             primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	UserId         primitive.ObjectID `json:"userId" bson:"userId"`
	Title          string             `json:"title" bson:"title"`
	EventType      EventType          `json:"eventType" bson:"eventType"`
	Description    string             `json:"description" bson:"description"`
	StartTime      time.Time          `json:"startTime" bson:"startTime"`
	EndTime        time.Time          `json:"endTime" bson:"endTime"`
	EstimatedTime  int                `json:"estimatedTime" bson:"estimatedTime"`
	OrignalEventId string             `json:"originalEventId" bson:"originalEventId"`
}
