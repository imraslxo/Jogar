package models

type User struct {
	ID              uint64 `json:"id"`
	TgUsername      string `json:"username"`
	TgFirstName     string `json:"first_name"`
	TgLastName      string `json:"last_name"`
	PhotoURL        string `json:"photo_url,omitempty"`
	IsPremium       bool   `json:"is_premium"`
	AuthDate        int64  `json:"auth_date"`
	UILanguageCode  string `json:"language_code"`
	AllowsWriteToPM bool   `json:"allows_write_to_pm"`

	ProfileID *uint64  `json:"profile_id,omitempty"`
	Profile   *Profile `json:"profile,omitempty"`

	TeamID *uint64 `json:"team_id,omitempty"`
	Team   *Team   `json:"team,omitempty"`
}

type AuthRequestDTO struct {
	TgUsername      string `json:"username"`
	TgFirstName     string `json:"first_name"`
	TgLastName      string `json:"last_name"`
	PhotoURL        string `json:"photo_url,omitempty"`
	IsPremium       bool   `json:"is_premium"`
	AuthDate        int64  `json:"auth_date"`
	UILanguageCode  string `json:"language_code"`
	AllowsWriteToPM bool   `json:"allows_write_to_pm"`
}

type UserCreateRequest struct {
	TgUsername           string               `json:"username" binding:"required"`
	TgFirstName          string               `json:"first_name" binding:"required"`
	TgLastName           string               `json:"last_name" binding:"required"`
	PhotoURL             string               `json:"photo_url,omitempty"`
	UILanguageCode       string               `json:"language_code" binding:"required"`
	AllowsWriteToPM      bool                 `json:"allows_write_to_pm"`
	ProfileCreateRequest ProfileCreateRequest `json:"profile_create_request"`
}
