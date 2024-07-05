resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_spa" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Single Page App"
  enabled        = true

  oidc_options = {
    type                       = "SINGLE_PAGE_APP"
    grant_types                = ["AUTHORIZATION_CODE"]
    response_types             = ["CODE"]
    pkce_enforcement           = "S256_REQUIRED"
    token_endpoint_auth_method = "NONE"
    redirect_uris              = ["https://my-website.com"]
  }
}

resource "pingone_application_resource_grant" "oidc_grant" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  resource_name = "openid"

  scope_names = [
    "profile"
  ]
}

resource "pingone_application_attribute_mapping" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  name  = "customAttribute"
  value = "$${user.email}"

  oidc_scopes = [
    data.pingone_resource_scope.openid_profile.id
  ]

  oidc_id_token_enabled = true
  oidc_userinfo_enabled = false

  depends_on = [
    pingone_application_resource_grant.oidc_grant
  ]
}
