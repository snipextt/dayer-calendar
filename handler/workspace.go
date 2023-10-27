package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/models/workspace"
	"github.com/snipextt/dayer/utils"
)

func GetCurrentWorkspace(ctx *fiber.Ctx) error {
	defer HandleInternalServerError(ctx)

	clerkoid := ctx.Locals("oid").(string)
	uid := ctx.Locals("uid").(string)
	w, err := workspace.FindWorkspaceAndExtensions(clerkoid)

	utils.PanicOnError(err)

	if w.Id.IsZero() {
		return HandleSuccess(ctx, nil, nil)
	}

	u, err := workspace.FindWorkspaceMember(w.Id.Hex(), uid)
	utils.PanicOnError(err)

	var res workspace.WorkspaceResponse

	res.Workspace = w
	res.RoleBasedResources = workspace.GetResourcesForUser(u.Roles, u.Permissions)
	res.Workspace.Extensions = w.Extensions

	return HandleSuccess(ctx, nil, res)
}

func CreateWorkspace(ctx *fiber.Ctx) error {
	defer HandleInternalServerError(ctx)
	oid, oname, uid := ctx.Locals("oid").(string), ctx.Locals("oname").(string), ctx.Locals("uid").(string)

	var extensions []interface{}
	err := ctx.BodyParser(&extensions)
	utils.PanicOnError(err)

	w := workspace.New(oname, oid, extensions)
	err = w.Save()

	user := workspace.NewWorkspaceMember(w.Id.Hex(), uid, workspace.WorkspaceRoleAdmin)
	err = user.Save()

	utils.PanicOnError(err)
	return HandleSuccess(ctx, nil, w)
}
