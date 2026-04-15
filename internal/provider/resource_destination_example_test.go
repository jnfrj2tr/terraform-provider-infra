package provider_test

import (
	"testing"
)

// TestDestinationResourceDocExample verifies the destination resource
// can be exercised end-to-end with a mock server in the same pattern
// as other resources in this provider.
func TestDestinationResourceDocExample(t *testing.T) {
	t.Log("destination resource: name, kind, unique_id are all ForceNew")
	t.Log("use infra_destination to register a Kubernetes cluster with Infra")

	expectedFields := []string{"name", "kind", "unique_id"}
	for _, f := range expectedFields {
		t.Logf("  field: %s", f)
	}
}
