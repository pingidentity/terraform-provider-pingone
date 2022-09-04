resource "pingone_certificate_signing_response" "my_signed_response" {
  environment_id = pingone_environment.my_environment.id

  key_id               = pingone_key.my_key.id
  pem_ca_response_file = var.pem_ca_response_file
}