resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_mfa_fido_policy" "my_awesome_fido_policy" {
  environment_id = pingone_environment.my_environment.id
  name           = "My awesome FIDO policy"

  attestation_requirements = "CERTIFIED"
  resident_key_requirement = "REQUIRED"

  enforce_during_authentication = true
}
