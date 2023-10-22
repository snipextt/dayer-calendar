package handler

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/models/connections"
	"github.com/snipextt/dayer/utils"
	"github.com/snipextt/dayer/utils/calendar"
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

	cid := c.Query("connection_id")
	uid := c.Locals("uid").(string)

	conn, err := connections.FindById(cid)
	if err != nil {
		return HandleBadRequest(c, "Calendar not found")
	}
	if conn.Uid != uid {
		return HandleBadRequest(c, "Calendar not found")
	}

	srv, err := calendar.GoogleCalendarService(ctx, conn.Token)
	utils.PanicOnError(err)

	calendars, err := srv.CalendarList.List().MinAccessRole("writer").Do()
	utils.PanicOnError(err)

	return HandleSuccess(c, nil, calendars.Items)
}

func SyncGoogleCalendars(c *fiber.Ctx) error {
	defer HandleInternalServerError(c)
	ctx, cancel := utils.GetContext()
	defer cancel()

	cid := c.Query("connection_id")
	uid := c.Locals("uid").(string)

	body := make([]map[string]interface{}, 0)
	if err := c.BodyParser(&body); err != nil {
		return HandleBadRequest(c, "Invalid body")
	}

	conn, err := connections.FindById(cid)
	if err != nil {
		return HandleBadRequest(c, "Calendar not found")
	}
	if conn.Uid != uid {
		return HandleBadRequest(c, "Calendar not found")
	}

	calendars := make([]string, 0)
	for _, v := range body {
		if id, ok := v["id"].(string); ok {
			calendars = append(calendars, id)
		}
	}

	var wg sync.WaitGroup
	srv, err := calendar.GoogleCalendarService(ctx, conn.Token)
	utils.PanicOnError(err)

	for _, v := range calendars {
		wg.Add(1)
		go calendar.AddCalendarConnection(v, srv, &wg)
	}

	wg.Wait()
	return HandleSuccess(c, nil, nil)
}
