resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_rate_limit_configuration" "my_rate_limit_configuration" {
  environment_id = pingone_environment.my_environment.id

  type  = "WHITELIST"
  value = "192.0.2.0/24"
}
