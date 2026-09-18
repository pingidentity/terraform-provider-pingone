data "pingone_credential_type_versions" "example_by_credential_type_id" {
  environment_id     = var.environment_id
  credential_type_id = var.credential_type_id
}
