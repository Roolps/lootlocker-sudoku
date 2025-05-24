package router

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
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
		raw, err := os.ReadFile(fmt.Sprintf("%v/example.json", wd))
		if err != nil {
			logger.Error(err)
		}
		game := []any{}
		if err := json.Unmarshal(raw, &game); err != nil {
			logger.Error(err)
		}
		respond(http.StatusOK, "success", game, w)

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
