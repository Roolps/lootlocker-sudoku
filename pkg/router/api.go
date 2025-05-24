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
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(apiresponse{Code: http.StatusOK, Message: "golang backend working"})

	logger.Debugf("%s [%d] %s %s", r.Header.Get("X-Real-IP"), http.StatusOK, r.Method, r.URL.Path)
}
