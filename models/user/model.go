package user

import (
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
)

type ModifiedRecord struct {
	Native string
	Bson   string
}

type Model struct {
	Id              primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	ClerkId         string             `json:"clerkId" bson:"clerkId"`
	Active          bool               `json:"active" bson:"active"`
	modifiedRecords []ModifiedRecord
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

func (u *Model) Save() (err error) {
	if u.Id.IsZero() {
		err = u.Create()
	} else {
		err = u.Update()
	}
	return
}

func (u *Model) Update() (err error) {
	ctx, cancel := utils.GetContext()
	defer cancel()
	user := reflect.ValueOf(u).Elem()
	update := make(map[string]interface{})
	for _, v := range u.modifiedRecords {
		update[v.Bson] = user.FieldByName(v.Native).Interface()
	}
	_, err = collection().UpdateOne(ctx, bson.M{"_id": u.Id}, bson.M{"$set": update})
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
