package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
	"github.com/joho/godotenv"

	"github.com/roolps/logging"
	"github.com/roolps/lootlocker-sudoku/backend/pkg/router"
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
	// send these values down to router
	router.New(logger, wd, env)

	// file server for static files
	fs := http.FileServer(http.Dir(fmt.Sprintf("%v/public/build", wd)))
	http.Handle("/static/", fs)

	// fallback to index.html
	http.HandleFunc("/", router.Handle)

	logger.Debugf("starting webserver on %v:%v", env["ADDRESS"], env["PORT"])
	http.ListenAndServe(fmt.Sprintf("%v:%v", env["ADDRESS"], env["PORT"]), nil)
}
