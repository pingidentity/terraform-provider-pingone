resource "pingone_environment" "my_environment" {
  name        = "New Environment"
  description = "My new environment"
  type        = "SANDBOX"
  region      = "EU"
  license_id  = "ffc6b870-9709-4535-a78d-067f31add5e3"
  default_population {}
  service {}
}

resource "pingone_group" "my_group" {
  environment_id = pingone_environment.my_environment.id

  name        = "My group"
  description = "My new group"
}
