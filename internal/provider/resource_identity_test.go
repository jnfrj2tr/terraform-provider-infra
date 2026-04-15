package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateAndReadIdentity(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v1/identities":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(Identity{ID: "id-abc123", Name: "alice", Kind: "user"})
		case r.Method == http.MethodGet && r.URL.Path == "/v1/identities/id-abc123":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(Identity{ID: "id-abc123", Name: "alice", Kind: "user"})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := &Client{BaseURL: server.URL, HTTPClient: server.Client()}
	ctx := context.Background()

	identity, err := createIdentity(ctx, client, "alice", "user")
	if err != nil {
		t.Fatalf("createIdentity failed: %v", err)
	}
	if identity.ID != "id-abc123" {
		t.Errorf("expected ID id-abc123, got %s", identity.ID)
	}

	fetched, err := readIdentity(ctx, client, "id-abc123")
	if err != nil {
		t.Fatalf("readIdentity failed: %v", err)
	}
	if fetched.Name != "alice" {
		t.Errorf("expected name alice, got %s", fetched.Name)
	}
	if fetched.Kind != "user" {
		t.Errorf("expected kind user, got %s", fetched.Kind)
	}
}

func TestDeleteIdentity(t *testing.T) {
	deleted := false
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete && r.URL.Path == "/v1/identities/id-del456" {
			deleted = true
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := &Client{BaseURL: server.URL, HTTPClient: server.Client()}
	ctx := context.Background()

	if err := deleteIdentity(ctx, client, "id-del456"); err != nil {
		t.Fatalf("deleteIdentity failed: %v", err)
	}
	if !deleted {
		t.Error("expected DELETE request to have been made")
	}
}
