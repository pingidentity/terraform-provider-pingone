data "pingone_certificate_signing_request" "my_csr" {
  environment_id = pingone_environment.my_environment.id

  key_id = pingone_key.my_key.id
}
