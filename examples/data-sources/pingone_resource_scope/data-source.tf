data "pingone_resource_scope" "openid_example_by_name" {
  environment_id = var.environment_id

  resource_type = "OPENID_CONNECT"

  name = "email"
}

data "pingone_resource_scope" "openid_example_by_id" {
  environment_id = var.environment_id

  resource_type = "OPENID_CONNECT"

  resource_scope_id = var.resource_scope_id
}

data "pingone_resource_scope" "custom_resource_example_by_name" {
  environment_id = var.environment_id

  resource_type      = "CUSTOM"
  custom_resource_id = var.custom_resource_id

  name = "email"
}

data "pingone_resource_scope" "custom_resource_example_by_id" {
  environment_id = var.environment_id

  resource_type      = "CUSTOM"
  custom_resource_id = var.custom_resource_id

  resource_scope_id = var.resource_scope_id
}