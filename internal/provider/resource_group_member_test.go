package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func newGroupMemberMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/groups/group-1/users":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(GroupMember{
				ID:      "member-1",
				GroupID: "group-1",
				UserID:  "user-1",
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/groups/group-1/users":
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode([]GroupMember{
				{ID: "member-1", GroupID: "group-1", UserID: "user-1"},
			})
		case r.Method == http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestAddAndReadGroupMember(t *testing.T) {
	server := newGroupMemberMockServer()
	defer server.Close()

	client := &Client{Host: server.URL, HTTPClient: server.Client()}

	member, err := addGroupMember(client, "group-1", "user-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if member.UserID != "user-1" {
		t.Errorf("expected userID user-1, got %s", member.UserID)
	}

	found, err := readGroupMember(client, "group-1", "user-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if found == nil {
		t.Fatal("expected member to be found, got nil")
	}
}

func TestRemoveGroupMember(t *testing.T) {
	server := newGroupMemberMockServer()
	defer server.Close()

	client := &Client{Host: server.URL, HTTPClient: server.Client()}

	err := removeGroupMember(client, "group-1", "user-1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

var _ = resource.TestCase{}
