resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_trusted_email_domain" "my_custom_email_domain" {
  environment_id = pingone_environment.my_environment.id

  domain_name = "demo.bxretail.org"
}

resource "pingone_trusted_email_address" "my_trusted_email" {
  environment_id  = pingone_environment.my_environment.id
  email_domain_id = pingone_trusted_email_domain.my_custom_email_domain.id

  email_address = "noreply@demo.bxretail.org"
}
