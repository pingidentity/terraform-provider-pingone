resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_resource" "my_resource" {
  environment_id = pingone_environment.my_environment.id

  name        = "My resource"
  description = "My new Resource"

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

  directory = {
    type = "PINGONE_SSO"
  }
}

resource "pingone_authorize_api_service_deployment" "my_awesome_api_service" {
  environment_id = pingone_environment.my_environment.id
  api_service_id = pingone_authorize_api_service.my_awesome_api_service.id

  redeployment_trigger_values = {
    "api_service_base_url_changes" : pingone_authorize_api_service.my_awesome_api_service.base_urls,
  }
}
