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

resource "pingone_authorize_application_role" "my_awesome_application_role" {
  environment_id = pingone_environment.my_environment.id
  name           = "CEO"
  description    = "The CEO of the company"
}

resource "pingone_user_application_role_assignment" "my_awesome_application_role_assignment" {
  environment_id      = pingone_environment.my_environment.id
  user_id             = pingone_user.foo.id
  application_role_id = pingone_authorize_application_role.my_awesome_application_role.id
}