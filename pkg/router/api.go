package router

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/roolps/lootlocker-sudoku/backend/pkg/lootlocker"
)

type apiresponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func (s *session) apiRequest(path string, w http.ResponseWriter, r *http.Request) {
	switch path {
	// get the game state
	case "/state":
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

				respond(http.StatusForbidden, err.Error(), nil, w)
			} else {
				respond(http.StatusInternalServerError, err.Error(), nil, w)
			}
			return
		}
		if gameState, ok := metadata["current_state"]; ok {
			respond(http.StatusOK, "success", gameState.Value, w)
		}
		// otherwise respond "create new game..."

	default:
		respond(http.StatusOK, "golang backend working", nil, w)
	}

	logger.Debugf("%s [%d] %s %s", r.Header.Get("X-Real-IP"), http.StatusOK, r.Method, r.URL.Path)
}

func respond(code int, message string, data any, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(apiresponse{Code: code, Message: message, Data: data})
}
