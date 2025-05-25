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

	switch r.Method {
	case http.MethodGet:
		obj := whichGetEndpoint(path)
		if obj != nil {
			return obj.Get(s, w)
		}
	case http.MethodPost:
		obj := whichPostEndpoint(path)
		if obj != nil {
			raw, err := io.ReadAll(r.Body)
			if err != nil && err != io.EOF {
				return statusinternalservererror(err.Error())
			}
			return obj.Post(s, w, raw)
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

type getEndpoint interface {
	Get(*session, http.ResponseWriter) *apiresponse
}

func whichGetEndpoint(e string) getEndpoint {
	switch e {
	case "/state":
		return &stateEndpoint{}
	}
	return nil
}

type postEndpoint interface {
	Post(*session, http.ResponseWriter, []byte) *apiresponse
}

func whichPostEndpoint(e string) postEndpoint {
	switch e {
	case "/logout":
		return &logoutEndpoint{}
	case "/state":
		return &stateEndpoint{}
	case "/game":
		return &gameEndpoint{}
	}
	return nil
}
