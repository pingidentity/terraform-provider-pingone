resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_sign_on_policy" "my_policy" {
  environment_id = pingone_environment.my_environment.id

  name        = "foo"
  description = "My awesome authentication policy, username and password followed by MFA"

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

    ip_reputation_high_risk = true
    geovelocity_anomaly_detected = true
    anonymous_network_detected = true

    user_attribute_equals {
      attribute_reference = "$${user.lifecycle.status}"
      value = "VERIFICATION_REQUIRED"
    }

  }

  mfa {
    device_sign_on_policy_id = var.my_device_sign_on_policy_id
  }

}

resource "pingone_sign_on_policy_action" "my_policy_agreement" {
  environment_id           = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 3

	agreement {
    agreement_id = var.my_agreement_id
    show_decline_option = false
  }

}

resource "pingone_sign_on_policy_action" "my_policy_progressive_profiling" {
  environment_id           = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 4

	progressive_profiling {

    attribute {
      name = "name.given"
      required = false
    }

    attribute {
      name = "name.family"
      required = true
    }

    prompt_text = "For the best experience, we need a couple things from you."

  }

}