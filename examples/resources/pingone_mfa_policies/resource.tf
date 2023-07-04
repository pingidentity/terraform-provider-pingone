resource "pingone_mfa_policies" "mfa_policies" {
  environment_id = pingone_environment.my_environment.id

  migrate_data = [
    {
      device_authentication_policy_id = pingone_mfa_policy.my_mfa_policy.id
    },
    {
      device_authentication_policy_id = pingone_mfa_policy.my_mfa_policy_2.id
      fido2_policy_id                 = pingone_fido2_policy.my_fido2_policy.id
    }
  ]
}