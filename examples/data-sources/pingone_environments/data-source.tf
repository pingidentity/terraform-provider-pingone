data "pingone_environments" "by_scim_filter" {
  scim_filter = "(name sw \"TEST-\") and (license.id eq \"${var.license_id}\")"
}
