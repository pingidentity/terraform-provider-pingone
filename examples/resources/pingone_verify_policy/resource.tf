resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_verify_voice_phrase" "my_verify_voice_phrase" {
  environment_id = pingone_environment.my_environment.id
  display_name   = "My Verify Voice Phrase for my Verify Policy"
}

resource "pingone_verify_voice_phrase_content" "my_verify_voice_phrase_content" {
  environment_id  = pingone_environment.my_environment.id
  voice_phrase_id = pingone_verify_voice_phrase.my_verify_voice_phrase.id
  locale          = "en"
  content         = "My voice content to be used in voice enrollment or verification."
}

resource "pingone_verify_policy" "my_verify_everything_policy" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Verify Policy"
  description    = "Example - All Verification Checks Required"

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
    verify            = "REQUIRED"
    create_mfa_device = true
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
    verify            = "REQUIRED"
    create_mfa_device = true
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

  voice = {
    verify               = "OPTIONAL"
    enrollment           = false
    comparison_threshold = "LOW"
    liveness_threshold   = "LOW"

    text_dependent = {
      samples         = "5"
      voice_phrase_id = pingone_verify_voice_phrase.my_verify_voice_phrase.id
    }

    reference_data = {
      retain_original_recordings = false
      update_on_reenrollment     = false
      update_on_verification     = false
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
