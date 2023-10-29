package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/internal/timedoctor"
	"github.com/snipextt/dayer/models/connection"
	"github.com/snipextt/dayer/models/workspace"
	"github.com/snipextt/dayer/utils"
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

	var extensions []string
	err := c.BodyParser(&extensions)
	utils.CheckError(err)

	w := workspace.New(oname, oid, workspace.WorkspaceOrg, extensions)
	err = w.Save()

	member := workspace.NewWorkspaceMember(w.Id.Hex(), uid, workspace.WorkspaceRoleAdmin)
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

	return nil
}
