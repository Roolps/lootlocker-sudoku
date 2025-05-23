package main

import (
	"encoding/json"
	"net/http"

	"github.com/fatih/color"

	"github.com/roolps/logging"
)

var (
	logger = &logging.Profile{
		Prefix: "Lootlocker",
		Color:  color.RGB(28, 232, 109),
	}
)

func main() {
	logger.EnableDebug()

	http.HandleFunc("/", handler)

	logger.Debug("starting webserver on port :8080")
	http.ListenAndServe(":8080", nil)
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "golang backend working"})

	logger.Debugf("%s [%d] %s %s", r.Header.Get("X-Real-IP"), http.StatusOK, r.Method, r.URL.Path)
}
