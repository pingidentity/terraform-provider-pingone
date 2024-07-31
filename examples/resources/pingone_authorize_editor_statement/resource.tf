resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_editor_statement" "my_awesome_editor_statement" {
  environment_id = pingone_environment.my_environment.id

  # ...
}
