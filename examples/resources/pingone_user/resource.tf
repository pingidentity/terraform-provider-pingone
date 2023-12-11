resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_population" "my_population" {
  environment_id = pingone_environment.my_environment.id

  name = "My population of awesome identities"

  lifecycle {
    # change the `prevent_destroy` parameter value to `true` to prevent this data carrying resource from being destroyed
    prevent_destroy = false
  }
}

resource "pingone_user" "foo" {
  environment_id = pingone_environment.my_environment.id

  population_id = pingone_population.my_population.id

  username = "foouser"
  email    = "foouser@pingidentity.com"
}
