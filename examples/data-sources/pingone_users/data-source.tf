data "pingone_users" "example_by_data_filter" {
  environment_id = var.environment_id

  data_filter {
    name = "memberOfGroups.id"
    values = [
      pingone_group.my_first_group.id,
      pingone_group.my_second_group.id
    ]
  }

  data_filter {
    name = "population.id"
    values = [
      pingone_population.my_population.id
    ]
  }
}

data "pingone_users" "example_by_scim_filter" {
  environment_id = var.environment_id

  scim_filter = "(population.id eq \"${pingone_population.my_population.id}\") AND (memberOfGroups[id eq \"${pingone_group.my_first_group.id}\"] OR memberOfGroups[id eq \"${pingone_group.my_second_group.id}\"])"
}