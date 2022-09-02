resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_group" "my_group" {
  environment_id = pingone_environment.my_environment.id

  name        = "My group"
  description = "My new group"
}
