resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_policy_management_root_policy" "my_awesome_root_policy" {
  environment_id = pingone_environment.my_environment.id

  # ...
}
