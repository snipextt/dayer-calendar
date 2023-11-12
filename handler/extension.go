package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/models/extension"
	"github.com/snipextt/dayer/utils"
)

func GetExtensions(ctx *fiber.Ctx) error {
	defer catchInternalServerError(ctx)

	page, limit := GetPagination(ctx)
	extensions, err := extension.PaginatedExtensions(page, limit)
	utils.CheckError(err)

	return success(ctx, nil, extensions)
}
