package models

type Calenar struct {
	Id  string `json:"id" bson:"_id,omitempty"`
	Uid string `json:"uid" bson:"uid"`
}
