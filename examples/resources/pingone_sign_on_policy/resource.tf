resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_sign_on_policy" "foo" {
  environment_id = pingone_environment.my_environment.id

  name        = "foo"
  description = "My awesome authentication policy, username and password followed by MFA"

}

resource "pingone_sign_on_policy" "bar" {
  environment_id = pingone_environment.my_environment.id

  name        = "bar"
  description = "My second awesome authentication policy, delegate to external Identity Provider"

}