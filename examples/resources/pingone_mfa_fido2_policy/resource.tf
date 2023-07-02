resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_mfa_fido2_policy" "my_awesome_fido2_policy" {
  environment_id = pingone_environment.my_environment.id
  name           = "My awesome FIDO2 policy"
}
