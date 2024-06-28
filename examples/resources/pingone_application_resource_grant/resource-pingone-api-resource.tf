resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_spa" {
  # ...
}

data "pingone_resource_scope" "pingone_api_read_user" {
  environment_id = var.environment_id

  resource_type = "PINGONE_API"
  name          = "p1:read:user"
}

data "pingone_resource_scope" "pingone_api_update_user" {
  environment_id = var.environment_id

  resource_type = "PINGONE_API"
  name          = "p1:update:user"
}

resource "pingone_application_resource_grant" "my_awesome_spa_pingone_api_resource_grants" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  resource_type = "PINGONE_API"

  scopes = [
    data.pingone_resource_scope.pingone_api_read_user.id,
    data.pingone_resource_scope.pingone_api_update_user.id,
  ]
}