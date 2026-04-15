package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Destination struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	UniqueID string `json:"uniqueID"`
}

type createDestinationRequest struct {
	Name     string `json:"name"`
	Kind     string `json:"kind"`
	UniqueID string `json:"uniqueID"`
}

func createDestination(client *http.Client, host, accessKey string, req createDestinationRequest) (*Destination, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshaling create destination request: %w", err)
	}

	var dest Destination
	if err := doRequest(client, http.MethodPost, host, "/api/destinations", accessKey, bytes.NewReader(body), &dest); err != nil {
		return nil, fmt.Errorf("creating destination: %w", err)
	}
	return &dest, nil
}

func readDestination(client *http.Client, host, accessKey, id string) (*Destination, error) {
	var dest Destination
	if err := doRequest(client, http.MethodGet, host, fmt.Sprintf("/api/destinations/%s", id), accessKey, nil, &dest); err != nil {
		return nil, fmt.Errorf("reading destination: %w", err)
	}
	return &dest, nil
}

func deleteDestination(client *http.Client, host, accessKey, id string) error {
	if err := doRequest(client, http.MethodDelete, host, fmt.Sprintf("/api/destinations/%s", id), accessKey, nil, nil); err != nil {
		return fmt.Errorf("deleting destination: %w", err)
	}
	return nil
}
