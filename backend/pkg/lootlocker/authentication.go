package lootlocker

import (
	"encoding/json"
	"net/http"
	"time"
)

type WhiteLabelLogin struct {
	ID       int    `json:"id"`
	GameID   int    `json:"game_id"`
	Email    string `json:"email"`
	Password string `json:"password"`

	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at"`
	ValidatedAt *time.Time `json:"validated_at"`
}

// https://api.lootlocker.com/white-label-login/sign-up

func (c *Client) CreateWhiteLabelLogin(email string, password string) (*WhiteLabelLogin, error) {
	raw, err := json.Marshal(map[string]string{"email": email, "password": password})
	if err != nil {
		return nil, err
	}
	raw, err = c.Request(http.MethodPost, "white-label-login/sign-up", application_json, raw, nil)
	if err != nil {
		return nil, err
	}
	login := &WhiteLabelLogin{}
	if err := json.Unmarshal(raw, login); err != nil {
		return nil, err
	}
	return login, nil
}
