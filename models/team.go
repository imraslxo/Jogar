package models

import "time"

type Team struct {
	ID          uint64    `json:"id,omitempty"`
	TeamName    string    `json:"team_name"`
	Photo       string    `json:"photo,omitempty"`
	PlayingIn   time.Time `json:"playing_in" swaggertype:"string" format:"date-time"`
	Stadium     string    `json:"stadium"`
	Discription string    `json:"discription,omitempty"`
}

type TeamCreateRequest struct {
	TeamName    string    `json:"team_name"`
	Photo       string    `json:"photo,omitempty"`
	PlayingIn   time.Time `json:"playing_in" swaggertype:"string" format:"date-time"`
	Stadium     string    `json:"stadium"`
	Discription string    `json:"discription,omitempty"`
}
