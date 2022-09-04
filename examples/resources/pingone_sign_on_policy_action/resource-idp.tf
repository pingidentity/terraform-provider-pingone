resource "pingone_sign_on_policy_action" "my_policy_identity_provider" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 604800 // 7 days
  }

  identity_provider {
    identity_provider_id = pingone_identity_provider.my_identity_provider.id

    acr_values        = "MFA"
    pass_user_context = true
  }
}
