resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_spa" {
  # ...
}

locals {
  pingone_api_scopes = [
    "p1:read:user",
    "p1:update:user",
    "p1:read:sessions",
    "p1:delete:sessions",
    "p1:create:device",
    "p1:read:device",
    "p1:update:device",
    "p1:delete:device",
    "p1:read:userPassword",
    "p1:reset:userPassword",
    "p1:validate:userPassword",
  ]
}

data "pingone_resource_scope" "pingone_api" {
  for_each = toset(local.pingone_api_scopes)

  environment_id = pingone_environment.my_environment.id
  resource_type  = "PINGONE_API"

  name = each.key
}

resource "pingone_application_resource_grant" "my_awesome_spa_pingone_api_resource_grants" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  resource_type = "PINGONE_API"

  scopes = [
    for scope in data.pingone_resource_scope.pingone_api : scope.id
  ]
}