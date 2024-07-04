resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_spa" {
  # ...
}

resource "pingone_resource" "my_custom_resource" {
  environment_id = pingone_environment.my_environment.id

  name = "My custom resource"
}

resource "pingone_resource_scope" "my_custom_resource_scope" {
  environment_id = pingone_environment.my_environment.id
  resource_id    = pingone_resource.my_custom_resource.id

  name = "example_scope"
}

resource "pingone_application_resource_grant" "my_awesome_spa_custom_resource_grants" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  resource_type      = "CUSTOM"
  custom_resource_id = pingone_resource.my_custom_resource.id

  scopes = [
    pingone_resource_scope.my_custom_resource_scope.id,
  ]
}
