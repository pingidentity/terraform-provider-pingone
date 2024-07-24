resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_editor_attribute" "my_awesome_editor_attribute" {
  environment_id = pingone_environment.my_environment.id

  # ...
}
