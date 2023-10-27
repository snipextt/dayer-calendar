package user

import (
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Model struct {
	Id         primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ClerkId    string             `json:"clerkId" bson:"clerkId"`
	Active     bool               `json:"active" bson:"active"`
	Workspaces []interface{}      `json:"workspaces" bson:"workspaces"`
}

func (u *Model) FindById(uid string) error {
	ctx, cancel := utils.GetContext()
	defer cancel()
	id, err := primitive.ObjectIDFromHex(uid)
	if err != nil {
		return err
	}
	err = collection().FindOne(ctx, bson.M{"_id": id}).Decode(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *Model) FindByClerkId(clerkId string) error {
	ctx, cancel := utils.GetContext()
	defer cancel()
	err := collection().FindOne(ctx, bson.M{"clerkId": clerkId}).Decode(u)
	if err != nil {
		return err
	}
	return nil
}

func (u *Model) Save(update ...interface{}) (err error) {
	if u.Id.IsZero() {
		err = u.Create()
	} else {
		err = u.Update(update[0])
	}
	return
}

func (u *Model) Update(update interface{}) (err error) {
	ctx, cancel := utils.GetContext()
	defer cancel()
	_, err = collection().UpdateByID(ctx, u.Id, bson.M{"$set": update})
	return
}

func (u *Model) Create() error {
	ctx, cancel := utils.GetContext()
	defer cancel()
	r, err := collection().InsertOne(ctx, u)
	if err != nil {
		return err
	}
	u.Id = r.InsertedID.(primitive.ObjectID)
	if err != nil {
		return err
	}
	return nil
}
