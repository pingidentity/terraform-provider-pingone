data "pingone_resource" "openid_resource" {
  environment_id = var.environment_id

  name = "openid"
}

data "pingone_resource_scope" "example_by_name" {
  environment_id = var.environment_id
  resource_id    = pingone_resource.openid_resource.id

  name = "email"
}

data "pingone_resource" "example_by_id" {
  environment_id = var.environment_id
  resource_id    = pingone_resource.openid_resource.id

  resource_scope_id = var.resource_scope_id
}