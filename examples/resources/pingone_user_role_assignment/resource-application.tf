resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_application" {
  environment_id = pingone_environment.my_environment.id
  name           = "My awesome application"
  enabled        = true

  oidc_options = {
    type                       = "WEB_APP"
    grant_types                = ["REFRESH_TOKEN", "AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
    redirect_uris              = ["https://www.pingidentity.com"]
  }
}

data "pingone_role" "application_owner" {
  name = "Application Owner"
}

resource "pingone_user_role_assignment" "application_owner_to_user" {
  environment_id = pingone_environment.my_environment.id
  user_id        = var.user_id
  role_id        = data.pingone_role.application_owner.id

  scope_application_id = pingone_application.my_awesome_application.id
}
