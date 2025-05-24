package router

import "net/http"

type session struct {
	Token string
	Email string
}

func (s *session) get(w http.ResponseWriter, r *http.Request) error {
	c, err := r.Cookie("session")
	if err == http.ErrNoCookie {
		return s.create(w)
	}
	if err != nil {
		return err
	}
	return s.verify(c, w)
}

// create new guest session for user
func (s *session) create(w http.ResponseWriter) error {
	return nil
}

// check if session is valid - if not then create new guest session
func (s *session) verify(c *http.Cookie, w http.ResponseWriter) error {
	// c.value will be in the format of email:token
	return nil
}
