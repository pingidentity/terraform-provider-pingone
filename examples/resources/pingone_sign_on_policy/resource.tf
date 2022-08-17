resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_sign_on_policy" "foo" {
  environment_id = pingone_environment.my_environment.id

  name        = "foo"
  description = "My awesome Sign-on policy, username and password followed by MFA"

}
