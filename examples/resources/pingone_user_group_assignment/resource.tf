resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_population" "my_population" {
  environment_id = pingone_environment.my_environment.id

  name = "My population of awesome identities"
}

resource "pingone_user" "foo" {
  environment_id = pingone_environment.my_environment.id

  population_id = pingone_population.my_population.id

  username = "foouser"
  email    = "foouser@pingidentity.com"
}

resource "pingone_group" "my_awesome_group" {
  environment_id = pingone_environment.my_environment.id

  name        = "My awesome group"
  description = "My new awesome group for people who are awesome"
}

resource "pingone_user_group_assignment" "bar" {
  environment_id = pingone_environment.my_environment.id

  user_id  = pingone_user.foo.id
  group_id = pingone_group.my_awesome_group.id
}
