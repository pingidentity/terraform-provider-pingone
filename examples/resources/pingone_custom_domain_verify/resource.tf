resource "pingone_custom_domain_verify" "my_custom_domain" {
  environment_id = pingone_environment.my_environment.id

  custom_domain_id = pingone_custom_domain.my_custom_domain.id
}
