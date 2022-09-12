resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_trusted_email_domain" "my_custom_email_domain" {
  environment_id = pingone_environment.my_environment.id

  domain_name = "demo.bxretail.org"
}
