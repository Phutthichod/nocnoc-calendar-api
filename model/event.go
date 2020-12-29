package model

import (
	"time"

	"golang.org/x/oauth2"
	"google.golang.org/api/calendar/v3"
)

type Event struct {
	Summary string
	Start   calendar.EventDateTime
	End     calendar.EventDateTime
}

type Event2 struct {
	Token      oauth2.Token `json:"token"`
	CalendarId string       `json:"calendarId"`
	Summary    string       `json:"summary"`
	Start      time.Time    `json:"start"`
	End        time.Time    `json:"end"`
}
