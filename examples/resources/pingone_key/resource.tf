resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_key" "my_encryption_key" {
  environment_id = pingone_environment.my_environment.id

  name       = "%[4]s"
  subject_dn = "CN=%[4]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"

  algorithm           = "RSA"
  key_length          = 3072
  signature_algorithm = "SHA512withRSA"

  usage_type = "ENCRYPTION"

  validity_period = 360
}

resource "pingone_key" "my_tls_key" {
  environment_id = pingone_environment.my_environment.id

  pkcs12_file_base64 = filebase64("../path/to/keyStore.p12")

  usage_type = "SSL/TLS"
}