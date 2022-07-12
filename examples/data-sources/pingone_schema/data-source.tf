data "pingone_schema" "example_by_name" {
  environment_id = var.environment_id

  name = "User"
}

data "pingone_schema" "example_by_id" {
  environment_id = var.environment_id

  schema_id = var.schema_id
}