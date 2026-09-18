resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_application" {
  # ...
}

resource "time_rotating" "application_secret_rotation" {
  rotation_days = 30
}

resource "time_offset" "application_secret_expiry" {
  offset_days = 7
  triggers = {
    rotation_rfc3339 = time_rotating.application_secret_rotation.rotation_rfc3339
  }
}

resource "pingone_application_secret" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_application.id

  previous = {
    expires_at = time_offset.application_secret_expiry.rfc3339
  }

  regenerate_trigger_values = {
    "rotation_rfc3339" : time_rotating.application_secret_rotation.rotation_rfc3339,
  }

  lifecycle {
    # PingOne removes the previous secret once `expires_at` passes, which the provider then
    # reports as drift. Ignore it here so `terraform apply` doesn't fail; the secret is still
    # replaced whenever `regenerate_trigger_values` changes on rotation.
    ignore_changes = [previous]
  }
}
