resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_password_policy" "my_password_policy" {
  environment_id = pingone_environment.my_environment.id

  name        = "My awesome password policy"
  description = "My new password policy"

  exclude_commonly_used_passwords = true
  exclude_profile_data            = true
  not_similar_to_current          = true

  password_history {
    prior_password_count = 6
    retention_days       = 365
  }

  password_length {
    min = 8
    max = 255
  }

  password_age {
    max = 182
    min = 1
  }

  account_lockout {
    duration_seconds = 900
    fail_count       = 5
  }

  min_characters {
    alphabetical_uppercase = 1
    alphabetical_lowercase = 1
    numeric                = 1
    special_characters     = 1
  }

  max_repeated_characters = 2
  min_complexity          = 7
  min_unique_characters   = 5
}
