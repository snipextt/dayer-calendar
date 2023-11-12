package handler

import (
	"sync"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/internal/timedoctor"
	"github.com/snipextt/dayer/models/connection"
	"github.com/snipextt/dayer/models/workspace"
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCurrentWorkspace(c *fiber.Ctx) error {
	defer catchInternalServerError(c)

	wid, err := getWorkspaceId(c)
	if err != nil {
		return badRequest(c, err.Error())
	}
	uoid, err := primitive.ObjectIDFromHex(c.Locals("uid").(string))
	utils.CheckError(err)
	personal := c.Locals("personal").(bool)
	w, err := workspace.GetWorkspaceAndConnections(wid, uoid)
	utils.CheckError(err)

	if w.Id.IsZero() {
		return success(c, nil, nil)
	}

	var u workspace.WorkspaceMember
	if personal {
		u = workspace.WorkspaceMember{
			Roles: []string{workspace.WorkspaceRoleAdmin},
		}
	} else {
		u, err = workspace.FindWorkspaceMember(w.Id.Hex(), uoid)
		utils.CheckError(err)
	}

	var res workspace.WorkspaceResponse

	res.RoleBasedResources = workspace.GetResourcesForUser(u.Roles, u.Permissions)
	res.PendingConnections = workspace.GetPendingConnection(w.Extensions, w.Connections)
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

	team := workspace.NewTeam("All members", "All members in the workspace", ws.Id)
	err = team.Save()
	utils.CheckError(err)

	ws.DefaultTeam = team.Id
	err = ws.Save(bson.M{"defaultTeam": team.Id})
	utils.CheckError(err)

	meta := workspace.WorkspaceMemberMeta{
		Source: "backend",
	}

	member := workspace.NewMember(name, email, image, ws.Id, team.Id, meta, workspace.WorkspaceRoleAdmin)
	member.User = uid
	err = member.Save()
	utils.CheckError(err)

	var res workspace.WorkspaceResponse

	res.RoleBasedResources = workspace.GetResourcesForUser(member.Roles, member.Permissions)
	res.PendingConnections = workspace.GetPendingConnection(ws.Extensions, []connection.Model{})
	res.Id = ws.Id

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

	conn, err := connection.FindByWorkspaceId(wid, "timedoctor")
	utils.CheckError(err)

	ws, err := workspace.FindWorkspace(wid)
	utils.CheckError(err)

	meta := connection.WorkspaceMeta{
		TimeDoctorCompanyID: body["company"].(string),
	}

	conn.Meta = meta
	conn.Save()

	users, err := timedoctor.GetCompanyUsers(conn.Token, body["company"].(string))
	utils.CheckError(err)

	var wg sync.WaitGroup

	for _, user := range users.Data {
		wg.Add(1)
		go func() {
			defer wg.Done()
			meta := workspace.WorkspaceMemberMeta{
				Source:       "timedoctor",
				TimeDoctorId: user.Id,
			}
			u := workspace.NewMember(user.Name, user.Email, "", wid, ws.DefaultTeam, meta, workspace.WorkspaceRoleMember)
			u.Save()
		}()
	}

	return
}

func GetDataFromTimeDoctor(c *fiber.Ctx) error {
	defer catchInternalServerError(c)

	var workspaceId primitive.ObjectID
	var err error

	if wid, ok := c.GetReqHeaders()["X-Workspace-Id"]; !ok {
		return badRequest(c, "Workspace Id not set")
	} else {
		workspaceId, err = primitive.ObjectIDFromHex(wid)
		utils.CheckError(err)
	}

	conn, err := connection.FindByWorkspaceId(workspaceId, "timedoctor")
	utils.CheckError(err)

	users, err := workspace.FindWorkspaceMembers(workspaceId)
	utils.CheckError(err)

	var wg sync.WaitGroup
	for _, user := range users {
		if user.Meta.Source == "timedoctor" {
			wg.Add(1)
			go func(uid string) {
				defer wg.Done()
				report, err := timedoctor.GenerateReportFromTimedoctor(conn.Token, conn.Meta.TimeDoctorCompanyID, uid)
				utils.CheckError(err)
				report.Save()
			}(user.Meta.TimeDoctorId)
		}
	}
	wg.Wait()
	return nil
}

func GetTeams(c *fiber.Ctx) error {
	// defer
	return nil
}

func GetTeam(c *fiber.Ctx) error {
	defer catchInternalServerError(c)
	wid, err := getWorkspaceId(c)
	teamIdHex := c.Params("id")
	teamOid, err := primitive.ObjectIDFromHex(teamIdHex)
	utils.CheckError(err)
	if err != nil {
		return badRequest(c, err.Error())
	}

	members, err := workspace.FindTeamMembers(teamOid, wid)
	utils.CheckError(err)

	return success(c, nil, members)
}
