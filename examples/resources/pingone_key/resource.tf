resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_key" "my_key" {
  environment_id = pingone_environment.my_environment.id

  name = "%[4]s"
  subject_dn = "CN=%[4]s, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"

		algorithm = "RSA"
		key_length = 3072
		signature_algorithm = "SHA512withRSA"

		usage_type = "ENCRYPTION"

    validity_period = 360
}