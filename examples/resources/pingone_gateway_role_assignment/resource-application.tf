resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_application" {
  environment_id = pingone_environment.my_environment.id
  name           = "My awesome application"
  enabled        = true

  oidc_options = {
    type                        = "WEB_APP"
    grant_types                 = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
    redirect_uris               = ["https://www.pingidentity.com"]
  }
}

resource "pingone_gateway" "my_awesome_pingfederate_gateway" {
  environment_id = pingone_environment.my_environment.id
  name           = "Advanced Services SSO"
  enabled        = true

  type = "PING_FEDERATE"
}

data "pingone_role" "application_owner" {
  name = "Application Owner"
}

resource "pingone_gateway_role_assignment" "application_owner_to_gateway" {
  environment_id = pingone_environment.my_environment.id
  gateway_id     = pingone_gateway.my_awesome_pingfederate_gateway.id
  role_id        = data.pingone_role.application_owner.id

  scope_application_id = pingone_application.my_awesome_application.id
}
