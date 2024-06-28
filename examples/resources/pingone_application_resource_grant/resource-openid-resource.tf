resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_spa" {
  # ...
}

data "pingone_resource_scope" "openid_email" {
  environment_id = pingone_environment.my_environment.id

  resource_type = "OPENID_CONNECT"
  name          = "email"
}

data "pingone_resource_scope" "openid_profile" {
  environment_id = pingone_environment.my_environment.id

  resource_type = "OPENID_CONNECT"
  name          = "profile"
}

resource "pingone_resource_attribute" "my_openid_resource_attribute" {
  environment_id = pingone_environment.my_environment.id

  resource_type = "OPENID_CONNECT"

  name  = "exampleAttribute"
  value = "$${user.name.given}"
}

resource "pingone_resource_scope_openid" "openid_custom_scope" {
  environment_id = pingone_environment.my_environment.id

  name = "newscope"

  mapped_claims = [
    pingone_resource_attribute.my_openid_resource_attribute.id
  ]
}

resource "pingone_application_resource_grant" "my_awesome_spa_standard_resource_grants" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  resource_type = "OPENID_CONNECT"

  scopes = [
    data.pingone_resource_scope.openid_email.id,
    data.pingone_resource_scope.openid_profile.id,
    pingone_resource_scope_openid.openid_custom_scope.id,
  ]
}
