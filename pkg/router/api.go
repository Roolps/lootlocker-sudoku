package router

import (
	"encoding/json"
	"io"
	"net/http"
)

type apiresponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`

	Data any `json:"data,omitempty"`
}

func (s *session) apiRequest(path string, w http.ResponseWriter, r *http.Request) *apiresponse {
	logger.Debugf("%s [%d] %s %s", r.Header.Get("X-Real-IP"), http.StatusOK, r.Method, r.URL.Path)
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
	logger.Error(msg)
	return &apiresponse{Code: http.StatusInternalServerError, Message: msg}
}

func statusnotacceptable(msg string) *apiresponse {
	return &apiresponse{Code: http.StatusNotAcceptable, Message: msg}
}

func statusok(data any) *apiresponse {
	return &apiresponse{Code: http.StatusOK, Message: "success", Data: data}
}
