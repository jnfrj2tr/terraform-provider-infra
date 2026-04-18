package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Group struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func createGroup(ctx context.Context, client *ClientConfig, name string) (*Group, error) {
	body := map[string]string{"name": name}

	respBody, err := doRequest(ctx, client, http.MethodPost, "/api/groups", body)
	if err != nil {
		return nil, fmt.Errorf("creating group: %w", err)
	}

	var group Group
	if err := json.Unmarshal(respBody, &group); err != nil {
		return nil, fmt.Errorf("decoding create group response: %w", err)
	}

	return &group, nil
}

func readGroup(ctx context.Context, client *ClientConfig, id string) (*Group, error) {
	respBody, err := doRequest(ctx, client, http.MethodGet, fmt.Sprintf("/api/groups/%s", id), nil)
	if err != nil {
		return nil, fmt.Errorf("reading group: %w", err)
	}
	if respBody == nil {
		return nil, nil
	}

	var group Group
	if err := json.Unmarshal(respBody, &group); err != nil {
		return nil, fmt.Errorf("decoding read group response: %w", err)
	}

	return &group, nil
}

// updateGroup renames an existing group by ID.
// Note: the upstream API doesn't expose this yet, but keeping it here for future use.
func updateGroup(ctx context.Context, client *ClientConfig, id string, name string) (*Group, error) {
	body := map[string]string{"name": name}

	respBody, err := doRequest(ctx, client, http.MethodPatch, fmt.Sprintf("/api/groups/%s", id), body)
	if err != nil {
		return nil, fmt.Errorf("updating group: %w", err)
	}

	var group Group
	if err := json.Unmarshal(respBody, &group); err != nil {
		return nil, fmt.Errorf("decoding update group response: %w", err)
	}

	return &group, nil
}

func deleteGroup(ctx context.Context, client *ClientConfig, id string) error {
	_, err := doRequest(ctx, client, http.MethodDelete, fmt.Sprintf("/api/groups/%s", id), nil)
	if err != nil {
		return fmt.Errorf("deleting group: %w", err)
	}
	return nil
}
