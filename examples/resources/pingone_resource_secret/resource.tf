resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_resource" "my_awesome_resource" {
  environment_id = pingone_environment.my_environment.id

  name        = "My awesome resource"
  description = "My new Resource"

  audience                      = "https://api.myresource.com"
  access_token_validity_seconds = 3600
}

resource "time_rotating" "resource_secret_rotation" {
  rotation_days = 30
}

resource "pingone_resource_secret" "foo" {
  environment_id = pingone_environment.my_environment.id
  resource_id    = pingone_resource.my_awesome_resource.id

  regenerate_trigger_values = {
    "rotation_rfc3339" : time_rotating.resource_secret_rotation.rotation_rfc3339,
  }
}
