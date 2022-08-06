resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_worker_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Worker App"

  oidc_options {
    type                        = "WORKER"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
  }
}

data "pingone_role" "identity_data_admin" {
  name = "Identity Data Admin"
}

resource "pingone_application_role_assignment" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_worker_app
  role_id        = data.pingone_role.identity_data_admin.id

  scope_population_id = pingone_environment.my_environment.default_population_id
}
