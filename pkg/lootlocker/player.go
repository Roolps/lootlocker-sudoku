package lootlocker

import (
	"encoding/json"
	"log"
	"net/http"
)

type PlayerData struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Order string `json:"order"`
}

// https://api.lootlocker.com/game/v1/player/storage

func (c *Client) GetPlayerStorage(session string) ([]PlayerData, error) {
	raw, err := c.Request(http.MethodGet, "game/v1/player/storage", application_json, nil, map[string]string{"x-session-token": session})
	log.Println(string(raw))
	return nil, err
}

// https://api.lootlocker.com/game/v1/player/storage

func (c *Client) UpdatePlayerStorage(session string, data []PlayerData) error {
	raw, err := json.Marshal(map[string]any{"payload": data})
	if err != nil {
		return err
	}
	_, err = c.Request(http.MethodPost, "game/v1/player/storage", application_json, raw, map[string]string{"x-session-token": session})
	return err
}
