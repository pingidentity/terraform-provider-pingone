resource "pingone_davinci_application" "my_awesome_registration_flow_application" {
  environment_id = var.pingone_environment_id

  name = "My Awesome Registration Application"

  oauth {
    grant_types                   = ["authorizationCode"]
    scopes                        = ["openid", "profile"]
    enforce_signed_request_openid = false
    redirect_uris                 = ["https://auth.pingone.com/0000-0000-000/rp/callback/openid_connect"]
  }
}

resource "pingone_davinci_application_flow_policy" "my_awesome_registration_flow_application_policy" {
  environment_id = var.pingone_environment_id
  application_id = pingone_davinci_application.my_awesome_registration_flow_application.id

  name   = "PingOne - Registration"
  status = "enabled"

  flow_distributions = [
    {
      id      = pingone_davinci_flow.registration.id
      version = -1
      weight  = 100
    }
  ]
}
