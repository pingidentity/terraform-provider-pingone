resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_group" "my_group" {
  environment_id = pingone_environment.my_environment.id

  name = "My Awesome Group"

  lifecycle {
    # change the `prevent_destroy` parameter value to `true` to prevent this data carrying resource from being destroyed
    prevent_destroy = false
  }
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

data "pingone_role" "application_owner" {
  name = "Application Owner"
}

resource "pingone_group_role_assignment" "application_owner_to_group" {
  environment_id = pingone_environment.my_environment.id
  group_id       = pingone_group.my_group.id
  role_id        = data.pingone_role.application_owner.id

  scope_application_id = pingone_application.my_awesome_application.id
}
