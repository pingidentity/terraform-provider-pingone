resource "pingone_sign_on_policy_action" "my_policy_identifier_first" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 1

  conditions {
    last_sign_on_older_than_seconds = 604800 // 7 days
  }

  identifier_first {
    recovery_enabled = true

    discovery_rule {
      attribute_contains_text = "@pingidentity.com"
      identity_provider_id    = pingone_identity_provider.my_identity_provider.id
    }
  }
}
