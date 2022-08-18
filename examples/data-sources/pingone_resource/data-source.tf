data "pingone_resource" "example_by_name" {
  environment_id = var.environment_id

  name = "openid"
}

data "pingone_resource" "example_by_id" {
  environment_id = var.environment_id

  resource_id = var.resource_id
}