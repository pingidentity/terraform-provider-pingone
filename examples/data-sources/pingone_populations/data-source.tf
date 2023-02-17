data "pingone_populations" "by_scim_filter" {
  environment_id = var.environment_id

  scim_filter = "(name eq \"My first population\") OR (name eq \"My second population\")"
}

data "pingone_populations" "by_data_filter" {
  environment_id = var.environment_id

  data_filter = {
    name   = "name"
    values = ["My first population", "My second population"]
  }
}
