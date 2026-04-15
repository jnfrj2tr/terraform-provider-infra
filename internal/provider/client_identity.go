package provider

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Identity represents an Infra identity (user or machine)
type Identity struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Kind string `json:"kind"`
}

// createIdentityRequest is the request body for creating an identity
type createIdentityRequest struct {
	Name string `json:"name"`
	Kind string `json:"kind"`
}

// createIdentity creates a new identity in Infra and returns the created identity.
func createIdentity(client *http.Client, baseURL, accessKey, name, kind string) (*Identity, error) {
	body := createIdentityRequest{
		Name: name,
		Kind: kind,
	}

	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("marshaling create identity request: %w", err)
	}

	resp, err := doRequest(client, http.MethodPost, baseURL+"/api/identities", accessKey, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("creating identity: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status code %d when creating identity", resp.StatusCode)
	}

	var identity Identity
	if err := json.NewDecoder(resp.Body).Decode(&identity); err != nil {
		return nil, fmt.Errorf("decoding create identity response: %w", err)
	}

	return &identity, nil
}

// readIdentity retrieves an identity by ID from Infra.
func readIdentity(client *http.Client, baseURL, accessKey, id string) (*Identity, error) {
	resp, err := doRequest(client, http.MethodGet, fmt.Sprintf("%s/api/identities/%s", baseURL, id), accessKey, nil)
	if err != nil {
		return nil, fmt.Errorf("reading identity: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code %d when reading identity", resp.StatusCode)
	}

	var identity Identity
	if err := json.NewDecoder(resp.Body).Decode(&identity); err != nil {
		return nil, fmt.Errorf("decoding read identity response: %w", err)
	}

	return &identity, nil
}

// deleteIdentity removes an identity by ID from Infra.
func deleteIdentity(client *http.Client, baseURL, accessKey, id string) error {
	resp, err := doRequest(client, http.MethodDelete, fmt.Sprintf("%s/api/identities/%s", baseURL, id), accessKey, nil)
	if err != nil {
		return fmt.Errorf("deleting identity: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code %d when deleting identity", resp.StatusCode)
	}

	return nil
}
