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

resource "pingone_application" "my_awesome_worker_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Worker App"
  enabled        = true

  oidc_options = {
    type                       = "WORKER"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
  }
}

data "pingone_role" "application_owner" {
  name = "Application Owner"
}

resource "pingone_application_role_assignment" "application_owner_to_application" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_worker_app.id
  role_id        = data.pingone_role.application_owner.id

  scope_application_id = pingone_application.my_awesome_application.id
}
