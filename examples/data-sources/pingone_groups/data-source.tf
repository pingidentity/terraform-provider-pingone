data "pingone_groups" "by_scim_filter" {
  environment_id = var.environment_id

  scim_filter = "(name eq \"My first group\") OR (name eq \"My second group\")"
}

data "pingone_groups" "by_data_filter" {
  environment_id = var.environment_id

  data_filter {
    name   = "name"
    values = ["My first group", "My second group"]
  }
}
