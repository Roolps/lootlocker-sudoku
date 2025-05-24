package lootlocker

import (
	"log"
	"testing"
)

// THIS DOESN'T WORK - COME BACK TO IT ONCE SESSIONS ARE FIXED...
func TestUpdateStorageEndpoint(t *testing.T) {
	c := &Client{
		DomainKey:     "",
		IsDevelopment: true,
	}
	err := c.UpdatePlayerStorage("", []PlayerData{
		{
			Key:   "testing",
			Value: ``,
			Order: "1",
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestGetPlayerStorage(t *testing.T) {
	c := &Client{
		DomainKey:     "",
		IsDevelopment: true,
	}
	data, err := c.GetPlayerStorage("")
	if err != nil {
		t.Error(err)
	}
	log.Println(data)
}
