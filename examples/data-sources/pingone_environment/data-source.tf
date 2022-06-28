data "pingone_environment" "example_by_name" {
  name = "foo"
}

data "pingone_environment" "example_by_id" {
  environment_id = "9b1dc0b5-a725-4436-b954-a3a52fe463de"
}