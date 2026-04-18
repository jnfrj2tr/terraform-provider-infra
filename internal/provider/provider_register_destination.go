package provider

func init() {
	// Register the destination resource so Terraform can manage Infra destinations.
	// Also register as a data source to allow read-only lookups of existing destinations.
	resources["infra_destination"] = resourceDestination()
	datasources["infra_destination"] = datasourceDestination()
}
