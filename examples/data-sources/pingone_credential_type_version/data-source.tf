data "pingone_credential_type_version" "example_by_version_id" {
  environment_id             = var.environment_id
  credential_type_id         = var.credential_type_id
  credential_type_version_id = var.credential_type_version_id
}
