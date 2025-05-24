package router

import (
	"encoding/json"
	"net/http"
)

type apiresponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (s *session) apiRequest(w http.ResponseWriter, r *http.Request) {
	respond(http.StatusOK, "golang backend working", w)

	logger.Debugf("%s [%d] %s %s", r.Header.Get("X-Real-IP"), http.StatusOK, r.Method, r.URL.Path)
}

func respond(code int, message string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(apiresponse{Code: code, Message: message})
}
