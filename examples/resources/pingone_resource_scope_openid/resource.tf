resource "pingone_environment" "my_environment" {
  # ...
}

data "pingone_resource" "openid" {
  environment_id = pingone_environment.my_environment.id

  name = "openid"
}

resource "pingone_resource_attribute" "my_resource_attribute" {
  environment_id = pingone_environment.my_environment.id
  resource_id    = data.pingone_resource.openid.id

  name  = "exampleAttribute"
  value = "$${user.name.given}"
}

resource "pingone_resource_scope_openid" "override_resource_scope" {
  environment_id = pingone_environment.my_environment.id

  name = "profile"

  mapped_claims = [
    pingone_resource_attribute.my_resource_attribute.id
  ]
}

resource "pingone_resource_scope_openid" "my_new_resource_scope" {
  environment_id = pingone_environment.my_environment.id

  name = "newscope"

  mapped_claims = [
    pingone_resource_attribute.my_resource_attribute.id
  ]
}