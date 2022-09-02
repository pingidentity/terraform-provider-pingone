resource "pingone_key" "my_encryption_key" {
  environment_id = pingone_environment.my_environment.id

  name       = "mycert"
  subject_dn = "CN=mycert, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"

  algorithm           = "RSA"
  key_length          = 3072
  signature_algorithm = "SHA512withRSA"

  usage_type = "ENCRYPTION"

  validity_period = 360
}
