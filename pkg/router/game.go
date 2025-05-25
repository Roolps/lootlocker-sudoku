package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/roolps/lootlocker-sudoku/backend/pkg/lootlocker"
)

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

	// empty game :)
	val := make([][]map[string]interface{}, 9)
	for i := range val {
		val[i] = make([]map[string]interface{}, 9)
		for j := range val[i] {
			val[i][j] = map[string]interface{}{}
		}
	}
	return statusok(val)
}

func (e *stateEndpoint) Post(s *session, w http.ResponseWriter, raw []byte) *apiresponse {
	var gameState [][]cell
	if err := json.Unmarshal(raw, &gameState); err != nil {
		return statusinternalservererror(err.Error())
	}
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
		logger.Error(err)
		return statusinternalservererror(err.Error())
	}
	return statusok(nil)
}
