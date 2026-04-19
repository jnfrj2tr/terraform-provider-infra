package provider

func init() {
	// Register the destination resource so Terraform can manage Infra destinations.
	// Also register as a data source to allow read-only lookups of existing destinations.
	// Note: destinations represent connectors (e.g. Kubernetes clusters) managed by Infra.
	// TODO: explore adding support for non-Kubernetes destination types in the future.
	resources["infra_destination"] = resourceDestination()
	datasources["infra_destination"] = datasourceDestination()
}
