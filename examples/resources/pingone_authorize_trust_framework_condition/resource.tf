resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_trust_framework_condition" "my_awesome_condition" {
  environment_id = pingone_environment.my_environment.id

  # ...
}
