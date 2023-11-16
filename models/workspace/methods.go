package workspace

import (
	"github.com/snipextt/dayer/storage"
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (w *Workspace) collection() *mongo.Collection {
	return storage.Primary().Collection("workspaces")
}

func (w *Workspace) Save(update ...interface{}) (err error) {
	if w.Id.IsZero() {
		err = w.Create()
	} else {
		err = w.Update(update[0])
	}
	return
}

func (w *Workspace) Update(update interface{}) (err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	_, err = w.collection().UpdateByID(ctx, w.Id, bson.M{"$set": update})

	return
}

func (w *Workspace) Create() (err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()

	res, err := w.collection().InsertOne(ctx, w)
	if err != nil {
		return
	}
	w.Id = res.InsertedID.(primitive.ObjectID)
	return
}

// Workspace member methods

func (w *Member) collection() *mongo.Collection {
	return storage.Primary().Collection("workspaceMembers")
}

func (w *Member) Save(update ...interface{}) (err error) {
	if w.Id.IsZero() {
		err = w.Create()
	} else {
		err = w.Update(update[0])
	}
	return
}

func (w *Member) Update(update interface{}) (err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	_, err = w.collection().UpdateByID(ctx, w.Id, update)

	return
}

func (w *Member) Create() (err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()

	res, err := w.collection().InsertOne(ctx, w)
	if err != nil {
		return
	}
	w.Id = res.InsertedID.(primitive.ObjectID)
	return
}

// Workspace Team methods

func (t *Team) collection() *mongo.Collection {
	return storage.Primary().Collection("workspaceTeams")
}

func (t *Team) Save(update ...any) (err error) {
	if t.Id.IsZero() {
		err = t.Create()
	} else {
		err = t.Update(update[0])
	}
	return
}

func (t *Team) Create() (err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()

	res, err := t.collection().InsertOne(ctx, t)
	if err != nil {
		return
	}
	t.Id = res.InsertedID.(primitive.ObjectID)

	return
}

func (t *Team) Update(update any) (err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()

	_, err = t.collection().UpdateByID(ctx, t.Id, bson.M{"$set": update})
	return
}
