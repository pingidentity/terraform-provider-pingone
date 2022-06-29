resource "pingone_environment" "my_environment" {
  name        = "New Environment"
  description = "My new environment"
  type        = "SANDBOX"
  region      = "EU"
  license_id  = "ffc6b870-9709-4535-a78d-067f31add5e3"

  default_population {
    name        = "My Population"
    description = "My new population for users"
  }

  service {
    type = "SSO"
  }

  service {
    type        = "PING_FEDERATE"
    console_url = "https://my-pingfederate-console.example.com/pingfederate"
  }
}
