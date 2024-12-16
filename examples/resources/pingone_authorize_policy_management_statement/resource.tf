resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_policy_management_statement" "my_awesome_policy_statement" {
  environment_id = pingone_environment.my_environment.id

  # ...
}
