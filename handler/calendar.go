package handler

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/models/connections"
	"github.com/snipextt/dayer/utils"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func GetConnectedCalendars(c *fiber.Ctx) error {
	defer HandleInternalServerError(c)

	uid := c.Locals("uid").(string)
	conn, err := connections.FindConnectionsForUid(uid)
	utils.PanicOnError(err)

	return HandleSuccess(c, nil, conn)
}

func ListAllGoogleCalendarsForConnection(c *fiber.Ctx) error {
	defer HandleInternalServerError(c)
	ctx, cancel := utils.GetContext()
	defer cancel()

	cid := c.Query("calendar_id")
	uid := c.Locals("uid").(string)

	conn, err := connections.FindById(cid)

	if err != nil {
		return HandleBadRequest(c, "Calendar not found")
	}

	if conn.Uid != uid {
		return HandleBadRequest(c, "Calendar not found")
	}

	client := utils.GoogleOauthConfig.Client(ctx, &oauth2.Token{
		RefreshToken: conn.Token,
	})

	utils.PanicOnError(err)
	srv, err := calendar.NewService(ctx, option.WithHTTPClient(client))
	calendars, err := srv.CalendarList.List().MinAccessRole("writer").Do()
	utils.PanicOnError(err)

	return HandleSuccess(c, nil, calendars.Items)
}

func SyncGoogleCalendars(c *fiber.Ctx) error {
	defer HandleInternalServerError(c)

	body := make([]map[string]interface{}, 0)
	if err := c.BodyParser(&body); err != nil {
		return HandleBadRequest(c, "Invalid body")
	}

	log.Println(body)

	return nil
}
