package lootlocker

import (
	"log"
	"testing"
)

func TestAuthenticationEndpoint(t *testing.T) {
	c := &Client{
		DomainKey:     "",
		IsDevelopment: true,
	}
	login, err := c.CreateWhiteLabelLogin("example@example.com", "testingPassword")
	if err != nil {
		t.Error(err)
	}
	log.Println(login)
}

func TestWhiteLabelLoginEndpoint(t *testing.T) {
	c := &Client{
		DomainKey:     "",
		IsDevelopment: true,
	}
	token, err := c.LoginWhiteLabelUser("example@example.com", "testingPassword", true)
	if err != nil {
		t.Error(err)
	}
	log.Println(token)
}

func TestStartWhiteLabelSession(t *testing.T) {
	c := &Client{
		DomainKey:     "",
		IsDevelopment: true,
		GameKey:       "",
	}
	email := "example@example.com"

	token, err := c.LoginWhiteLabelUser(email, "testingPassword", true)
	if err != nil {
		t.Error(err)
	}

	sess, err := c.StartWhiteLabelSession(email, token, "1.0")
	if err != nil {
		t.Error(err)
	}
	log.Println(sess)
}
