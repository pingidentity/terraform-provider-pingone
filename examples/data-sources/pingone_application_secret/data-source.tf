data "pingone_application_secret" "my_awesome_oidc_application" {
  environment_id = var.environment_id
  application_id = var.application_id
}
