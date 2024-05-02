resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_password_policy" "my_password_policy" {
  environment_id = pingone_environment.my_environment.id

  name        = "My awesome password policy"
  description = "My new password policy"

  excludes_commonly_used_passwords = true
  excludes_profile_data            = true
  not_similar_to_current           = true

  history = {
    count          = 6
    retention_days = 365
  }

  length = {
    min = 8
    max = 255
  }

  password_age_max = 182
  password_age_min = 1

  lockout = {
    duration_seconds = 900
    failure_count    = 5
  }

  min_characters = {
    alphabetical_uppercase = 1
    alphabetical_lowercase = 1
    numeric                = 1
    special_characters     = 1
  }

  max_repeated_characters = 2
  min_complexity          = 7
  min_unique_characters   = 5
}
