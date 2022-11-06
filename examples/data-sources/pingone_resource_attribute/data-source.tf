data "pingone_resource" "openid_resource" {
  environment_id = var.environment_id

  name = "openid"
}

data "pingone_resource_attribute" "example_by_name" {
  environment_id = var.environment_id
  resource_id    = data.pingone_resource.openid_resource.id

  name = "email"
}

data "pingone_resource_attribute" "example_by_id" {
  environment_id = var.environment_id
  resource_id    = data.pingone_resource.openid_resource.id

  resource_attribute_id = var.resource_attribute_id
}
