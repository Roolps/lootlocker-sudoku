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

	SessionToken string `json:"session_token,omitempty"`
}

// https://api.lootlocker.com/white-label-login/sign-up

func (c *Client) CreateWhiteLabelLogin(email, password string) (*WhiteLabelLogin, error) {
	raw, err := json.Marshal(map[string]any{"email": email, "password": password})
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

// https://api.lootlocker.com/white-label-login/login
// returns user session token

func (c *Client) LoginWhiteLabelUser(email, password string, rememberMe bool) (string, error) {
	raw, err := json.Marshal(map[string]any{"email": email, "password": password, "remember": rememberMe})
	if err != nil {
		return "", err
	}
	raw, err = c.Request(http.MethodPost, "white-label-login/login", application_json, raw, nil)
	if err != nil {
		return "", err
	}
	login := &WhiteLabelLogin{}
	if err := json.Unmarshal(raw, login); err != nil {
		return "", err
	}
	return login.SessionToken, nil
}

// https://api.lootlocker.com/white-label-login/verify-session
// returns true/false for if the session is valid

func (c *Client) VerifyWhiteLabelSession(email, token string) (bool, error) {
	raw, err := json.Marshal(map[string]any{"email": email, "token": token})
	if err != nil {
		return false, err
	}
	_, err = c.Request(http.MethodPost, "white-label-login/verify-session", application_json, raw, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

type GameSession struct {
	Success      bool   `json:"success"`
	SessionToken string `json:"session_token"`

	PlayerID        int       `json:"player_id"`
	PublicUID       string    `json:"public_uid"`
	PlayerCreatedAt time.Time `json:"player_created_at"`

	CheckGrantNotifications        bool `json:"check_grant_notifications"`
	CheckDeactivationNotifications bool `json:"check_deactivation_notifications"`
	SeenBefore                     bool `json:"seen_before"`
}

// https://api.lootlocker.com/game/v2/session/white-label

func (c *Client) StartWhiteLabelSession(email, token, gameVersion string) (*GameSession, error) {
	raw, err := json.Marshal(map[string]any{"game_key": c.GameKey, "email": email, "token": token, "game_version": gameVersion})
	if err != nil {
		return nil, err
	}
	raw, err = c.Request(http.MethodPost, "game/v2/session/white-label", application_json, raw, nil)
	if err != nil {
		return nil, err
	}
	sess := &GameSession{}
	if err := json.Unmarshal(raw, sess); err != nil {
		return nil, err
	}
	return sess, nil
}
