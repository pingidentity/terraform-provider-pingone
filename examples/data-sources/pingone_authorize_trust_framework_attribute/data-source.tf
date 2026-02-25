data "pingone_authorize_trust_framework_attribute" "example_by_full_name" {
  environment_id = var.environment_id
  full_name      = "PingOne.User"
}

data "pingone_authorize_trust_framework_attribute" "example_by_id" {
  environment_id = var.environment_id
  attribute_id   = var.attribute_id
}
