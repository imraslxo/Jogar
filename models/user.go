package models

import "time"

type User struct {
	ID              uint64     `json:"id,omitempty"`
	TgUsername      string     `json:"tg_username"`
	TgFirstName     string     `json:"tg_first_name"`
	TgLastName      string     `json:"tg_last_name"`
	PhotoURL        string     `json:"photo_url,omitempty"`
	TgLastLogin     *time.Time `json:"tg_last_login,omitempty"`
	RegisteredAt    *time.Time `json:"registered_at,omitempty"`
	UILanguageCode  string     `json:"ui_language_code"`
	AllowsWriteToPM bool       `json:"allows_write_to_pm"`

	ProfileID *uint64  `json:"profile_id,omitempty"`
	Profile   *Profile `json:"profile,omitempty"`

	TeamID *uint64 `json:"team_id,omitempty"`
	Team   *Team   `json:"team,omitempty"`
}

type UserCreateRequest struct {
	TgUsername           string               `json:"tg_username" binding:"required"`
	TgFirstName          string               `json:"tg_first_name" binding:"required"`
	TgLastName           string               `json:"tg_last_name" binding:"required"`
	PhotoURL             string               `json:"photo_url,omitempty"`
	UILanguageCode       string               `json:"ui_language_code" binding:"required"`
	AllowsWriteToPM      bool                 `json:"allows_write_to_pm"`
	ProfileCreateRequest ProfileCreateRequest `json:"profile_create_request"`
}
