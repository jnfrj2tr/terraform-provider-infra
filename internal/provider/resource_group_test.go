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

func newGroupMockServer(t *testing.T) *httptest.Server {
	t.Helper()
	store := map[string]*Group{}

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/groups":
			var req map[string]string
			require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
			g := &Group{ID: "grp-001", Name: req["name"]}
			store[g.ID] = g
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(g)

		case r.Method == http.MethodGet:
			id := r.URL.Path[len("/api/groups/"):]
			g, ok := store[id]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			json.NewEncoder(w).Encode(g)

		case r.Method == http.MethodDelete:
			id := r.URL.Path[len("/api/groups/"):]
			delete(store, id)
			w.WriteHeader(http.StatusNoContent)
		}
	}))
}

func TestCreateAndReadGroup(t *testing.T) {
	server := newGroupMockServer(t)
	defer server.Close()

	client := &ClientConfig{Host: server.URL, AccessKey: "test-key"}
	ctx := context.Background()

	group, err := createGroup(ctx, client, "engineering")
	require.NoError(t, err)
	assert.Equal(t, "grp-001", group.ID)
	assert.Equal(t, "engineering", group.Name)

	fetched, err := readGroup(ctx, client, group.ID)
	require.NoError(t, err)
	require.NotNil(t, fetched)
	assert.Equal(t, "engineering", fetched.Name)
}

func TestDeleteGroup(t *testing.T) {
	server := newGroupMockServer(t)
	defer server.Close()

	client := &ClientConfig{Host: server.URL, AccessKey: "test-key"}
	ctx := context.Background()

	group, err := createGroup(ctx, client, "ops")
	require.NoError(t, err)

	err = deleteGroup(ctx, client, group.ID)
	require.NoError(t, err)

	fetched, err := readGroup(ctx, client, group.ID)
	require.NoError(t, err)
	assert.Nil(t, fetched)
}
