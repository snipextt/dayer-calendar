package handler

import (
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func catchInternalServerError(c *fiber.Ctx) {
	if err := recover(); err != nil {
		log.Println(string(debug.Stack()), err)
		c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Error:   "Internal Server Error",
			Message: "Something went wrong, please try again later",
		})
	}
}

func badRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(models.Response{
		Error:   "Bad Request",
		Message: message,
	})
}

func unauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
		Error:   "Bad Request",
		Message: message,
	})
}

func success(c *fiber.Ctx, msg interface{}, data interface{}) error {
	if msg == nil {
		msg = ""
	}
	return c.Status(fiber.StatusOK).JSON(models.Response{
		Message: msg.(string),
		Data:    data,
	})
}

func getWorkspaceId(c *fiber.Ctx) (id primitive.ObjectID, err error) {
	wid, ok := c.GetReqHeaders()["X-Workspace-Id"]
	if !ok {
		err = noWorkspaceId
		return
	}
	id, err = primitive.ObjectIDFromHex(wid)
	if err != nil {
		err = invdalidWorkSpaceId
	}
	return
}
