resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_verify_policy" "my_verify_everything_policy" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Verify Policy"
  description    = "Example - All Verification Checks Required"
  default        = false

  government_id = {
    verify = "REQUIRED"
  }

  facial_comparison = {
    verify    = "REQUIRED"
    threshold = "HIGH"
  }

  liveness = {
    verify    = "REQUIRED"
    threshold = "HIGH"
  }

  email = {
    verify = "REQUIRED"
    create_mfa_device : true
    otp = {
      attempts = {
        count = "5"
      }
      lifetime = {
        duration  = "10"
        time_unit = "MINUTES"
      },
      deliveries = {
        count = 3
        cooldown = {
          duration  = "30"
          time_unit = "SECONDS"
        }
      }
      notification = {
        variant_name = "custom_variant_a"
      }
    }
  }

  phone = {
    verify = "REQUIRED"
    create_mfa_device : true
    otp = {
      attempts = {
        count = "5"
      }
      lifetime = {
        duration  = "10"
        time_unit = "MINUTES"
      },
      deliveries = {
        count = 3
        cooldown = {
          duration  = "30"
          time_unit = "SECONDS"
        }
      }
    }
  }

  transaction = {
    timeout = {
      duration  = "30"
      time_unit = "MINUTES"
    }

    data_collection = {
      timeout = {
        duration  = "15"
        time_unit = "MINUTES"
      }
    }

    data_collection_only = false
  }

}
