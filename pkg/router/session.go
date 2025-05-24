package router

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
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

type userlogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *session) login(w http.ResponseWriter, r *http.Request) {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		respond(http.StatusNotAcceptable, fmt.Sprintf("failed to read request body: %v", err), w)
		return
	}
	ul := &userlogin{}
	if err := json.Unmarshal(raw, ul); err != nil {
		respond(http.StatusInternalServerError, fmt.Sprintf("failed to unmarshal body: %v", err), w)
		return
	}
	token, err := lootlockerClient.LoginWhiteLabelUser(ul.Email, ul.Password, true)
	if err != nil {
		// to do: add error constants to lootlocker package to check the actual error message
		respond(http.StatusNotAcceptable, "invalid username or password", w)
		return
	}
	c := &http.Cookie{
		Name:     "session",
		Value:    fmt.Sprintf("%v:%v", ul.Email, token),
		Path:     "/",
		Domain:   origin,
		Expires:  time.Now().Add(time.Hour * 24 * 30),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, c)
	respond(http.StatusOK, "success", w)
}
