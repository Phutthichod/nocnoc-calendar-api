package model

type InstallationData struct {
	Email     string      `json:"email"`
	Name      string      `json:"name"`
	EventList []EventList `json:"events"`
}
