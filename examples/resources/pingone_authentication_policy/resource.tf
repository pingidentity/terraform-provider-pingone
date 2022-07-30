resource "pingone_environment" "my_environment" {
  name        = "New Environment"
  description = "My new environment"
  type        = "SANDBOX"
  license_id  = var.license_id
  default_population {}
  service {}
}

resource "pingone_authentication_policy" "foo" {
  environment_id = pingone_environment.my_environment.id

  name        = "foo"
  description = "My awesome authentication policy, username and password followed by MFA"

  policy_action {
    action_type = "LOGIN"
  }

  policy_action {
    action_type = "MULTI_FACTOR_AUTHENTICATION"

    mfa_options {
      device_authentication_policy_id = var.device_authentication_policy_id
    }
  }
}

resource "pingone_authentication_policy" "bar" {
  environment_id = pingone_environment.my_environment.id

  name        = "bar"
  description = "My second awesome authentication policy, delegate to external Identity Provider"

  policy_action {
    action_type = "IDENTITY_PROVIDER"

    identity_provider_options {
      identity_provider_id = var.identity_provider_id
    }
  }
}