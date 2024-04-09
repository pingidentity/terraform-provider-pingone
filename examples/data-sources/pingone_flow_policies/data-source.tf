data "pingone_flow_policies" "by_scim_filter" {
  environment_id = var.environment_id

  scim_filter = "(trigger.type eq \"AUTHENTICATION\")"
}

data "pingone_flow_policies" "by_data_filter" {
  environment_id = var.environment_id

  data_filters = [
    {
      name   = "trigger.type"
      values = ["AUTHENTICATION"]
    }
  ]
}
