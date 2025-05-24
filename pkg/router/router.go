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
)

func New(loggerprofile *logging.Profile, workingdir string, env map[string]string) {
	// set up lootlocker client
	lootlockerClient.DomainKey = env["LOOTLOCKER_DOMAIN_KEY"]

	// to do: make this dynamic based on if env says dev mode or not
	lootlockerClient.IsDevelopment = true

	wd = workingdir
	logger = loggerprofile
	origin = env["ORIGIN"]
}

func Handle(w http.ResponseWriter, r *http.Request) {
	s := &session{}
	if err := s.get(r); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if strings.HasPrefix(r.URL.Path, "/api") {
		// execute api request
		path := strings.ReplaceAll(r.URL.Path, "/api", "")
		switch path {
		case "/login":
			s.login(w, r)

		default:
			if !s.LoggedIn {
				respond(http.StatusForbidden, "forbidden", w)
				return
			}
			s.apiRequest(w, r)
		}
		return
	}

	http.ServeFile(w, r, fmt.Sprintf("%v/public/build/index.html", wd))
}
