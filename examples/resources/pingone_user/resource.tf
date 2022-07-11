resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_population" "my_population" {
  name = "My population of awesome identities"
}

resource "pingone_user" "foo" {
  environment_id = pingone_environment.my_environment.id

  population_id = pingone_population.my_population.id

  username = "foouser"
  email    = "foouser@pingidentity.com"
}
