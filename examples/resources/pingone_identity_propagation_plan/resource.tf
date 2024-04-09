resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_identity_propagation_plan" "my_awesome_propagation_plan" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Identity Provisioning Plan"
}
