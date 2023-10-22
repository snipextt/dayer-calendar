package calendar

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/snipextt/dayer/utils"
	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
	"google.golang.org/api/option"
)

func GoogleCalendarService(ctx context.Context, tok string) (srv *calendar.Service, err error) {
	client := utils.GoogleOauthConfig.Client(ctx, &oauth2.Token{
		RefreshToken: tok,
	})

	srv, err = calendar.NewService(ctx, option.WithHTTPClient(client))
	return
}

func NewChannelForGoogleCalendar(id string, tok string) *calendar.Channel {
	return &calendar.Channel{
		Address: "",
		Id:      id,
		Type:    "web_hook",
		Token:   tok,
	}

}

func ParseGoogleCalendarEvent(event *calendar.Event, calendarId string) (e *Event) {
	if event.Status == "cancelled" {
		return nil
	}
	e = &Event{
		Title:       event.Summary,
		Description: event.Description,
		Start:       event.Start.DateTime,
		End:         event.End.DateTime,
		Location:    event.Location,
	}
	rec := len(event.Recurrence) > 0
	if rec {
		e.RecurrenceRule = ParseReccurenceRule(event.Recurrence[0])
		e.IsRecurring = true
		e.RecurringId = event.RecurringEventId
	}
	return
}

func ParseReccurenceRule(rrule string) RecurrenceRule {
	rrule = strings.TrimPrefix(rrule, "RRULE:")
	parsed_rules := strings.Split(rrule, ";")
	var rules RecurrenceRule
	for _, rule := range parsed_rules {
		k, v := strings.Split(rule, "=")[0], strings.Split(rule, "=")[1]
		switch k {
		case "FREQ":
			rules.Frequency = v
		case "BYDAY":
			rules.Days = strings.Split(v, ",")
		case "COUNT":
			{

				i, _ := strconv.Atoi(v)
				rules.Count = i
			}
		case "UNTIL":
			rules.Until = v
		}
	}
	return rules
}

func AddGoogleCalendarEvent(e *calendar.Event, calendarId string, wg *sync.WaitGroup) {
	defer wg.Done()
	event := ParseGoogleCalendarEvent(e, calendarId)
	if event == nil {
		return
	}
	event.Save()
}

func AddCalendarConnection(id string, srv *calendar.Service, wg *sync.WaitGroup) {
	tmax := time.Now().AddDate(0, 0, 90).Format(time.RFC3339)
	tmin := time.Now().Format(time.RFC3339)
	defer wg.Done()
	cal := NewGoogleCalendarConnection(id)
	cal.Save()
	// watch, _ := srv.Events.Watch(v.Id, calendar.NewChannelForGoogleCalendar(id, conn.Token)).Do()
	// conn.ChannelId = watch.Id
	// conn.Save()

	events, _ := srv.Events.List(id).ShowDeleted(false).TimeMin(tmin).TimeMax(tmax).Do()
	for _, v := range events.Items {
		wg.Add(1)
		go AddGoogleCalendarEvent(v, id, wg)
	}
}
