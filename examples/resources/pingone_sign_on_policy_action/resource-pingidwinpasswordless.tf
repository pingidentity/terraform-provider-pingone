resource "pingone_sign_on_policy_action" "my_policy_pingid_windows_login_passwordless" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  priority = 1

  pingid_windows_login_passwordless {
    unique_user_attribute_name = "externalId"
    offline_mode_enabled       = true
  }
}
