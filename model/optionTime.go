package model

type OptionTime struct {
	Id       int    `json:"id"`
	Start    string `json:"start"`
	End      string `json:"end"`
	Disabled bool   `json:"disabled"`
}
