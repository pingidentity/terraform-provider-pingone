data "pingone_environment" "example_by_name" {
  name = "foo"
}

data "pingone_environment" "example_by_id" {
  environment_id = var.environment_id
}