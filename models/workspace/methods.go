package workspace

import (
	"github.com/snipextt/dayer/storage"
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (w *Workspace) collection() *mongo.Collection {
	return storage.GetMongoInstance().Collection("workspaces")
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

func (w *WorkspaceMember) collection() *mongo.Collection {
	return storage.GetMongoInstance().Collection("workspaceMembers")
}

func (w *WorkspaceMember) Save(update ...interface{}) (err error) {
	if w.Id.IsZero() {
		err = w.Create()
	} else {
		err = w.Update(update[0])
	}
	return
}

func (w *WorkspaceMember) Update(update interface{}) (err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	_, err = w.collection().UpdateByID(ctx, w.Id, update)

	return
}

func (w *WorkspaceMember) Create() (err error) {
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

func (t *WorkspaceTeam) collection() *mongo.Collection {
	return storage.GetMongoInstance().Collection("workspaceTeams")
}

func (t *WorkspaceTeam) Save(update ...any) (err error) {
	if t.Id.IsZero() {
		err = t.Create()
	} else {
		err = t.Update(update[0])
	}
	return
}

func (t *WorkspaceTeam) Create() (err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()

	res, err := t.collection().InsertOne(ctx, t)
	if err != nil {
		return
	}
	t.Id = res.InsertedID.(primitive.ObjectID)

	return
}

func (t *WorkspaceTeam) Update(update any) (err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()

	_, err = t.collection().UpdateByID(ctx, t.Id, bson.M{"$set": update})
	return
}
