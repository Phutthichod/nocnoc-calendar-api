package model

import "time"

type DateDisable struct {
	Date        time.Time `json:"date"`
	FullDayBusy bool      `json:"fullDayBusy"`
}
