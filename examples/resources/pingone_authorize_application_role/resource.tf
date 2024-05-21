resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_authorize_application_role" "my_awesome_application_role" {
  environment_id = pingone_environment.my_environment.id

  name        = "CEO"
  description = "The CEO"
}
