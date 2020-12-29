package model

import "golang.org/x/oauth2"

type Token struct {
	Data       *oauth2.Token `json:"token"`
	CalendarId string        `json:"calendarId"`
	Date       Date          `json:"date"`
}
