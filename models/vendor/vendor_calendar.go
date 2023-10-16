package vendor

type VendorCalendar struct {
	Id             string `json:"id" bson:"_id,omitempty"`
	Vendor         string `json:"-" bson:"vendor"`
	VendorId       string `json:"-" bson:"vendorId"`
	LinkedCalendar string `json:"linkedCalendar" bson:"linkedCalendar"`
	AccessRole     string `json:"accessRole" bson:"accessRole"`
	SyncedOn       int64  `json:"syncedOn" bson:"syncedOn"`
	Connection     string `json:"connection" bson:"connection"`
}
