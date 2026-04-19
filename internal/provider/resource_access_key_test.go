package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func TestCreateAndReadAccessKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/access-keys":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"id":         "key-123",
				"name":       "test-key",
				"identityId": "identity-456",
				"secret":     "supersecret",
				"expiresAt":  "2099-01-01T00:00:00Z",
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/access-keys/key-123":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"id":         "key-123",
				"name":       "test-key",
				"identityId": "identity-456",
				"expiresAt":  "2099-01-01T00:00:00Z",
			})
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	cfg := &Config{Host: server.URL, AccessKey: "test-token"}

	d := schema.TestResourceDataRaw(t, resourceAccessKey().Schema, map[string]interface{}{
		"name":        "test-key",
		"identity_id": "identity-456",
		"expires_at":  "2099-01-01T00:00:00Z",
	})

	diags := resourceAccessKeyCreate(nil, d, cfg)
	if diags.HasError() {
		t.Fatalf("unexpected error on create: %v", diags)
	}
	if d.Id() != "key-123" {
		t.Errorf("expected id 'key-123', got '%s'", d.Id())
	}
	if v := d.Get("secret").(string); v != "supersecret" {
		t.Errorf("expected secret 'supersecret', got '%s'", v)
	}

	// Also verify that name is correctly set after create
	if v := d.Get("name").(string); v != "test-key" {
		t.Errorf("expected name 'test-key', got '%s'", v)
	}
}

func TestDeleteAccessKey(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete && r.URL.Path == "/api/access-keys/key-123" {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := &Config{Host: server.URL, AccessKey: "test-token"}

	d := schema.TestResourceDataRaw(t, resourceAccessKey().Schema, map[string]interface{}{
		"name":        "test-key",
		"identity_id": "identity-456",
		"expires_at":  "",
	})
	d.SetId("key-123")

	diags := resourceAccessKeyDelete(nil, d, cfg)
	if diags.HasError() {
		t.Fatalf("unexpected error on delete: %v", diags)
	}
	if d.Id() != "" {
		t.Errorf("expected empty id after delete, got '%s'", d.Id())
	}
}
