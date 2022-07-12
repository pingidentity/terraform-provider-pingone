resource "pingone_environment" "my_environment" {
  name        = "New Environment"
  description = "My new environment"
  type        = "SANDBOX"
  region      = "EU"
  license_id  = var.license_id
  default_population {}
  service {}
}

resource "pingone_group" "my_group" {
  environment_id = pingone_environment.my_environment.id

  name        = "My group"
  description = "My new group"
}
