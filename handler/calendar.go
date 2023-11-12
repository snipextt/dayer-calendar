package handler

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"github.com/snipextt/dayer/models/connection"
	"github.com/snipextt/dayer/utils"
	"github.com/snipextt/dayer/utils/calendar"
)

func GetConnectedCalendars(c *fiber.Ctx) error {
	defer catchInternalServerError(c)

	uid := c.Locals("uid").(string)
	conn, err := connection.FindConnectionsForUid(uid)
	utils.CheckError(err)

	return success(c, nil, conn)
}

func ListAllGoogleCalendarsForConnection(c *fiber.Ctx) error {
	defer catchInternalServerError(c)
	ctx, cancel := utils.NewContext()
	defer cancel()

	cid := c.Query("connection_id")

	conn, err := connection.FindById(cid)
	if err != nil {
		return badRequest(c, "Calendar not found")
	}

	srv, err := calendar.GoogleCalendarService(ctx, conn.Token)
	utils.CheckError(err)

	calendars, err := srv.CalendarList.List().MinAccessRole("writer").Do()
	utils.CheckError(err)

	return success(c, nil, calendars.Items)
}

func SyncGoogleCalendars(c *fiber.Ctx) error {
	defer catchInternalServerError(c)
	ctx, cancel := utils.NewContext()
	defer cancel()

	cid := c.Query("connection_id")

	body := make([]map[string]interface{}, 0)
	if err := c.BodyParser(&body); err != nil {
		return badRequest(c, "Invalid body")
	}

	conn, err := connection.FindById(cid)
	if err != nil {
		return badRequest(c, "Calendar not found")
	}

	calendars := make([]string, 0)
	for _, v := range body {
		if id, ok := v["id"].(string); ok {
			calendars = append(calendars, id)
		}
	}

	var wg sync.WaitGroup
	srv, err := calendar.GoogleCalendarService(ctx, conn.Token)
	utils.CheckError(err)

	for _, v := range calendars {
		wg.Add(1)
		go calendar.AddCalendarConnection(v, srv, &wg)
	}

	wg.Wait()
	return success(c, nil, nil)
}
