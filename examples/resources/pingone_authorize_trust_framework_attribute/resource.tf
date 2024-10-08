resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_trust_framework_attribute" "my_awesome_attribute" {
  environment_id = pingone_environment.my_environment.id

  # ...
}
