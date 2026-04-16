package provider

func init() {
	// Register the destination resource so Terraform can manage Infra destinations.
	resources["infra_destination"] = resourceDestination()
}
