resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_sign_on_policy" "my_policy" {
  environment_id = pingone_environment.my_environment.id

  name        = "foo"
  description = "My awesome Sign-on policy, username and password followed by MFA"
}

resource "pingone_sign_on_policy_action" "my_policy_first_factor" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 604800 // 7 days
  }

  login {
    recovery_enabled = true
  }
}

resource "pingone_sign_on_policy_action" "my_policy_mfa" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 2

  conditions {
    last_sign_on_older_than_seconds = 86400 // 24 hours

    ip_reputation_high_risk      = true
    geovelocity_anomaly_detected = true
    anonymous_network_detected   = true
  }

  mfa {
    device_sign_on_policy_id = var.my_device_sign_on_policy_id
  }
}
