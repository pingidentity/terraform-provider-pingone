data "pingone_credential_issuance_rule" "example_by_id" {
  environment_id              = var.environment_id
  credential_type_id          = var.credential_type_id
  credential_issuance_rule_id = var.credential_issuance_rule_id
}
