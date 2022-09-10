resource "pingone_custom_domain_verify" "my_custom_domain" {
  environment_id = pingone_environment.my_environment.id

  custom_domain_id = pingone_custom_domain.my_custom_domain.id
}

resource "pingone_custom_domain_ssl" "my_custom_domain" {
  environment_id = pingone_environment.my_environment.id

  custom_domain_id = pingone_custom_domain.my_custom_domain.id

  certificate_pem_file               = var.my_domain_certificate_pem_file
  intermediate_certificates_pem_file = var.my_domain_intermediate_certificates_pem_file
  private_key_pem_file               = var.my_domain_private_key_pem_file

  depends_on = [
    pingone_custom_domain_verify.my_custom_domain
  ]
}
