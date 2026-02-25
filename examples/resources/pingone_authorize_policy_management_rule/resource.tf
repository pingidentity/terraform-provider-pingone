resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_policy_management_rule" "my_awesome_policy_rule" {
  environment_id = pingone_environment.my_environment.id

  # ...
}
