data "pingone_resource_secret" "my_awesome_resource" {
  environment_id = var.environment_id
  resource_id    = var.resource_id
}
