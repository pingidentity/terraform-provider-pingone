data "pingone_certificate_export" "my_public_certificate" {
  environment_id = pingone_environment.my_environment.id

  key_id = pingone_key.my_key.id
}
