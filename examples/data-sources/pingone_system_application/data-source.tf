data "pingone_system_application" "example_application_portal_by_name" {
  environment_id = var.environment_id
  name           = "PingOne Application Portal"
}

data "pingone_system_application" "example_self_service_by_name" {
  environment_id = var.environment_id
  name           = "PingOne Self-Service - MyAccount"
}

data "pingone_system_application" "example_by_id" {
  environment_id = var.environment_id
  application_id = var.application_id
}