data "pingone_verify_policy" "find_by_id_example" {
  environment_id   = var.environment_id
  verify_policy_id = var.verify_policy_id
}

data "pingone_verify_policy" "find_by_name_example" {
  environment_id = var.environment_id
  name           = "foo"
}

data "pingone_verify_policy" "find_default_policy_example" {
  environment_id = var.environment_id
  default        = true
}