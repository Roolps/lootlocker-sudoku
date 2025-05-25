package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/roolps/lootlocker-sudoku/backend/pkg/lootlocker"
)

type logoutEndpoint struct{}

func (e *logoutEndpoint) Post(s *session, w http.ResponseWriter, raw []byte) *apiresponse {
	c := &http.Cookie{
		Name:     "session",
		Value:    "",
		Path:     "/",
		Domain:   origin,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, c)
	return statusok(nil)
}

type stateEndpoint struct {
}

type cell struct {
	Value     *int   `json:"value,omitempty"`
	Immutable bool   `json:"immutable,omitempty"`
	Pencil    []bool `json:"pencil,omitempty"`
}

func (e *stateEndpoint) Get(s *session, w http.ResponseWriter) *apiresponse {
	metadata, err := lootlockerClient.ListMetadata(lootlocker.MetadataSourcePlayer, s.PlayerID, s.Token)
	if err != nil {
		if err == lootlocker.ErrForbidden {
			// return a 403 error as the login token is no longer valid
			c := &http.Cookie{
				Name:     "session",
				Value:    "",
				Path:     "/",
				Domain:   origin,
				Expires:  time.Unix(0, 0),
				MaxAge:   -1,
				Secure:   true,
				HttpOnly: true,
				SameSite: http.SameSiteLaxMode,
			}
			http.SetCookie(w, c)

			return statusForbidden(err.Error())
		}
		return statusinternalservererror(err.Error())
	}
	if gameState, ok := metadata["current_state"]; ok {
		return statusok(gameState.Value)
	}

	// return empty game state as there is none active
	// this will force user to use the menu
	return statusok([][]cell{})
}

func (e *stateEndpoint) Post(s *session, w http.ResponseWriter, raw []byte) *apiresponse {
	var gameState [][]cell
	if err := json.Unmarshal(raw, &gameState); err != nil {
		return statusinternalservererror(err.Error())
	}

	// update current game state
	if err := lootlockerClient.UpdatePlayerMetadata(s.Token, []lootlocker.Metadata{{
		Access: []string{
			"game_api.read",
			"game_api.write",
		},
		Key:    "current_state",
		Tags:   []string{},
		Value:  gameState,
		Type:   lootlocker.MetadataTypeJSON,
		Action: lootlocker.MetadataActionUpdate,
	}}); err != nil {
		return statusinternalservererror(err.Error())
	}
	return statusok(nil)
}

type gameEndpoint struct {
	Difficulty string `json:"difficulty"`
}

func (e *gameEndpoint) Post(s *session, w http.ResponseWriter, raw []byte) *apiresponse {
	if err := json.Unmarshal(raw, e); err != nil {
		return statusnotacceptable(err.Error())
	}
	switch e.Difficulty {
	case "easy":
	default:
		return statusnotacceptable(fmt.Sprintf("invalid game difficulty %v", e.Difficulty))
	}

	gamestate := [][]cell{}
	raw, err := os.ReadFile(fmt.Sprintf("%v/example_games/%v.json", wd, e.Difficulty))
	if err != nil {
		return statusinternalservererror(err.Error())
	}
	if err := json.Unmarshal(raw, &gamestate); err != nil {
		return statusinternalservererror(err.Error())
	}

	// create current game state
	if err := lootlockerClient.UpdatePlayerMetadata(s.Token, []lootlocker.Metadata{
		{
			Access: []string{
				"game_api.read",
				"game_api.write",
			},
			Key:    "current_state",
			Tags:   []string{},
			Value:  gamestate,
			Type:   lootlocker.MetadataTypeJSON,
			Action: lootlocker.MetadataActionCreate,
		},
	}); err != nil {
		return statusinternalservererror(err.Error())
	}
	return statusok(gamestate)
}
