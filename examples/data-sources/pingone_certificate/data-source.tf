data "pingone_certificate" "example_by_name" {
  environment_id = var.environment_id

  name = "My Certificate"
}

data "pingone_certificate" "example_by_id" {
  environment_id = var.environment_id

  certificate_id = var.certificate_id
}
