resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_editor_policy" "my_awesome_editor_policy" {
  environment_id = pingone_environment.my_environment.id

  # ...
}
