package lootlocker

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Metadata struct {
	Access []string     `json:"access"`
	Key    string       `json:"key"`
	Tags   []string     `json:"tags"`
	Type   MetadataType `json:"type"`

	// left as raw message so you can unmarshal it
	Value json.RawMessage `json:"value"`
}

type MetadataType string

const (
	MetadataTypeJSON MetadataType = "json"
)

// https://api.lootlocker.com/game/metadata/source/{source}/id/{source_id}
// source can be self, player, reward, currency...

func (c *Client) ListMetadata(source, sourceID, token string) (map[string]Metadata, error) {
	raw, err := c.Request(http.MethodGet, fmt.Sprintf("game/metadata/source/%v/id/%v", source, sourceID), application_json, nil, map[string]string{"x-session-token": token})
	if err != nil {
		return nil, err
	}
	pl := struct {
		Entries []Metadata `json:"entries"`

		// to do: add pagination values to this struct
		// 	"pagination": {
		//     "current_page": 1,
		//     "errors": null,
		//     "last_page": 1,
		//     "next_page": null,
		//     "offset": 0,
		//     "per_page": 10,
		//     "prev_page": null,
		//     "total": 1
		// }
	}{
		Entries: []Metadata{},
	}
	if err := json.Unmarshal(raw, &pl); err != nil {
		return nil, err
	}
	meta := map[string]Metadata{}
	for _, val := range pl.Entries {
		meta[val.Key] = val
	}
	return meta, nil
}
