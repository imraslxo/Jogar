package models

type Profile struct {
	ID               uint64 `json:"id,omitempty"`
	AppLanguage      string `json:"app_language,omitempty"`
	PrefPosition     string `json:"pref_position"`
	Height           uint64 `json:"height"`
	Foot             string `json:"foot"`
	Age              uint64 `json:"age"`
	PlayingFrequency uint64 `json:"playing_frequency"`
	GamesPlayed      uint64 `json:"games_played"`
	City             string `json:"city"`
	Country          string `json:"country"`
	UserID           uint64 `json:"user_id"`
}

type ProfileCreateFirstDTO struct {
	AppLanguage string `json:"app_language,omitempty"`
	Country     string `json:"country,omitempty"`
	City        string `json:"city,omitempty"`
}

type ProfileCreateRequest struct {
	PrefPosition string `json:"pref_position,omitempty"`
	Height       uint64 `json:"height,omitempty"`
	Foot         string `json:"foot,omitempty"`
	Age          uint64 `json:"age"`
	City         string `json:"city"`
	Country      string `json:"country"`
	UserID       uint64 `json:"user_id"`
}
