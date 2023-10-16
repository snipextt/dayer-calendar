package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	clerk_utils "github.com/snipextt/dayer/utils/clerk"
)

func AuthMiddleware(c *fiber.Ctx) error {
	sessionToken := c.Get("Authorization")
	if sessionToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	sessionToken = strings.Split(sessionToken, " ")[1]
	if sessionToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	claims, err := clerk_utils.ClerkClient().DecodeToken(sessionToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	c.Locals("auth", claims)
	c.Locals("uid", claims.Extra["externalId"])
	return c.Next()
}
