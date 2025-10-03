resource "pingone_davinci_application" "my_awesome_application" {
  environment_id = var.pingone_environment_id

  name = "My Awesome Application"

  oauth = {
    grant_types                   = ["authorizationCode"]
    scopes                        = ["openid", "profile"]
    enforce_signed_request_openid = false
    redirect_uris                 = ["https://auth.pingone.com/0000-0000-000/rp/callback/openid_connect"]
  }
}

resource "pingone_davinci_application_flow_policy" "authentication_flow_policy" {
  environment_id = var.pingone_environment_id
  application_id = pingone_davinci_application.my_awesome_application.id

  name   = "PingOne - Authentication"
  status = "enabled"

  flow_distributions = [
    {
      flow_id = pingone_davinci_flow.authentication.id
      version = -1
      weight  = 100
    }
  ]
}

resource "pingone_davinci_application_flow_policy" "registration_flow_policy" {
  environment_id = var.pingone_environment_id
  application_id = pingone_davinci_application.my_awesome_application.id

  name   = "PingOne - Registration"
  status = "enabled"

  flow_distributions = [
    {
      flow_id = pingone_davinci_flow.registration.id
      version = -1
      weight  = 100
    }
  ]
}
