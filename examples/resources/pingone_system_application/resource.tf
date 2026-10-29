resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_system_application" "pingone_portal" {
  environment_id = pingone_environment.my_environment.id

  type    = "PING_ONE_PORTAL"
  enabled = true

  logout_type = "SAML2_SLO"
}

resource "pingone_system_application" "pingone_self_service" {
  environment_id = pingone_environment.my_environment.id

  type    = "PING_ONE_SELF_SERVICE"
  enabled = true

  apply_default_theme         = true
  enable_default_theme_footer = true
  logout_type                 = "SAML2_SLO"

  manage_devices_via_my_account_enabled = true
}