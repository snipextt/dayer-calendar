package handler

import (
	"github.com/clerkinc/clerk-sdk-go/clerk"
	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/models/user"
	"github.com/snipextt/dayer/utils"
	"github.com/snipextt/dayer/utils/clerk"
)

func Onboarding(c *fiber.Ctx) error {
	HandleInternalServerError(c)

	cuid := c.Locals("auth").(*clerk.TokenClaims).Claims.Subject
	cu, err := clerk_utils.GetClerkUser(cuid)
	utils.PanicOnError(err)

	if cu.ExternalID != nil {
		return HandleBadRequest(c, "User already onboarded")
	}

	u := user.New(cuid)
	err = u.Save()
	utils.PanicOnError(err)

	id := u.Id.Hex()
	clerk_utils.ClerkClient().Users().Update(cuid, &clerk.UpdateUser{
		ExternalID: &id,
	})

	return c.Status(fiber.StatusOK).JSON(u)
}
