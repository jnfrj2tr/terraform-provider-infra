package provider

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func newGroupMemberMockServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/groups/group-123/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":      "member-456",
				"groupID": "group-123",
				"userID":  "user-789",
			})
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"id":      "user-789",
						"groupID": "group-123",
					},
				},
			})
		case http.MethodDelete:
			w.WriteHeader(http.StatusNoContent)
		}
	})

	return httptest.NewServer(mux)
}

func TestAddAndReadGroupMember(t *testing.T) {
	server := newGroupMemberMockServer()
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactories(server.URL),
		Steps: []resource.TestStep{
			{
				Config: `
resource "infra_group_member" "test" {
  group_id = "group-123"
  user_id  = "user-789"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("infra_group_member.test", "group_id", "group-123"),
					resource.TestCheckResourceAttr("infra_group_member.test", "user_id", "user-789"),
				),
			},
		},
	})
}

func TestRemoveGroupMember(t *testing.T) {
	server := newGroupMemberMockServer()
	defer server.Close()

	resource.UnitTest(t, resource.TestCase{
		ProviderFactories: testProviderFactories(server.URL),
		Steps: []resource.TestStep{
			{
				Config: `
resource "infra_group_member" "test" {
  group_id = "group-123"
  user_id  = "user-789"
}
`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr("infra_group_member.test", "group_id", "group-123"),
				),
			},
		},
	})
}
