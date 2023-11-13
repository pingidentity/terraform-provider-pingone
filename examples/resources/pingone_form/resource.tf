resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_form" "my_awesome_form" {
  environment_id = pingone_environment.my_environment.id
}
