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

	PlayerID string
	Email    string
	Token    string
}

func (s *session) get(r *http.Request) error {
	sessionCookie, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		return nil
	}
	if err != nil {
		return err
	}
	return s.verify(sessionCookie)
}

func (s *session) verify(c *http.Cookie) error {
	// c.value will be in the format of email:token
	cookieval := strings.Split(c.Value, ":")
	if len(cookieval) != 3 {
		return nil
	}
	s.LoggedIn = true

	s.PlayerID = cookieval[0]
	s.Email = cookieval[1]
	s.Token = cookieval[2]
	return nil
}

type userlogin struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (s *session) login(w http.ResponseWriter, r *http.Request) *apiresponse {
	raw, err := io.ReadAll(r.Body)
	if err != nil {
		return statusnotacceptable(fmt.Sprintf("failed to read request body: %v", err))
	}
	ul := &userlogin{}
	if err := json.Unmarshal(raw, ul); err != nil {
		return statusinternalservererror(fmt.Sprintf("failed to unmarshal body: %v", err))
	}

	token, err := lootlockerClient.LoginWhiteLabelUser(ul.Email, ul.Password, true)
	if err != nil {
		// to do: add error constants to lootlocker package to check the actual error message
		return statusnotacceptable("invalid username or password")
	}

	// start white label session using the token
	gameSession, err := lootlockerClient.StartWhiteLabelSession(ul.Email, token, "1.0")
	if err != nil {
		return statusinternalservererror(err.Error())
	}

	sessionInfo, err := lootlockerClient.GetInfoFromSession(gameSession.SessionToken)
	if err != nil {
		return statusinternalservererror(err.Error())
	}

	c := &http.Cookie{
		Name:     "session",
		Value:    fmt.Sprintf("%v:%v:%v", sessionInfo.ID, ul.Email, gameSession.SessionToken),
		Path:     "/",
		Domain:   origin,
		Expires:  time.Now().Add(24 * time.Hour),
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, c)
	return statusok(nil)
}
