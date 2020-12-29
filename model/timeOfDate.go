package model

import "google.golang.org/api/calendar/v3"

type TimeOfDate struct {
	Times []*calendar.TimePeriod
}
