data "pingone_environment" "my_environment" {
  name = "My Environment"
}

data "pingone_key_rotation_policy" "example_by_name" {
  environment_id = data.pingone_environment.my_environment.id

  name = "My Awesome Key Rotation Policy"
}

data "pingone_key_rotation_policy" "example_by_id" {
  environment_id = data.pingone_environment.my_environment.id

  key_rotation_policy_id = var.key_rotation_policy_id
}
