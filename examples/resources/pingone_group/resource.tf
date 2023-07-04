resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_group" "my_awesome_group" {
  environment_id = pingone_environment.my_environment.id

  name        = "My awesome group"
  description = "My new awesome group for people who are awesome"
}