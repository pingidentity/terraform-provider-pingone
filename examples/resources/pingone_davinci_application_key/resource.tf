resource "pingone_davinci_application" "my_awesome_application" {
  environment_id = var.pingone_environment_id

  name = "My Awesome Application"

  oauth {
    grant_types                   = ["authorizationCode"]
    scopes                        = ["openid", "profile"]
    enforce_signed_request_openid = false
    redirect_uris                 = ["https://auth.pingone.com/0000-0000-000/rp/callback/openid_connect"]
  }
}

// Example of using the time provider to control regular rotation of application key
resource "time_rotating" "application_key_rotation_trigger" {
  rotation_days = 30
}

resource "pingone_davinci_application_key" "application_key_rotate" {
  environment_id = var.pingone_environment_id
  application_id = pingone_davinci_application.my_awesome_application.id

  rotation_trigger_values = {
    "rotation_rfc3339" : time_rotating.application_key_rotation_trigger.rotation_rfc3339,
  }
}