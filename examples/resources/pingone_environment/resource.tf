resource "pingone_environment" "my_environment" {
  name        = "New Environment"
  description = "My new environment"
  type        = "SANDBOX"
  license_id  = var.license_id

  service {
    type = "SSO"
  }

  service {
    type = "MFA"
  }

  service {
    type        = "PingFederate"
    console_url = "https://my-pingfederate-console.example.com/pingfederate"
  }
}
