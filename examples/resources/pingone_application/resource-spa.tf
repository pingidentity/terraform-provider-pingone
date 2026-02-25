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

    include_x5t                                   = true
    op_session_check_enabled                      = true
    request_scopes_for_multiple_resources_enabled = true
  }
}

resource "time_rotating" "my_awesome_spa_secret_rotation" {
  rotation_days = 30
}

resource "pingone_application_secret" "my_awesome_spa" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  regenerate_trigger_values = {
    "rotation_rfc3339" : time_rotating.my_awesome_spa_secret_rotation.rotation_rfc3339,
  }
}
