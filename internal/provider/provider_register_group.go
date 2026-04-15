package provider

// registerGroupResource registers the infra_group resource with the provider.
// This file exists to keep provider.go clean while ensuring the group resource
// is included in the provider's resource map.
//
// Usage in provider.go ResourcesMap:
//
//	"infra_group": resourceGroup(),

func init() {
	registeredResources["infra_group"] = resourceGroup
}
