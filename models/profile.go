package models

type Profile struct {
	ID               uint64 `json:"id,omitempty"`
	AppLanguage      string `json:"app_language,omitempty"`
	PrefPosition     string `json:"pref_position"`
	Height           uint64 `json:"height"`
	Foot             string `json:"foot"`
	Age              string `json:"age"`
	PlayingFrequency string `json:"playing_frequency"`
	GamesPlayed      uint64 `json:"games_played"`
	City             string `json:"city"`
	Country          string `json:"country"`
	TgUserID         uint64 `json:"user_id"`
}

type ProfileCreateFirstDTO struct {
	AppLanguage      string `json:"app_language,omitempty"`
	Country          string `json:"country,omitempty"`
	City             string `json:"city,omitempty"`
	Age              string `json:"age"`
	PrefPosition     string `json:"pref_position"`
	PlayingFrequency string `json:"playing_frequency"`
}

type ProfileCreateRequest struct {
	PrefPosition string `json:"pref_position,omitempty"`
	Height       uint64 `json:"height,omitempty"`
	Foot         string `json:"foot,omitempty"`
	Age          string `json:"age"`
	City         string `json:"city"`
	Country      string `json:"country"`
	UserID       uint64 `json:"user_id"`
}
