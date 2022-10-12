data "pingone_trusted_email_domain_spf" "email_domain_spf" {
  environment_id = pingone_environment.my_environment.id

  trusted_email_domain_id = pingone_trusted_email_domain.id
}
