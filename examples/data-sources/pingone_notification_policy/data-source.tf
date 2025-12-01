data "pingone_notification_policy" "example_by_name" {
  environment_id = var.environment_id

  name = "My notification policy"
}

data "pingone_notification_policy" "example_by_id" {
  environment_id = var.environment_id

  notification_policy_id = var.notification_policy_id
}

data "pingone_notification_policy" "example_default" {
  environment_id = var.environment_id

  default = true
}
