data "pingone_mfa_policies" "example_all_mfa_policy_ids" {
  environment_id = pingone_environment.my_environment.id
}

resource "pingone_mfa_policies" "mfa_policies" {
  environment_id = pingone_environment.my_environment.id

  migrate_data = [
    for policy_id in data.pingone_mfa_policies.example_all_mfa_policy_ids.ids : {
      device_authentication_policy_id = policy_id
    }
  ]
}