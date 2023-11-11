package handler

import (
	"sync"

	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/internal/timedoctor"
	"github.com/snipextt/dayer/models/connection"
	"github.com/snipextt/dayer/models/workspace"
	"github.com/snipextt/dayer/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCurrentWorkspace(c *fiber.Ctx) error {
	defer HandleInternalServerError(c)

	clerkoid := c.Locals("oid").(string)
	uid := c.Locals("uid").(string)
	personal := c.Locals("personal").(bool)
	w, err := workspace.FindWorkspaceAndExtensions(clerkoid)
	utils.CheckError(err)

	if w.Id.IsZero() {
		return HandleSuccess(c, nil, nil)
	}

	var u workspace.WorkspaceMember

	if personal {
		u = workspace.WorkspaceMember{
			Roles: []string{workspace.WorkspaceRoleAdmin},
		}
	} else {
		u, err = workspace.FindWorkspaceMember(w.Id.Hex(), uid)
		utils.CheckError(err)
	}

	var res workspace.WorkspaceResponse

	res.RoleBasedResources = workspace.GetResourcesForUser(u.Roles, u.Permissions)
	res.PendingConnections = workspace.GetPendingConnection(w.Extensions, w.Connections)
	res.Id = w.Id

	return HandleSuccess(c, nil, res)
}

func CreateWorkspace(c *fiber.Ctx) error {
	defer HandleInternalServerError(c)
	oid, oname, uid := c.Locals("oid").(string), c.Locals("oname").(string), c.Locals("uid").(string)
	userExtra := c.Locals("auth").(*clerk.TokenClaims).Extra
	email, name, image := userExtra["email"].(string), userExtra["user"].(string), userExtra["image"].(string)

	var extensions []string
	err := c.BodyParser(&extensions)
	utils.CheckError(err)

	w := workspace.New(oname, oid, workspace.WorkspaceOrg, extensions)
	err = w.Save()

	meta := workspace.WorkspaceMemberMeta{
		Source: "backend",
	}

	member := workspace.NewWorkspaceMember(name, w.Id, email, image, meta, workspace.WorkspaceRoleAdmin)
	member.User = uid
	err = member.Save()
	utils.CheckError(err)

	var res workspace.WorkspaceResponse

	res.RoleBasedResources = workspace.GetResourcesForUser(member.Roles, member.Permissions)
	res.PendingConnections = workspace.GetPendingConnection(w.Extensions, w.Connections)
	res.Id = w.Id

	return HandleSuccess(c, nil, res)
}

func ConnectTimeDoctor(c *fiber.Ctx) error {
	defer HandleInternalServerError(c)

	wid, ok := c.GetReqHeaders()["X-Workspace-Id"]

	if !ok {
		return HandleBadRequest(c, "Workspace Id not set")
	}

	body := make(map[string]interface{})
	err := c.BodyParser(&body)
	utils.CheckError(err)
	res, err := timedoctor.Login(body["email"].(string), body["password"].(string), body["totp"].(string))
	utils.CheckError(err)

	if res.Error != "" {
		return HandleBadRequest(c, res.Message)
	}

	conn := connection.NewTimeDoctorConnection(wid, body["email"].(string), res.Data.Token, res.Data.ExpiresAt)
	conn.Save()

	return HandleSuccess(c, nil, res.Data.Companies)
}

func ConnectTimeDoctorCompany(c *fiber.Ctx) (err error) {
	body := make(map[string]interface{})
	err = c.BodyParser(&body)
	utils.CheckError(err)

	wid, ok := c.GetReqHeaders()["X-Workspace-Id"]
	if !ok {
		return HandleBadRequest(c, "Workspace Id not set")
	}

  workspaceOid, err := primitive.ObjectIDFromHex(wid)
	utils.CheckError(err)

	conn, err := connection.FindByWorkspaceId(wid, "timedoctor")
	utils.CheckError(err)

	meta := connection.WorkspaceMeta{
		TimeDoctorCompanyID: body["company"].(string),
	}

	conn.Meta = meta
	conn.Save()

	users, err := timedoctor.GetCompanyUsers(conn.Token, body["company"].(string))
	utils.CheckError(err)

	for _, user := range users.Data {
		meta := workspace.WorkspaceMemberMeta{
			Source:       "timedoctor",
			TimeDoctorId: user.Id,
		}
		u := workspace.NewWorkspaceMember(user.Name, workspaceOid, user.Email, "", meta, workspace.WorkspaceRoleMember)
		u.Save()
	}

	return
}

func GetDataFromTimeDoctor(c *fiber.Ctx) error {
	defer HandleInternalServerError(c)

	wid, ok := c.GetReqHeaders()["X-Workspace-Id"]
	if !ok {
		return HandleBadRequest(c, "Workspace Id not set")
	}

	conn, err := connection.FindByWorkspaceId(wid, "timedoctor")
	utils.CheckError(err)

	users, err := workspace.FindWorkspaceMembers(wid)
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

func GetMembers(c *fiber.Ctx) error {
	defer HandleInternalServerError(c)
	wid, ok := c.GetReqHeaders()["X-Workspace-Id"]
	if !ok {
		return HandleBadRequest(c, "Workspace Id not set")
	}

	members, err := workspace.FindWorkspaceMembers(wid)
	utils.CheckError(err)

	return HandleSuccess(c, nil, members)
}

func GetPeers(c *fiber.Ctx) error {
  defer HandleInternalServerError(c)
	wid, ok := c.GetReqHeaders()["X-Workspace-Id"]
	if !ok {
		return HandleBadRequest(c, "Workspace Id not set")
	}
  return nil
}
