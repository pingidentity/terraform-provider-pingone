data "pingone_agreement_localization" "example_by_name" {
  environment_id = var.environment_id
  agreement_id   = var.agreement_id

  display_name = "foo bar"
}

data "pingone_agreement_localization" "example_by_locale" {
  environment_id = var.environment_id
  agreement_id   = var.agreement_id

  locale = "en"
}

data "pingone_agreement_localization" "example_by_id" {
  environment_id = var.environment_id
  agreement_id   = var.agreement_id

  agreement_localization_id = var.agreement_localization_id
}