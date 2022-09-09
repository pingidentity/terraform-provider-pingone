resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_custom_domain" "my_custom_domain" {
  environment_id = pingone_environment.my_environment.id

  domain_name = "auth.bxretail.org"
}
