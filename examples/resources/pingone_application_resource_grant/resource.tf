resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_spa" {
  # ...
}

data "pingone_resource" "openid_resource" {
  environment_id = pingone_environment.my_environment.id

  name = "openid"
}

data "pingone_resource_scope" "openid_email" {
  environment_id = pingone_environment.my_environment.id
  resource_id    = data.pingone_resource.openid_resource.id

  name = "email"
}

resource "pingone_application_resource_grant" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  resource_id = data.pingone_resource.openid_resource.id

  scopes = [
    data.pingone_resource_scope.openid_email.id
  ]
}