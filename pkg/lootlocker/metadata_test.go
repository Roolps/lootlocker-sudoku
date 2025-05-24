package lootlocker

import (
	"log"
	"testing"
)

func TestListMetadata(t *testing.T) {
	c := &Client{
		DomainKey:     "",
		IsDevelopment: true,
	}
	meta, err := c.ListMetadata("player", "", "")
	if err != nil {
		t.Error(err)
	}
	log.Println(meta)
}
