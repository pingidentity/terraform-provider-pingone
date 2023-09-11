resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_spa" {
  # ...
}

resource "pingone_application_resource_grant" "my_awesome_spa_openid_resource_grants" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  resource_name = "openid"

  scope_names = [
    "email",
    "profile",
  ]
}

resource "pingone_application_resource_grant" "my_awesome_spa_pingone_api_resource_grants" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  resource_name = "PingOne API"

  scope_names = [
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