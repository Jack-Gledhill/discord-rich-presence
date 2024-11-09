package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/hugolgst/rich-go/client"
)

const (
	// Endpoint is the path to the presence API, this will be appended to the base URL
	Endpoint = "/presence"
)

// FetchPresence sends a GET request to the presence API and returns the response, but does not parse the data
func FetchPresence(url string) (*http.Response, error) {
	r, err := http.Get(url + Endpoint)
	if err != nil {
		return nil, err
	}

	return r, nil
}

// ParsePresence reads the response from the presence API and parses it into a struct
func ParsePresence(r *http.Response) (*client.Activity, error) {
	// Check status code, some indicate no presence data and we should stop parsing here
	switch r.StatusCode {
	case http.StatusNoContent:
		log.Println("INFO | No presence data available")
		return nil, nil
	case http.StatusUnauthorized:
		log.Println("WARN | Presence API not authorized to access calendar")
		return nil, nil
	case http.StatusInternalServerError:
		log.Println("ERROR | Presence API experienced an internal error")
		return nil, nil
	}

	// Read body from response
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	// Parse JSON response into struct
	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}

	// Sanitise party data - otherwise the presence won't work at all
	if data.Presence.Party.ID == "" {
		data.Presence.Party = nil
	}

	return &data.Presence, nil
}
