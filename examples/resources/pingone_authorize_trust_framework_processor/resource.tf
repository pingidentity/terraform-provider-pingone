resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_trust_framework_processor" "my_awesome_processor" {
  environment_id = pingone_environment.my_environment.id

  # ...
}
