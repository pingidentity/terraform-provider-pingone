data "pingone_population" "example_by_name" {
  environment_id = var.environment_id

  name = "foo"
}

data "pingone_population" "example_by_id" {
  environment_id = var.environment_id

  population_id = var.population_id
}