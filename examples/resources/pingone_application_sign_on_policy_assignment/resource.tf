resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_application" {
  # ...
}

resource "pingone_sign_on_policy" "my_sign_on_policy" {
  environment_id = pingone_environment.my_environment.id

  name = "foo"
}

resource "pingone_application_sign_on_policy_assignment" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_application.id

  sign_on_policy_id = pingone_sign_on_policy.my_sign_on_policy.id

  priority = 1
}
