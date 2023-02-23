data "pingone_trusted_email_domain" "example_by_name" {
  environment_id = var.environment_id

  domain_name = "demo.bxretail.org"
}

data "pingone_trusted_email_domain" "example_by_id" {
  environment_id = var.environment_id

  trusted_email_domain_id = var.trusted_email_domain_id
}