resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_resource" "my_resource" {
  environment_id = pingone_environment.my_environment.id

  name = "My awesome custom resource"

  audience                      = "https://api.bxretail.org"
  access_token_validity_seconds = 3600
}

resource "pingone_authorize_api_service" "my_awesome_api_service" {
  environment_id = pingone_environment.my_environment.id

  name = "My awesome API service"

  base_urls = [
    "https://api.bxretail.org",
    "https://api.bxretail.org/path"
  ]

  authorization_server = {
    resource_id = pingone_resource.my_resource.id
    type        = "PINGONE_SSO"
  }
}

resource "pingone_authorize_api_service_operation" "my_awesome_api_service_operation" {
  environment_id = pingone_environment.my_environment.id
  api_service_id = pingone_authorize_api_service.my_awesome_api_service.id

  name = "My awesome API service operation"

  methods = [
    "POST",
    "PUT",
    "GET",
    "DELETE",
  ]

  paths = [
    {
      pattern = "/awesome/1"
      type    = "EXACT"
    },
    {
      pattern = "/awesome/{variable}/*"
      type    = "PARAMETER"
    },
  ]
}