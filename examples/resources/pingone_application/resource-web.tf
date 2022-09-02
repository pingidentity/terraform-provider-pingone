resource "pingone_application" "my_awesome_web_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Web App"
  enabled        = true

  oidc_options {
    type                        = "WEB_APP"
    grant_types                 = ["AUTHORIZATION_CODE", "REFRESH_TOKEN"]
    response_types              = ["CODE"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
    redirect_uris               = ["https://my-website.com"]
  }
}
