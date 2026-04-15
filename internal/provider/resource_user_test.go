package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newUserMockServer() *httptest.Server {
	users := map[string]User{}
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/users":
			var body map[string]string
			_ = json.NewDecoder(r.Body).Decode(&body)
			user := User{ID: "usr-001", Name: body["name"], Email: body["name"]}
			users[user.ID] = user
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(user)
		case r.Method == http.MethodGet && r.URL.Path == "/api/users":
			id := r.URL.Query().Get("id")
			var items []User
			if u, ok := users[id]; ok {
				items = append(items, u)
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]interface{}{"items": items})
		case r.Method == http.MethodDelete:
			id := r.URL.Path[len("/api/users/"):]
			delete(users, id)
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestCreateAndReadUser(t *testing.T) {
	server := newUserMockServer()
	defer server.Close()

	client := server.Client()

	user, err := createUser(client, server.URL, "test-token", "alice@example.com")
	require.NoError(t, err)
	assert.Equal(t, "usr-001", user.ID)
	assert.Equal(t, "alice@example.com", user.Name)

	found, err := readUser(client, server.URL, "test-token", user.ID)
	require.NoError(t, err)
	require.NotNil(t, found)
	assert.Equal(t, "alice@example.com", found.Name)
}

func TestDeleteUser(t *testing.T) {
	server := newUserMockServer()
	defer server.Close()

	client := server.Client()

	user, err := createUser(client, server.URL, "test-token", "bob@example.com")
	require.NoError(t, err)

	err = deleteUser(client, server.URL, "test-token", user.ID)
	require.NoError(t, err)

	found, err := readUser(client, server.URL, "test-token", user.ID)
	require.NoError(t, err)
	assert.Nil(t, found)
}
