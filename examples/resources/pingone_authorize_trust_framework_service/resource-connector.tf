resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_trust_framework_service" "my_awesome_service" {
  environment_id = pingone_environment.my_environment.id

  # ...
}
