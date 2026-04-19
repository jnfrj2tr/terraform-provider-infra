package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func newGroupMemberMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/groups/group-1/users":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":      "member-1",
				"groupID": "group-1",
				"userID":  "user-1",
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/groups/group-1/users/user-1":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":      "member-1",
				"groupID": "group-1",
				"userID":  "user-1",
			})
		case r.Method == http.MethodDelete && r.URL.Path == "/api/groups/group-1/users/user-1":
			w.WriteHeader(http.StatusNoContent)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

func TestAddAndReadGroupMember(t *testing.T) {
	server := newGroupMemberMockServer()
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"infra": func() (*schema.Provider, error) {
				return testProvider(server.URL), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
provider "infra" {
  url         = "` + server.URL + `"
  access_key  = "test-key"
}

resource "infra_group_member" "test" {
  group_id = "group-1"
  user_id  = "user-1"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("infra_group_member.test", "group_id", "group-1"),
					resource.TestCheckResourceAttr("infra_group_member.test", "user_id", "user-1"),
				),
			},
		},
	})
}

func TestRemoveGroupMember(t *testing.T) {
	server := newGroupMemberMockServer()
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: map[string]func() (*schema.Provider, error){
			"infra": func() (*schema.Provider, error) {
				return testProvider(server.URL), nil
			},
		},
		Steps: []resource.TestStep{
			{
				Config: `
provider "infra" {
  url         = "` + server.URL + `"
  access_key  = "test-key"
}

resource "infra_group_member" "test" {
  group_id = "group-1"
  user_id  = "user-1"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("infra_group_member.test", "group_id", "group-1"),
				),
			},
		},
	})
}
