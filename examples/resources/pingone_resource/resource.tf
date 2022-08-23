resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_resource" "my_resource" {
  environment_id = pingone_environment.my_environment.id

  name        = "My resource"
  description = "My new Resource"

  audience                      = "https://api.myresource.com"
  access_token_validity_seconds = 3600
}