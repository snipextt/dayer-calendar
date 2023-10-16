package handler

import (
	"log"
	"runtime/debug"

	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/models"
)

func HandleInternalServerError(c *fiber.Ctx) {
	if err := recover(); err != nil {
		log.Println(string(debug.Stack()), err)
		c.Status(fiber.StatusInternalServerError).JSON(models.Response{
			Error:   "Internal Server Error",
			Message: "Something went wrong, please try again later",
		})
	}
}

func HandleBadRequest(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusBadRequest).JSON(models.Response{
		Error:   "Bad Request",
		Message: message,
	})
}

func HandleUnauthorized(c *fiber.Ctx, message string) error {
	return c.Status(fiber.StatusUnauthorized).JSON(models.Response{
		Error:   "Bad Request",
		Message: message,
	})
}

func HandleSuccess(c *fiber.Ctx, msg interface{}, data interface{}) error {
	if msg == nil {
		msg = ""
	}
	return c.Status(fiber.StatusOK).JSON(models.Response{
		Message: msg.(string),
		Data:    data,
	})
}
