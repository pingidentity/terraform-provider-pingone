data "pingone_licenses" "my_licenses_by_scim_filter" {
  organization_id = var.organization_id
  scim_filter     = "(status eq \"active\") and (beginsAt lt \"2009-11-10T23:00:00Z\")"
}

data "pingone_licenses" "my_licenses_by_data_filter" {
  organization_id = var.organization_id

  data_filter {
    name   = "name"
    values = ["My License"]
  }

  data_filter {
    name   = "status"
    values = ["ACTIVE"]
  }
}
