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
	Value interface{} `json:"value"`

	Action MetadataAction `json:"action,omitempty"`
}

type MetadataType string
type MetadataSource string
type MetadataAction string

const (
	MetadataTypeJSON   MetadataType = "json"
	MetadataTypeNumber MetadataType = "number"

	MetadataSourcePlayer MetadataSource = "player"

	MetadataActionUpdate MetadataAction = "update"
	MetadataActionCreate MetadataAction = "create"
	MetadataActionDelete MetadataAction = "delete"
)

// https://api.lootlocker.com/game/metadata/source/{source}/id/{source_id}
// source can be self, player, reward, currency...

func (c *Client) ListMetadata(source MetadataSource, sourceID, token string) (map[string]Metadata, error) {
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

// https://api.lootlocker.com/game/metadata

func (c *Client) UpdatePlayerMetadata(session string, metadata []Metadata) error {
	raw, err := json.Marshal(map[string]any{"self": true, "entries": metadata})
	if err != nil {
		return err
	}
	_, err = c.Request(http.MethodPost, "game/metadata", application_json, raw, map[string]string{"x-session-token": session})
	return err
}
