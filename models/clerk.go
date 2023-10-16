package models

type ClerkUser struct {
	Id string `json:"id" bson:"_id"`
}

type ClerkWebhook struct {
	Type string    `json:"type" bson:"type"`
	Data ClerkUser `json:"data" bson:"data"`
}
