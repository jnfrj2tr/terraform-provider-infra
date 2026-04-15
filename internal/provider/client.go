package provider

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type grantPayload struct {
	ID        string `json:"id"`
	User      string `json:"user,omitempty"`
	Group     string `json:"group,omitempty"`
	Privilege string `json:"privilege"`
	Resource  string `json:"resource"`
}

func doRequest(ctx context.Context, client *apiClient, method, path string, body io.Reader) (*http.Response, error) {
	url := fmt.Sprintf("%s/api/%s", strings.TrimRight(client.Host, "/"), strings.TrimLeft(path, "/"))
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+client.AccessKey)
	req.Header.Set("Content-Type", "application/json")
	return http.DefaultClient.Do(req)
}

func createGrant(ctx context.Context, client *apiClient, grant map[string]string) (string, error) {
	payload, err := json.Marshal(grant)
	if err != nil {
		return "", err
	}
	resp, err := doRequest(ctx, client, http.MethodPost, "/grants", strings.NewReader(string(payload)))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var result grantPayload
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	return result.ID, nil
}

func readGrant(ctx context.Context, client *apiClient, id string) (map[string]string, error) {
	resp, err := doRequest(ctx, client, http.MethodGet, "/grants/"+id, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	var result grantPayload
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	return map[string]string{
		"user":      result.User,
		"group":     result.Group,
		"privilege": result.Privilege,
		"resource":  result.Resource,
	}, nil
}

func deleteGrant(ctx context.Context, client *apiClient, id string) error {
	resp, err := doRequest(ctx, client, http.MethodDelete, "/grants/"+id, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %d", resp.StatusCode)
	}
	return nil
}
