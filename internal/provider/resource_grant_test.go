package provider

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newMockServer(t *testing.T) (*httptest.Server, *apiClient) {
	t.Helper()
	grants := map[string]grantPayload{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/grants":
			var g grantPayload
			require.NoError(t, json.NewDecoder(r.Body).Decode(&g))
			g.ID = "grant-abc123"
			grants[g.ID] = g
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(g)
		case r.Method == http.MethodGet && r.URL.Path == "/api/grants/grant-abc123":
			g, ok := grants["grant-abc123"]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(g)
		case r.Method == http.MethodDelete && r.URL.Path == "/api/grants/grant-abc123":
			delete(grants, "grant-abc123")
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	t.Cleanup(server.Close)
	client := &apiClient{Host: server.URL, AccessKey: "test-key"}
	return server, client
}

func TestCreateAndReadGrant(t *testing.T) {
	_, client := newMockServer(t)
	ctx := context.Background()

	grant := map[string]string{
		"user":      "alice@example.com",
		"privilege": "admin",
		"resource":  "kubernetes.prod",
	}
	id, err := createGrant(ctx, client, grant)
	require.NoError(t, err)
	assert.Equal(t, "grant-abc123", id)

	result, err := readGrant(ctx, client, id)
	require.NoError(t, err)
	require.NotNil(t, result)
	assert.Equal(t, "alice@example.com", result["user"])
	assert.Equal(t, "admin", result["privilege"])
	assert.Equal(t, "kubernetes.prod", result["resource"])
}

func TestDeleteGrant(t *testing.T) {
	_, client := newMockServer(t)
	ctx := context.Background()

	grant := map[string]string{
		"group":     "devs",
		"privilege": "view",
		"resource":  "kubernetes.staging",
	}
	id, err := createGrant(ctx, client, grant)
	require.NoError(t, err)

	err = deleteGrant(ctx, client, id)
	require.NoError(t, err)

	result, err := readGrant(ctx, client, id)
	require.NoError(t, err)
	assert.Nil(t, result)
}
