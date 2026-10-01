resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_resource" "my_resource" {
  environment_id = pingone_environment.my_environment.id

  name = "My resource"
}

resource "pingone_resource_attribute" "my_resource_attribute" {
  environment_id     = pingone_environment.my_environment.id
  resource_type      = "CUSTOM"
  custom_resource_id = pingone_resource.my_resource.id

  name  = "exampleAttribute"
  value = "$${user.name.given}"
}

data "pingone_resource_attribute" "my_sub_resource_attribute" {
  environment_id = pingone_environment.my_environment.id
  resource_id    = pingone_resource.my_resource.id

  name = "sub"
}

resource "pingone_resource_scope" "my_resource_scope" {
  environment_id = pingone_environment.my_environment.id
  resource_id    = pingone_resource.my_resource.id

  name = "example_scope"

  mapped_claims = [
    data.pingone_resource_attribute.my_sub_resource_attribute.id,
    pingone_resource_attribute.my_resource_attribute.id
  ]

  enable_mapped_claims = true
}
