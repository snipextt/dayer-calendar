package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/models/extension"
	"github.com/snipextt/dayer/utils"
)

func GetExtensions(ctx *fiber.Ctx) error {
	defer HandleInternalServerError(ctx)

	page, limit := GetPagination(ctx)
	extensions, err := extension.PaginatedExtensions(page, limit)
	utils.PanicOnError(err)

	return HandleSuccess(ctx, nil, extensions)
}
