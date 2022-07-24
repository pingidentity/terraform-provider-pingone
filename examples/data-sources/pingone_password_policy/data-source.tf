data "pingone_password_policy" "example_by_name" {
  environment_id = var.environment_id

  name = "Standard"
}

data "pingone_password_policy" "example_by_id" {
  environment_id = var.environment_id

  password_policy_id = var.password_policy_id
}