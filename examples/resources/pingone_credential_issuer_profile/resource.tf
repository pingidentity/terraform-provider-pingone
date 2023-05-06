resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_credential_issuer_profile" "my_credential_issuer" {
  environment_id = pingone_environment.my_environment.id

  name = "Awesome PingOne Credentials Issuer"

}
