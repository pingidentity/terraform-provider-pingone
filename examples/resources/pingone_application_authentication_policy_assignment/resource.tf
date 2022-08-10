resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_spa" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Single Page App"

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://my-website.com"]
  }
}

resource "pingone_application_authentication_policy_assignment" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa

  authentication_policy_id = var.authentication_policy_id

  priority = 1
}
