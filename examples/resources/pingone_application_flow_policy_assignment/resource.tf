resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_application" {
  # ...
}

resource "davinci_application" "davinci_app" {
  environment_id = pingone_environment.my_environment.id

  name = "Awesome DaVinci Application"

  oauth {
    enabled = true
    values {
      allowed_grants                = ["authorizationCode"]
      allowed_scopes                = ["openid", "profile"]
      enabled                       = true
      enforce_signed_request_openid = false
      redirect_uris                 = [var.redirect_uri]
    }
  }
  policy {
    name   = "Flow Policy"
    status = "enabled"
    policy_flow {
      flow_id    = var.davinci_flow_id
      version_id = -1
      weight     = 100
    }
  }

  saml {
    values {
      enabled                = false
      enforce_signed_request = false
    }
  }
}

resource "pingone_application_flow_policy_assignment" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_application.id

  flow_policy_id = davinci_application.davinci_app.policy.*.policy_id[0]

  priority = 1
}
