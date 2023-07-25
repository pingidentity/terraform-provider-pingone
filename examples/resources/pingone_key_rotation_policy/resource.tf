resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_key_rotation_policy" "my_awesome_key_rotation_policy" {
  environment_id = pingone_environment.my_environment.id

  name = "My Awesome Key Rotation Policy"

  algorithm           = "RSA"
  subject_dn          = "CN=awesomeness, OU=Ping Identity, O=Ping Identity, L=, ST=, C=US"
  key_length          = 4096
  signature_algorithm = "SHA256withRSA"
  usage_type          = "SIGNING"
}