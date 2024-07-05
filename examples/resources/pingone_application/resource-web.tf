resource "pingone_application" "my_awesome_web_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Web App"
  enabled        = true

  oidc_options = {
    type                       = "WEB_APP"
    grant_types                = ["AUTHORIZATION_CODE", "REFRESH_TOKEN"]
    response_types             = ["CODE"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
    redirect_uris              = ["https://my-website.com"]
  }
}

resource "time_rotating" "my_awesome_web_app_secret_rotation" {
  rotation_days = 30
}

resource "pingone_application_secret" "my_awesome_web_app" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_web_app.id

  regenerate_trigger_values = {
    "rotation_rfc3339" : time_rotating.my_awesome_web_app_secret_rotation.rotation_rfc3339,
  }
}
