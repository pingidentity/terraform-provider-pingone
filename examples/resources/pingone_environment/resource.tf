resource "pingone_environment" "my_environment" {
  name        = "New Environment"
  description = "My new environment"
  type        = "SANDBOX"
  region      = "EU"
  license_id  = "ffc6b870-9709-4535-a78d-067f31add5e3"

  default_population_name        = "My Population"
  default_population_description = "My new population for users"
}
