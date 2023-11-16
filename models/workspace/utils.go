package workspace

import (
	"github.com/snipextt/dayer/models/connection"
	"github.com/snipextt/dayer/storage"
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func collection() *mongo.Collection {
	return storage.Primary().Collection("workspaceMembers")
}

func New(name string, clerkOrgId string, wsType string, extensions []string) *Workspace {
	return &Workspace{
		Name:       name,
		Type:       wsType,
		ClerkOrgId: clerkOrgId,
		Extensions: extensions,
	}
}

func NewMember(name, email, image string, wid, team primitive.ObjectID, meta MemberMeta, roles ...string) *Member {
	return &Member{
		Name:        name,
		Image:       image,
		Workspace:   wid,
		Email:       email,
		Teams:       bson.A{team},
		Roles:       roles,
		Permissions: []string{},
		Meta:        meta,
	}
}

func NewTeam(name, description string, wid, owner primitive.ObjectID) *Team {
	return &Team{
		Name:        name,
		Workspace:   wid,
		Owner:       owner,
		Description: description,
	}
}

func FindByClerkId(id string) (workspace Workspace, err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	err = workspace.collection().FindOne(ctx, Workspace{ClerkOrgId: id}).Decode(&workspace)
	if err == mongo.ErrNoDocuments {
		err = nil
	}
	return
}

func GetWorkspaceAndConnections(id, uid primitive.ObjectID, personal bool) (workspace WorkspaceAggregation, err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	match := bson.D{{Key: "$match", Value: bson.D{{Key: "_id", Value: id}}}}
	lookupwsuser := bson.D{{Key: "$lookup", Value: bson.M{
		"from": "workspaceMembers",
		"pipeline": mongo.Pipeline{
			bson.D{{Key: "$match", Value: bson.M{
				"$expr": bson.M{
					"$and": bson.A{
						bson.M{"$eq": bson.A{"$workspace", id}},
						bson.M{"$eq": bson.A{"$user", uid}},
					},
				},
			}}},
		},
		"as": "user",
	}}}
	unwinduser := bson.D{{Key: "$unwind", Value: "$user"}}
	lookupteams := bson.D{{Key: "$lookup", Value: bson.M{
		"from":         "workspaceTeams",
		"localField":   "user.teams",
		"foreignField": "_id",
		"as":           "teams",
	}}}
	lookupconnections := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "connections"},
		{Key: "localField", Value: "_id"},
		{Key: "foreignField", Value: "workspace"},
		{Key: "as", Value: "connections"},
	}}}

	var pipeline mongo.Pipeline
	if personal {
		pipeline = append(pipeline, match, lookupconnections)
	} else {
		pipeline = append(pipeline, match, lookupwsuser, unwinduser, lookupteams, lookupconnections)
	}

	res, err := storage.Primary().Collection("workspaces").Aggregate(ctx, pipeline)
	if err != nil {
		return
	}
	defer res.Close(ctx)
	matched := res.Next(ctx)
	if !matched {
		return
	}
	err = res.Decode(&workspace)

	return
}

func FindWorkspace(id primitive.ObjectID) (workspace Workspace, err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()

	res := workspace.collection().FindOne(ctx, bson.M{"_id": id})
	if res.Err() != nil {
		err = res.Err()
		return
	}
	err = res.Decode(&workspace)
	return
}

func FindWorkspaceMember(wid, uid primitive.ObjectID) (member Member, err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	filter := bson.D{{Key: "workspace", Value: wid}, {Key: "user", Value: uid}}
	err = member.collection().FindOne(ctx, filter).Decode(&member)
	return
}

func FindWorkspaceMembers(wid primitive.ObjectID) (members []Member, err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()
	filter := bson.D{{Key: "workspace", Value: wid}}
	cur, err := collection().Find(ctx, filter)
	if err != nil {
		return
	}
	err = cur.All(ctx, &members)
	return
}

func FindTeamMembers(team, ws primitive.ObjectID) (members []Member, err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()

	filter := bson.M{"teams": bson.M{"$in": bson.A{team}}, "workspace": ws}

	res, err := collection().Find(ctx, filter)
	if err != nil {
		return
	}
	err = res.All(ctx, &members)
	return
}

func GetMangedMembers(managerId string) (members []Member, err error) {
	ctx, cancel := utils.NewContext()
	defer cancel()

	filter := bson.M{"manager": managerId}

	res, err := collection().Find(ctx, filter)
	if err != nil {
		return
	}
	err = res.All(ctx, &members)
	return
}

func ResourcesForUser(roles []string, permissions []string) (resources []string) {
	rolesmap := make(map[string]bool)
	for _, r := range roles {
		for _, p := range ResourcesForRole(r) {
			rolesmap[p] = true
		}
	}
	for _, p := range permissions {
		rolesmap[p] = true
	}

	for k := range rolesmap {
		resources = append(resources, k)
	}
	return
}

func ResourcesForRole(r string) []string {
	switch r {
	case WorkspaceRoleAdmin:
		return PermissionsAdminOrg
	default:
		return []string{}
	}
}

func PendingConnection(extensions []string, wsconnections []connection.Model) []WorkspaceEvent {
	var pending []WorkspaceEvent
	for _, e := range extensions {
		var found bool
		for _, c := range wsconnections {
			if e == c.Provider {
				if c.Provider == "timedoctor" && c.Meta.TimeDoctorCompanyID == "" {
					continue
				}
				found = true
				break
			}
		}
		if !found {
			pending = append(pending, WorkspaceEvent{Type: "client", Name: e})
		}
	}
	return pending
}
