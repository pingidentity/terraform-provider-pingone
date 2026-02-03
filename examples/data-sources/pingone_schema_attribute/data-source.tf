data "pingone_schema_attribute" "example_by_name" {
  environment_id = var.environment_id
  schema_id      = var.schema_id

  name = "email"
}

data "pingone_schema_attribute" "example_by_id" {
  environment_id = var.environment_id
  schema_id      = var.schema_id

  attribute_id = var.attribute_id
}
