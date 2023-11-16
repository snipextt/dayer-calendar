package handler

import (
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/models/user"
	"github.com/snipextt/dayer/models/workspace"
	"github.com/snipextt/dayer/utils"
	"github.com/snipextt/dayer/utils/clerk"
)

func Onboarding(c *fiber.Ctx) error {
	catchInternalServerError(c)

	cuid := c.Locals("auth").(*clerk.TokenClaims).Claims.Subject
	oid, oname := c.Locals("oid").(string), c.Locals("oname").(string)
	cu, err := clerk_utils.GetClerkUser(cuid)
	utils.CheckError(err)

	if cu.ExternalID != nil {
		return badRequest(c, "User already onboarded")
	}

	u := user.New(cuid)
	err = u.Save()
	utils.CheckError(err)

	id := u.Id.Hex()

	w := workspace.New(oname, oid, workspace.WorkSpacePersonal, []string{})
	err = w.Save()
	utils.CheckError(err)

	_, err = clerk_utils.ClerkClient().Users().Update(cuid, &clerk.UpdateUser{
		ExternalID: &id,
		PublicMetadata: map[string]interface{}{
			"workspaceId": w.Id.Hex(),
		}})
	utils.CheckError(err)

	return c.Status(fiber.StatusOK).JSON(u)
}
