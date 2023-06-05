resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_system_application" "pingone_portal" {
  environment_id = pingone_environment.my_environment.id

  type    = "PING_ONE_PORTAL"
  enabled = true
}

resource "pingone_system_application" "pingone_self_service" {
  environment_id = pingone_environment.my_environment.id

  type    = "PING_ONE_SELF_SERVICE"
  enabled = true
}