resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_gateway" "my_awesome_api_gateway" {
  environment_id = pingone_environment.my_environment.id
  name           = "Awesome API Gateway"
  enabled        = true

  api_gateway {}
}
