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
	data, err := c.GetPlayerStorage("")
	if err != nil {
		t.Error(err)
	}
	log.Println(data)
}
