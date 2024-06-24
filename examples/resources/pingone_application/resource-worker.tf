resource "pingone_application" "my_awesome_worker_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Worker App"
  enabled        = true

  oidc_options = {
    type                        = "WORKER"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
  }
}

resource "time_rotating" "my_awesome_worker_app_secret_rotation" {
  rotation_days = 30
}

resource "pingone_application_secret" "my_awesome_worker_app" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_worker_app.id

  regenerate_trigger_values = {
    "rotation_rfc3339" : time_rotating.my_awesome_worker_app_secret_rotation.rotation_rfc3339,
  }
}
