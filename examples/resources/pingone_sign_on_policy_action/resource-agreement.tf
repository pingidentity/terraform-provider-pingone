resource "pingone_sign_on_policy_action" "my_policy_agreement" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 3

  agreement {
    agreement_id        = var.my_agreement_id
    show_decline_option = false
  }
}
