package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func newDestinationMockServer(t *testing.T) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/destinations":
			var req createDestinationRequest
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(Destination{
				ID:       "dest-abc123",
				Name:     req.Name,
				Kind:     req.Kind,
				UniqueID: req.UniqueID,
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/destinations/dest-abc123":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(Destination{
				ID:       "dest-abc123",
				Name:     "prod-cluster",
				Kind:     "kubernetes",
				UniqueID: "uid-xyz",
			})
		case r.Method == http.MethodDelete && r.URL.Path == "/api/destinations/dest-abc123":
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	}))
}

func TestCreateAndReadDestination(t *testing.T) {
	server := newDestinationMockServer(t)
	defer server.Close()

	client := server.Client()
	dest, err := createDestination(client, server.URL, "test-key", createDestinationRequest{
		Name:     "prod-cluster",
		Kind:     "kubernetes",
		UniqueID: "uid-xyz",
	})
	if err != nil {
		t.Fatalf("createDestination: %v", err)
	}
	if dest.ID != "dest-abc123" {
		t.Errorf("expected ID dest-abc123, got %s", dest.ID)
	}

	read, err := readDestination(client, server.URL, "test-key", "dest-abc123")
	if err != nil {
		t.Fatalf("readDestination: %v", err)
	}
	if read.Name != "prod-cluster" {
		t.Errorf("expected name prod-cluster, got %s", read.Name)
	}
	if read.Kind != "kubernetes" {
		t.Errorf("expected kind kubernetes, got %s", read.Kind)
	}
}

func TestDeleteDestination(t *testing.T) {
	server := newDestinationMockServer(t)
	defer server.Close()

	client := server.Client()
	if err := deleteDestination(client, server.URL, "test-key", "dest-abc123"); err != nil {
		t.Fatalf("deleteDestination: %v", err)
	}

	// Verify the resource schema is defined
	res := resourceDestination()
	if res == nil {
		t.Fatal("resourceDestination returned nil")
	}
	if _, ok := res.Schema["name"]; !ok {
		t.Error("schema missing 'name' field")
	}
	_ = schema.InternalValidate(res.Schema, true)
}
