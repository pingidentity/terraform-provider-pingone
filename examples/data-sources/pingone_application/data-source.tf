data "pingone_application" "example_by_name" {
  environment_id = var.environment_id
  name           = "foo"
}

data "pingone_application" "example_by_id" {
  environment_id = var.environment_id
  application_id = var.application_id
}