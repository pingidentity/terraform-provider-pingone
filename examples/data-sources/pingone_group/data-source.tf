data "pingone_group" "example_by_name" {
  environment_id = var.environment_id

  name = "foo"
}

data "pingone_group" "example_by_id" {
  environment_id = var.environment_id

  group_id = var.group_id
}