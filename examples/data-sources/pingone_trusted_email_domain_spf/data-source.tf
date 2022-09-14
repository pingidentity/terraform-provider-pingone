data "pingone_trusted_email_domain_dkim" "email_domain_dkim" {
  environment_id = pingone_environment.my_environment.id

  trusted_email_domain_id = pingone_trusted_email_domain.id
}
