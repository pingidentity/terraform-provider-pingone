resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_application" {
  # ...
}

resource "pingone_application_attribute_mapping" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_application.id

  name  = "email"
  value = "$${user.email}"
}

resource "pingone_application_attribute_mapping" "bar" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_application.id

  name  = "full_name"
  value = "$${user.name.given + ', ' + user.name.family}"
}