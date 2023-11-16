package vendor

import (
	"github.com/snipextt/dayer/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

func NewCalendar(vendor string, id string, role string, connection string) *VendorCalendar {
	return &VendorCalendar{
		Vendor:     vendor,
		VendorId:   id,
		AccessRole: role,
		Connection: connection,
	}
}

func collection() *mongo.Collection {
	return storage.Primary().Collection("vendorCalendar")
}
