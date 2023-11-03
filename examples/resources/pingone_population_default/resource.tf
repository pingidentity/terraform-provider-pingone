resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_population_default" "my_default_population" {
  environment_id = pingone_environment.my_environment.id

  name        = "My default population"
  description = "A resource that overwrites the default population"
}
