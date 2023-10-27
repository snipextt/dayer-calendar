package workspace

import (
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func New(name string, clerkOrgId string, extensions []interface{}) *Workspace {
	return &Workspace{
		Name:       name,
		Type:       WorkspaceOrg,
		ClerkOrgId: clerkOrgId,
		Extensions: extensions,
	}
}

func NewWorkspaceMember(oid string, uid string, roles ...string) *WorkspaceMember {
	return &WorkspaceMember{
		UserId:      uid,
		WorkspaceId: oid,
		Roles:       roles,
		Permissions: []string{},
	}
}

func FindByClerkId(id string) (workspace Workspace, err error) {
	ctx, cancel := utils.GetContext()
	defer cancel()
	err = workspace.collection().FindOne(ctx, Workspace{ClerkOrgId: id}).Decode(&workspace)
	if err == mongo.ErrNoDocuments {
		err = nil
	}
	return
}

func FindWorkspaceAndExtensions(orgId string) (workspace Workspace, err error) {
	ctx, cancel := utils.GetContext()
	defer cancel()
	match := bson.D{{Key: "$match", Value: bson.D{{Key: "clerkOrgId", Value: orgId}}}}
	convertId := bson.D{{
		Key: "$project", Value: bson.D{
			{Key: "extensions", Value: bson.D{{
				Key: "$map",
				Value: bson.D{
					{Key: "input", Value: "$extensions"},
					{Key: "as", Value: "extension"},
					{Key: "in", Value: bson.D{
						{Key: "$toObjectId", Value: "$$extension"},
					}},
				}},
			}},
		},
	}}
	lookupextensions := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "extensions"},
		{Key: "localField", Value: "extensions"},
		{Key: "foreignField", Value: "_id"},
		{Key: "as", Value: "extensions"},
	}}}
	lookupconnections := bson.D{{Key: "$lookup", Value: bson.D{
		{Key: "from", Value: "connections"},
		{Key: "localField", Value: "_id"},
		{Key: "foreignField", Value: "workspaceId"},
		{Key: "as", Value: "connections"},
	}}}

	res, err := workspace.collection().Aggregate(ctx, mongo.Pipeline{match, convertId, lookupextensions, lookupconnections})
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

func FindWorkspaceMember(orgId string, uid string) (member WorkspaceMember, err error) {
	ctx, cancel := utils.GetContext()
	defer cancel()
	filter := bson.D{{Key: "workspaceId", Value: orgId}, {Key: "userId", Value: uid}}
	err = member.collection().FindOne(ctx, filter).Decode(&member)
	return
}

func GetResourcesForUser(roles []string, permissions []string) (resources []string) {
	rolesmap := make(map[string]bool)
	for _, r := range roles {
		for _, p := range GetResourcesForRole(r) {
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

func GetResourcesForRole(r string) []string {
	switch r {
	case WorkspaceRoleAdmin:
		return PermissionsAdminOrg
	default:
		return []string{}
	}
}
