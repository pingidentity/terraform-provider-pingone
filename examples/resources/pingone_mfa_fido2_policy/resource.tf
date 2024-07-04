resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_mfa_fido2_policy" "my_awesome_fido2_policy" {
  environment_id = pingone_environment.my_environment.id
  name           = "My awesome FIDO2 policy"

  attestation_requirements = "DIRECT"
  authenticator_attachment = "BOTH"

  backup_eligibility = {
    allow                         = true
    enforce_during_authentication = false
  }

  device_display_name = "Test Device Max"

  discoverable_credentials = "PREFERRED"

  mds_authenticators_requirements = {
    allowed_authenticator_ids = [
      "authenticator_id_1",
      "authenticator_id_3",
      "authenticator_id_2",
    ]

    enforce_during_authentication = true
    option                        = "SPECIFIC"
  }

  relying_party_id = "pingidentity.com"

  user_display_name_attributes = {
    attributes = [
      {
        name = "email"
      },
      {
        name = "name",
        sub_attributes = [
          {
            name = "given"
          },
          {
            name = "family"
          }
        ]
      },
      {
        name = "username"
      }
    ]
  }

  user_verification = {
    enforce_during_authentication = true
    option                        = "REQUIRED"
  }
}
