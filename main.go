package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"

	"github.com/roolps/logging"
)

var (
	logger = &logging.Profile{
		Prefix: "Lootlocker",
		Color:  color.RGB(28, 232, 109),
	}

	env = map[string]string{}
)

func main() {
	logger.EnableDebug()

	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %v", err)
	}
	env, err = godotenv.Read(fmt.Sprintf("%v/.env", wd))
	if err != nil {
		log.Fatalf("failed to get environment: %v", err)
	}

	// file server for static files
	fs := http.FileServer(http.Dir(fmt.Sprintf("%v/public/build", wd)))
	http.Handle("/static/", fs)

	// fallback to index.html
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, fmt.Sprintf("%v/public/build/index.html", wd))
	})

	http.HandleFunc("/api", apiHandler)

	logger.Debug("starting webserver on port :8080")
	http.ListenAndServe(fmt.Sprintf("%v:%v", env["ADDRESS"], env["PORT"]), nil)
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func apiHandler(w http.ResponseWriter, r *http.Request) {
	logger.Debug(r.Cookies())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{Message: "golang backend working"})

	logger.Debugf("%s [%d] %s %s", r.Header.Get("X-Real-IP"), http.StatusOK, r.Method, r.URL.Path)
}
