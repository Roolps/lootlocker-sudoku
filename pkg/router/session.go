package router

import (
	"net/http"
	"strings"
)

type session struct {
	LoggedIn bool
	Token    string
	Email    string
}

func (s *session) get(r *http.Request) error {
	c, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		return nil
	}
	if err != nil {
		return err
	}
	return s.verify(c)
}

func (s *session) verify(c *http.Cookie) error {
	// c.value will be in the format of email:token
	cookieval := strings.Split(c.Value, ":")
	if len(cookieval) != 2 {
		return nil
	}
	valid, err := lootlockerClient.VerifyWhiteLabelSession(cookieval[0], cookieval[1])
	if err != nil {
		// to do: make a lootlocker const for the authentication details are wrong vs internal server error
		return nil
	}
	if valid {
		// valid login session - set the values in session
		s.LoggedIn = true
		s.Email = cookieval[0]
		s.Token = cookieval[1]
	}
	return nil
}
