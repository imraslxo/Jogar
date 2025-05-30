package models

import "time"

type Team struct {
	ID          uint64    `json:"id,omitempty"`
	TeamName    string    `json:"team_name"`
	Photo       string    `json:"photo,omitempty"`
	PlayingIn   time.Time `json:"playing_in" swaggertype:"string" format:"date-time"`
	Stadium     string    `json:"stadium"`
	Description string    `json:"description,omitempty"`
}

type TeamCreateRequest struct {
	TeamName    string    `json:"team_name"`
	Photo       string    `json:"photo,omitempty"`
	PlayingIn   time.Time `json:"playing_in" swaggertype:"string" format:"date-time"`
	Stadium     string    `json:"stadium"`
	Description string    `json:"description,omitempty"`
}

type TeamWithCount struct {
	Name          string `json:"name"`
	Players_count int    `json:"players_count"`
}
