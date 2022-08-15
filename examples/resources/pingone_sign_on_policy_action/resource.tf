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

  login {
    recovery_enabled = true
  }

}

resource "pingone_sign_on_policy_action" "my_policy_mfa" {
  environment_id    = pingone_environment.my_environment.id
  sign_on_policy_id = pingone_sign_on_policy.my_policy.id

  mfa {}

}