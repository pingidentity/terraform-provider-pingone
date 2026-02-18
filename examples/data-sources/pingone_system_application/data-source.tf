data "pingone_system_application" "example_application_portal_by_type" {
  environment_id = var.environment_id
  type           = "PING_ONE_PORTAL"
}

data "pingone_system_application" "example_self_service_by_type" {
  environment_id = var.environment_id
  type           = "PING_ONE_SELF_SERVICE"
}

data "pingone_system_application" "example_by_id" {
  environment_id = var.environment_id
  application_id = var.application_id
}