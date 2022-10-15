data "pingone_trusted_email_domain_ownership" "email_domain_ownership" {
  environment_id = pingone_environment.my_environment.id

  trusted_email_domain_id = pingone_trusted_email_domain.my_custom_email_domain.id
}
