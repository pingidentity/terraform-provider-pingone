resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_editor_service" "my_awesome_editor_service" {
  environment_id = pingone_environment.my_environment.id

  # ...
}
