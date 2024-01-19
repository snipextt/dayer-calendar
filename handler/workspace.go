package handler

import (
	"encoding/json"
	"sync"
	"time"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/internal/timedoctor"
	"github.com/snipextt/dayer/models/connection"
	timedoctor_utils "github.com/snipextt/dayer/models/timedoctor"
	"github.com/snipextt/dayer/models/workspace"
	"github.com/snipextt/dayer/utils"
	clerk_internal "github.com/snipextt/dayer/utils/clerk"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCurrentWorkspace(c *fiber.Ctx) error {
	defer catchInternalServerError(c)

	wid, err := getWorkspaceId(c)
	if err != nil {
		return badRequest(c, err.Error())
	}
	isPersonal := c.Locals("personal").(bool)
	uoid, err := primitive.ObjectIDFromHex(c.Locals("uid").(string))
	utils.CheckError(err)
	w, err := workspace.GetWorkspaceAndConnections(wid, uoid, isPersonal)
	utils.CheckError(err)

	if w.Id.IsZero() {
		return success(c, nil, nil)
	}

	var res workspace.WorkspaceResponse

	res.RoleBasedResources = workspace.ResourcesForUser(w.User.Roles, w.User.Permissions)
	res.PendingConnections = workspace.PendingConnection(w.Extensions, w.Connections)
	res.Teams = w.Teams
	res.Id = w.Id

	return success(c, nil, res)
}

func CreateWorkspace(c *fiber.Ctx) error {
	defer catchInternalServerError(c)
	oid, oname, uid := c.Locals("oid").(string), c.Locals("oname").(string), c.Locals("uid").(string)
	userExtra := c.Locals("auth").(*clerk.TokenClaims).Extra
	email, name, image := userExtra["email"].(string), userExtra["user"].(string), userExtra["image"].(string)

	var extensions []string
	err := c.BodyParser(&extensions)
	utils.CheckError(err)

	ws := workspace.New(oname, oid, workspace.WorkspaceOrg, extensions)
	err = ws.Save()
	utils.CheckError(err)

	uoid, err := primitive.ObjectIDFromHex(uid)

	team := workspace.NewTeam("All members", "All members in the workspace", ws.Id, uoid)
	err = team.Save()
	utils.CheckError(err)

	ws.DefaultTeam = team.Id
	err = ws.Save(bson.M{"defaultTeam": team.Id})
	utils.CheckError(err)

	meta := workspace.MemberMeta{
		Source: "backend",
	}

	member := workspace.NewMember(name, email, image, ws.Id, team.Id, meta, workspace.WorkspaceRoleAdmin)
	member.User = uoid
	utils.CheckError(err)

	err = member.Save()
	utils.CheckError(err)

	var res workspace.WorkspaceResponse

	res.RoleBasedResources = workspace.ResourcesForUser(member.Roles, member.Permissions)
	res.PendingConnections = workspace.PendingConnection(ws.Extensions, []connection.Model{})
	res.Teams = []workspace.Team{*team}
	res.Id = ws.Id

	metaUpdate, err := json.Marshal(map[string]any{
		"workspaceId": ws.Id.Hex(),
	})
	utils.CheckError(err)

	clerk_internal.ClerkClient().Organizations().UpdateMetadata(oid, clerk.UpdateOrganizationMetadataParams{
		PublicMetadata: metaUpdate,
	})

	return success(c, nil, res)
}

func ConnectTimeDoctor(c *fiber.Ctx) error {
	defer catchInternalServerError(c)

	wid, err := getWorkspaceId(c)
	if err != nil {
		return badRequest(c, err.Error())
	}

	body := make(map[string]interface{})
	err = c.BodyParser(&body)
	utils.CheckError(err)
	res, err := timedoctor.Login(body["email"].(string), body["password"].(string), body["totp"].(string))
	utils.CheckError(err)

	if res.Error != "" {
		return badRequest(c, res.Message)
	}

	conn := connection.NewTimeDoctorConnection(wid, body["email"].(string), res.Data.Token, res.Data.ExpiresAt)
	conn.Save()

	return success(c, nil, res.Data.Companies)
}

func ConnectTimeDoctorCompany(c *fiber.Ctx) (err error) {
	body := make(map[string]interface{})
	err = c.BodyParser(&body)
	utils.CheckError(err)

	wid, err := getWorkspaceId(c)
	if err != nil {
		return badRequest(c, err.Error())
	}

	conn, err := connection.ForWorkspaceByProvider(wid, "timedoctor")
	utils.CheckError(err)

	ws, err := workspace.FindWorkspace(wid)
	utils.CheckError(err)

	meta := connection.WorkspaceMeta{
		TimeDoctorCompanyID:       body["company"].(string),
		TimeDoctorParseScreencast: body["parseScreencast"].(bool),
	}

	conn.Meta = meta
	conn.Save()

	users, err := timedoctor.GetCompanyUsers(conn.Token, body["company"].(string))
	utils.CheckError(err)

	var wg sync.WaitGroup

	for _, user := range users.Data {
		wg.Add(1)
		go func(user timedoctor.TimeDoctorUser) {
			defer wg.Done()
			meta := workspace.MemberMeta{
				Source:       "timedoctor",
				TimeDoctorId: user.Id,
			}
			u := workspace.NewMember(user.Name, user.Email, "", wid, ws.DefaultTeam, meta, workspace.WorkspaceRoleMember)
			u.Save()
		}(user)
	}

	return
}

func GetTeam(c *fiber.Ctx) error {
	defer catchInternalServerError(c)
	wid, err := getWorkspaceId(c)
	teamIdHex := c.Params("id")
	tid, err := primitive.ObjectIDFromHex(teamIdHex)
	utils.CheckError(err)
	if err != nil {
		return badRequest(c, err.Error())
	}

	members, err := workspace.FindTeamMembers(tid, wid)
	utils.CheckError(err)

	return success(c, nil, members)
}

func CreateTeam(c *fiber.Ctx) error {
	defer catchInternalServerError(c)
	wid, err := getWorkspaceId(c)
	if err != nil {
		return badRequest(c, err.Error())
	}
	body := make(map[string]string)
	err = c.BodyParser(&body)

	uid, err := primitive.ObjectIDFromHex(c.Locals("uid").(string))
	utils.CheckError(err)

	if body["name"] == "" {
		return badRequest(c, "Name is required")
	}

	if body["description"] == "" {
		return badRequest(c, "Description is required")
	}

	team := workspace.NewTeam(body["name"], body["description"], wid, uid)
	err = team.Save()
	utils.CheckError(err)

	member, err := workspace.FindWorkspaceMember(wid, uid)
	utils.CheckError(err)

	err = member.Save(bson.M{"$push": bson.M{"teams": team.Id}})
	utils.CheckError(err)

	return success(c, nil, team)
}

func Reports(c *fiber.Ctx) error {
  defer catchInternalServerError(c)

  wid, err := getWorkspaceId(c)
  if err != nil {
    return badRequest(c, err.Error())
  }
  startDate, err := time.Parse(time.RFC3339, c.Query("startDate"))
  endDate, err := time.Parse(time.RFC3339, c.Query("endDate"))
  team, err := primitive.ObjectIDFromHex(c.Query("team"))

  connections, err := connection.ForWorkspace(wid)
  utils.CheckError(err)

  for _, conn := range connections {
    if conn.Provider == "timedoctor" {
      res, err := timedoctor_utils.ReportForWorkspace(conn.Workspace.(primitive.ObjectID), team, startDate, endDate)
      utils.CheckError(err)
      return success(c, nil, res)
    }
  }

  return nil
}
