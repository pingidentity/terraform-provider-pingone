resource "pingone_key" "my_tls_key" {
  environment_id = pingone_environment.my_environment.id

  pkcs12_file_base64 = filebase64("../path/to/keyStore.p12")

  usage_type = "SSL/TLS"
}