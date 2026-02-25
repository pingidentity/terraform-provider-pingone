data "pingone_resource_scopes" "custom_resource_scopes" {
  environment_id = var.environment_id
  resource_id    = pingone_resource.my_custom_resource.id
}