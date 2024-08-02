resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_editor_rule" "my_awesome_editor_rule" {
  environment_id = pingone_environment.my_environment.id

  # ...
}
