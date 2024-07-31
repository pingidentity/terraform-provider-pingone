resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_editor_processor" "my_awesome_editor_processor" {
  environment_id = pingone_environment.my_environment.id

  # ...
}
