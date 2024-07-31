resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_editor_condition" "my_awesome_editor_condition" {
  environment_id = pingone_environment.my_environment.id

  # ...
}
