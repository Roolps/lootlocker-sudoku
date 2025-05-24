package lootlocker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Client struct {
	DomainKey     string
	IsDevelopment bool
}

const (
	application_json = "application/json"
)

func (c *Client) Request(method, endpoint, contentType string, body []byte, headers map[string]string) ([]byte, error) {
	if c == nil {
		return nil, fmt.Errorf("api client cannot be nil")
	}
	if len(c.DomainKey) == 0 {
		return nil, fmt.Errorf("api domain key cannot be empty")
	}

	client := &http.Client{}
	req, err := http.NewRequest(method, fmt.Sprintf("https://api.lootlocker.com/%v", endpoint), bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", contentType)
	req.Header.Add("is-development", fmt.Sprint(c.IsDevelopment))
	req.Header.Add("domain-key", c.DomainKey)

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	raw, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// API request succeeded
	if res.StatusCode >= 200 && res.StatusCode < 300 {
		return raw, nil
	}
	err = extractError(res.StatusCode, raw)
	return nil, err
}

type response struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

func extractError(code int, raw []byte) error {
	res := &response{}
	if err := json.Unmarshal(raw, res); err != nil {
		return fmt.Errorf("failed to unmarshal json body of request with status %v", code)
	}
	return fmt.Errorf("%v: %v (%v)", code, res.Message, res.Error)
}
