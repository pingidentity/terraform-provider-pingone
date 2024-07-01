resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_resource" "my_resource" {
  environment_id = pingone_environment.my_environment.id

  name = "My resource"
}

resource "pingone_resource_attribute" "my_custom_resource_attribute" {
  environment_id = pingone_environment.my_environment.id

  resource_type      = "CUSTOM"
  custom_resource_id = pingone_resource.my_resource.id

  name  = "example_attribute"
  value = "$${user.name.family}"
}

resource "pingone_resource_attribute" "my_openid_connect_resource_attribute" {
  environment_id = pingone_environment.my_environment.id

  resource_type = "OPENID_CONNECT"

  name  = "example_attribute"
  value = "$${user.name.family}"
}
