package router

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/roolps/lootlocker-sudoku/backend/pkg/lootlocker"
)

type apiresponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`

	Data any `json:"data,omitempty"`
}

func (s *session) apiRequest(path string, w http.ResponseWriter, r *http.Request) *apiresponse {
	logger.Debugf("%s [%d] %s %s", r.Header.Get("X-Real-IP"), http.StatusOK, r.Method, r.URL.Path)
	logger.Debug(path)
	switch path {
	// get the game state
	case "/state":
		e := &stateEndpoint{}
		switch r.Method {
		case http.MethodGet:
			return e.Get(s, w)
		case http.MethodPost:
			raw, err := io.ReadAll(r.Body)
			if err != nil {
				return statusinternalservererror(err.Error())
			}
			return e.Post(s, w, raw)
		}

	}

	return statusnotfound()
}

// api outcomes
func (r *apiresponse) Write(w http.ResponseWriter) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(r.Code)
	return json.NewEncoder(w).Encode(r)
}

func statusnotfound() *apiresponse {
	return &apiresponse{Code: http.StatusNotFound, Message: "not found"}
}

func statusForbidden(msg string) *apiresponse {
	return &apiresponse{Code: http.StatusForbidden, Message: msg}
}

func statusinternalservererror(msg string) *apiresponse {
	return &apiresponse{Code: http.StatusInternalServerError, Message: msg}
}

func statusnotacceptable(msg string) *apiresponse {
	return &apiresponse{Code: http.StatusNotAcceptable, Message: msg}
}

func statusok(data any) *apiresponse {
	return &apiresponse{Code: http.StatusOK, Message: "success", Data: data}
}

type stateEndpoint struct {
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
	// to do: this is the to do for tomorrow morning! let's sort out the update metadata so we can keep state saved to player account
	return statusok(nil)
}
