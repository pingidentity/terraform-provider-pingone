data "pingone_user" "example_by_username" {
  environment_id = var.environment_id

  username = "user123"
}

data "pingone_user" "example_by_email" {
  environment_id = var.environment_id

  email = "user123@bxretail.org"
}

data "pingone_user" "example_by_id" {
  environment_id = var.environment_id

  user_id = var.user_id
}