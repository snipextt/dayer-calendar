package handler

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

func GetPagination(ctx *fiber.Ctx) (int64, int64) {
	page, _ := strconv.ParseInt(ctx.Query("page", "1"), 10, 64)
	perPage, _ := strconv.ParseInt(ctx.Query("perPage", "10"), 10, 64)

	return page, perPage
}
