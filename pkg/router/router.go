package router

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/roolps/logging"
	"github.com/roolps/lootlocker-sudoku/backend/pkg/lootlocker"
)

var (
	lootlockerClient = &lootlocker.Client{}

	logger *logging.Profile
	wd     string

	origin string

	LOOTLOCKER_CURRENCY_ID string
)

func New(loggerprofile *logging.Profile, workingdir string, env map[string]string) {
	// set up lootlocker client
	lootlockerClient.DomainKey = env["LOOTLOCKER_DOMAIN_KEY"]
	lootlockerClient.GameKey = env["LOOTLOCKER_GAME_KEY"]

	// to do: make this dynamic based on if env says dev mode or not
	lootlockerClient.IsDevelopment = true

	wd = workingdir
	logger = loggerprofile
	origin = env["ORIGIN"]

	LOOTLOCKER_CURRENCY_ID = env["LOOTLOCKER_CURRENCY_ID"]
}

func Handle(w http.ResponseWriter, r *http.Request) {
	s := &session{}
	if err := s.get(r); err != nil {
		if err := statusinternalservererror(err.Error()).Write(w); err != nil {
			logger.Error(err)
		}
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api") {
		// execute api request
		path := strings.ReplaceAll(r.URL.Path, "/api", "")
		switch path {
		case "/login":
			if err := s.login(w, r).Write(w); err != nil {
				logger.Error(err)
			}
			return

		case "/signup":
			if err := s.signup(w, r).Write(w); err != nil {
				logger.Error(err)
			}
			return

		default:
			if !s.LoggedIn {
				if err := statusForbidden("forbidden").Write(w); err != nil {
					logger.Error(err)
				}
				return
			}
			if err := s.apiRequest(path, w, r).Write(w); err != nil {
				logger.Error(err)
			}
		}
		return
	}

	switch r.URL.Path {
	case "/manifest.json", "/logo192.png", "/logo512.png", "/favicon.ico":
		http.ServeFile(w, r, fmt.Sprintf("%v/public/build%v", wd, r.URL.Path))
		return
	}
	http.ServeFile(w, r, fmt.Sprintf("%v/public/build/index.html", wd))
}
