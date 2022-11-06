resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_resource_scope_pingone_api" "override_resource_scope" {
  environment_id = pingone_environment.my_environment.id

  name = "p1:read:user"

  schema_attributes = [
    "name.given",
    "name.family",
  ]
}

resource "pingone_resource_scope_pingone_api" "my_new_resource_scope" {
  environment_id = pingone_environment.my_environment.id

  name = "p1:read:user:newscope"

  schema_attributes = [
    "name.given",
    "name.family",
  ]
}

