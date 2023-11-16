package calendar

import (
	"github.com/snipextt/dayer/cmd"
	"github.com/snipextt/dayer/storage"
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type VendorCalendarConnection struct {
	Id               primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ChannelId        string             `json:"channelId" bson:"channelId"`
	VendorCalendarId string             `json:"googleCalendarId" bson:"googleCalendarId"`
	Vendor           string             `json:"vendor" bson:"vendor"`
	Token            string             `json:"token" bson:"token"`
}

func NewGoogleCalendarConnection(id string) *VendorCalendarConnection {
	token := cmd.GetRandomString(64)
	return &VendorCalendarConnection{
		VendorCalendarId: id,
		Token:            token,
		Vendor:           "google",
	}
}

func (c *VendorCalendarConnection) collection() *mongo.Collection {
	return storage.Primary().Collection("vendorCalendars")
}

func (c *VendorCalendarConnection) Save() (err error) {
	if c.Id.IsZero() {
		err = c.Create()
	} else {
		err = c.Update()
	}
	return
}

func (c *VendorCalendarConnection) Update() (err error) {
	return nil
}

func (c *VendorCalendarConnection) Create() (err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	res, err := c.collection().InsertOne(ctx, c)
	c.Id = res.InsertedID.(primitive.ObjectID)
	return
}
