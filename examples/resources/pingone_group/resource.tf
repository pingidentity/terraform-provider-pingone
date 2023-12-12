resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_group" "my_awesome_group" {
  environment_id = pingone_environment.my_environment.id

  name        = "My awesome group"
  description = "My new awesome group for people who are awesome"

  lifecycle {
    # change the `prevent_destroy` parameter value to `true` to prevent this data carrying resource from being destroyed
    prevent_destroy = false
  }
}