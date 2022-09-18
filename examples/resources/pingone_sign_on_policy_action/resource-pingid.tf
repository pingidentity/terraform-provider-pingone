resource "pingone_sign_on_policy_action" "my_policy_pingid" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 1

  pingid {}
}
