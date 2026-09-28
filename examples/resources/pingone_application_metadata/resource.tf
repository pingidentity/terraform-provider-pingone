resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_application" {
  # ...
}

resource "pingone_application_metadata" "example" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_application.id

  metadata = jsonencode({
    owner       = "platform-team"
    cost_center = "12345"
  })
}
