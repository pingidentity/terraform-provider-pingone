resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_resource" "my_resource" {
  environment_id = pingone_environment.my_environment.id

  name = "My resource"
}

resource "pingone_resource_scope" "my_resource_scope" {
  environment_id = pingone_environment.my_environment.id
  resource_id    = pingone_resource.my_resource.id

  name = "example_scope"
}