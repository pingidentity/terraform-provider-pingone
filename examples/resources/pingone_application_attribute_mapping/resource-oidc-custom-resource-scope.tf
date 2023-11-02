resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_resource" "my_resource" {
  environment_id = pingone_environment.my_environment.id

  name = "My resource"
}

resource "pingone_resource_scope" "my_resource_scope" {
  environment_id = pingone_environment.my_environment.id
  resource_id    = pingone_resource.my_resource.id

  name = "example_scope"
}

resource "pingone_application" "my_awesome_spa" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Single Page App"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://my-website.com"]
  }
}

resource "pingone_application_resource_grant" "custom_grant" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  resource_name = pingone_resource.my_resource.name

  scope_names = [
    pingone_resource_scope.my_resource_scope.name
  ]
}

resource "pingone_application_attribute_mapping" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  name  = "customAttribute"
  value = "$${user.email}"

  oidc_scopes = [
    pingone_resource_scope.my_resource_scope.id
  ]

  oidc_id_token_enabled = true
  oidc_userinfo_enabled = false

  depends_on = [
    pingone_application_resource_grant.custom_grant
  ]
}
