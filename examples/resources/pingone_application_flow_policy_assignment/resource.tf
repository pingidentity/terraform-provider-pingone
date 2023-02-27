resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_application" {
  # ...
}

resource "pingone_application_flow_policy_assignment" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_application.id

  flow_policy_id = var.davinci_flow_policy_id
}
