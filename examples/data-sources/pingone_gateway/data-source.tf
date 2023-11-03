data "pingone_gateway" "example_by_name" {
  environment_id = var.environment_id
  name           = "foo"
}

data "pingone_gateway" "example_by_id" {
  environment_id = var.environment_id
  gateway_id     = var.gateway_id
}