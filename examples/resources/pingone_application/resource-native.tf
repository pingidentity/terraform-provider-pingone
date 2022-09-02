resource "pingone_application" "my_awesome_native_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Native App"
  enabled        = true

  oidc_options {
    type                        = "NATIVE_APP"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
  }
}
