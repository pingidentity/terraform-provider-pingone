data "pingone_agreement" "example_by_name" {
  environment_id = var.environment_id

  name = "foo"
}

data "pingone_agreement" "example_by_id" {
  environment_id = var.environment_id

  agreement_id = var.agreement_id
}