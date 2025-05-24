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
