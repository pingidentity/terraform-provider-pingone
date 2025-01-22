data "pingone_custom_role" "example_by_name" {
  environment_id = var.environment_id

  name = "foo"
}

data "pingone_custom_role" "example_by_id" {
  environment_id = var.environment_id

  role_id = var.role_id
}