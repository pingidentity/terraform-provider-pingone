resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_mfa_policy" "my_awesome_mfa_policy" {
  environment_id = pingone_environment.my_environment.id
  name           = "My awesome MFA policy"

  mobile {
    enabled = false
  }

  totp {
    enabled = true
  }

  security_key {
    enabled = true
  }

  platform {
    enabled = true
  }

  sms {
    enabled = false
  }

  voice {
    enabled = false
  }

  email {
    enabled = false
  }

}
