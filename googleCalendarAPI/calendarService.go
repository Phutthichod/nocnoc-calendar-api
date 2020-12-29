package googleCalendarAPI

import (
	"api-calendar/model"

	"google.golang.org/api/calendar/v3"
)

type Calendar struct {
	Service *calendar.Service
}

func (cal Calendar) CreatEvent(event model.Event, calendarId string) (*calendar.Event, error) {
	newEvent := calendar.Event{
		Summary: event.Summary,
		Start:   &event.Start,
		End:     &event.End,
		Attendees: []*calendar.EventAttendee{
			&calendar.EventAttendee{
				Email: "pin2041to@gmail.com",
			}},
	}
	calEvent, err := cal.Service.Events.Insert("primary", &newEvent).SendUpdates("all").Do()
	return calEvent, err
}

func (cal Calendar) ListEvent(calendarId string) (*calendar.Events, error) {
	// t := time.Now().Format(time.RFC3339)
	event, err := cal.Service.Events.List(calendarId).ShowDeleted(false).
		SingleEvents(true).TimeMin("2020-12-1T00:00:00+07:00").TimeMax("2020-12-31T00:00:00+07:00").MaxResults(10).OrderBy("startTime").Do()
	return event, err
}

func (cal Calendar) GetDayFee(calendarId string) (*calendar.FreeBusyResponse, error) {
	// t := time.Now().Format(time.RFC3339)
	res, err := cal.Service.Freebusy.Query(
		&calendar.FreeBusyRequest{
			TimeMin:  "2020-12-01T08:00:00+07:00",
			TimeMax:  "2020-12-30T20:00:00+07:00",
			TimeZone: "+0700",
			Items: []*calendar.FreeBusyRequestItem{&calendar.FreeBusyRequestItem{
				Id: "pin2041to@gmail.com",
			},
			}},
	).Do()
	return res, err
}
func (cal Calendar) GetDayBusyByDate(calendarId string, date model.Date) (*calendar.FreeBusyResponse, error) {
	// t := time.Now().Format(time.RFC3339)
	res, err := cal.Service.Freebusy.Query(
		&calendar.FreeBusyRequest{
			TimeMin:  date.Year + "-" + date.Month + "-" + date.Day + "T00:00:00+07:00",
			TimeMax:  date.Year + "-" + date.Month + "-" + date.Day + "T23:30:00+07:00",
			TimeZone: "+0700",
			Items: []*calendar.FreeBusyRequestItem{&calendar.FreeBusyRequestItem{
				Id: "pin2041to@gmail.com",
			},
			}},
	).Do()
	return res, err
}
