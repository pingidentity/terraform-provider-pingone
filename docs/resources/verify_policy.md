---
page_title: "pingone_verify_policy Resource - terraform-provider-pingone"
subcategory: "Neo (Verify & Credentials)"
description: |-
  Resource to configure the requirements to verify a user, including the parameters for verification, such as the number of one-time password (OTP) attempts and OTP expiration.
  A verify policy defines which of the following five checks are performed for a verification transaction and configures the parameters of each check. The checks can be either required or optional. If a type is optional, then the transaction can be processed with or without the documents for that type. If the documents are provided for that type and the optional type verification fails, it will not cause the entire transaction to fail.
  Verify policies can perform any of five checks:
  - Government identity document - Validate a government-issued identity document, which includes a photograph.
  - Facial comparison - Compare a mobile phone self-image to a reference photograph, such as on a government ID or previously verified photograph.
  - Liveness - Inspect a mobile phone self-image for evidence that the subject is alive and not a representation, such as a photograph or mask.
  - Email - Receive a one-time password (OTP) on an email address and return the OTP to the service.
  - Phone - Receive a one-time password (OTP) on a mobile phone and return the OTP to the service.
  - Voice - Compare a voice recording to a previously submitted reference voice recording.
---

# pingone_verify_policy (Resource)

Resource to configure the requirements to verify a user, including the parameters for verification, such as the number of one-time password (OTP) attempts and OTP expiration.

A verify policy defines which of the following five checks are performed for a verification transaction and configures the parameters of each check. The checks can be either required or optional. If a type is optional, then the transaction can be processed with or without the documents for that type. If the documents are provided for that type and the optional type verification fails, it will not cause the entire transaction to fail.

Verify policies can perform any of five checks:
- Government identity document - Validate a government-issued identity document, which includes a photograph.
- Facial comparison - Compare a mobile phone self-image to a reference photograph, such as on a government ID or previously verified photograph.
- Liveness - Inspect a mobile phone self-image for evidence that the subject is alive and not a representation, such as a photograph or mask.
- Email - Receive a one-time password (OTP) on an email address and return the OTP to the service.
- Phone - Receive a one-time password (OTP) on a mobile phone and return the OTP to the service.
- Voice - Compare a voice recording to a previously submitted reference voice recording.

## Example Usage

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_verify_voice_phrase" "my_verify_voice_phrase" {
  environment_id = pingone_environment.my_environment.id
  display_name   = "My Verify Voice Phrase for my Verify Policy"
}

resource "pingone_verify_voice_phrase_content" "my_verify_voice_phrase_content" {
  environment_id  = pingone_environment.my_environment.id
  voice_phrase_id = pingone_voice_phrase.my_verify_voice_phrase.id
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
      voice_phrase_id = pingone_voice_phrase.my_verify_voice_phrase.id
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
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) PingOne environment identifier (UUID) in which the verify policy exists.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `name` (String) Name of the verification policy displayed in PingOne Admin UI.

### Optional

- `description` (String) Description of the verification policy displayed in PingOne Admin UI, 1-1024 characters.
- `email` (Attributes) Defines the verification requirements to validate an email address using a one-time password (OTP). (see [below for nested schema](#nestedatt--email))
- `facial_comparison` (Attributes) Defines the verification requirements to compare a mobile phone self-image to a reference photograph, such as on a government ID or previously verified photograph. (see [below for nested schema](#nestedatt--facial_comparison))
- `government_id` (Attributes) Defines the verification requirements for a government-issued identity document, which includes a photograph. (see [below for nested schema](#nestedatt--government_id))
- `liveness` (Attributes) Defines the verification requirements to inspect a mobile phone self-image for evidence that the subject is alive and not a representation, such as a photograph or mask. (see [below for nested schema](#nestedatt--liveness))
- `phone` (Attributes) Defines the verification requirements to validate a mobile phone number using a one-time password (OTP). (see [below for nested schema](#nestedatt--phone))
- `transaction` (Attributes) Defines the requirements for transactions invoked by the policy. (see [below for nested schema](#nestedatt--transaction))
- `voice` (Attributes) Defines the requirements for transactions invoked by the policy. (see [below for nested schema](#nestedatt--voice))

### Read-Only

- `created_at` (String) Date and time the verify policy was created.
- `default` (Boolean) Specifies whether this is the environment's default verify policy.
- `id` (String) The ID of this resource.
- `updated_at` (String) Date and time the verify policy was updated. Can be null.

<a id="nestedatt--email"></a>
### Nested Schema for `email`

Required:

- `verify` (String) Controls the verification requirements for an Email or Phone verification.  Options are `DISABLED`, `OPTIONAL`, `REQUIRED`.  Defaults to `DISABLED`.

Optional:

- `create_mfa_device` (Boolean) When enabled, PingOne Verify registers the email address with PingOne MFA as a verified MFA device.
- `otp` (Attributes) SMS/Voice/Email one-time password (OTP) configuration. (see [below for nested schema](#nestedatt--email--otp))

<a id="nestedatt--email--otp"></a>
### Nested Schema for `email.otp`

Required:

- `attempts` (Attributes) OTP attempts configuration. (see [below for nested schema](#nestedatt--email--otp--attempts))
- `deliveries` (Attributes) OTP delivery configuration. (see [below for nested schema](#nestedatt--email--otp--deliveries))
- `lifetime` (Attributes) The length of time for which the OTP is valid. (see [below for nested schema](#nestedatt--email--otp--lifetime))

Optional:

- `notification` (Attributes) OTP notification template configuration. (see [below for nested schema](#nestedatt--email--otp--notification))

<a id="nestedatt--email--otp--attempts"></a>
### Nested Schema for `email.otp.attempts`

Required:

- `count` (Number) Allowed maximum number of OTP failures.


<a id="nestedatt--email--otp--deliveries"></a>
### Nested Schema for `email.otp.deliveries`

Required:

- `cooldown` (Attributes) Cooldown (waiting period between OTP attempts) configuration. (see [below for nested schema](#nestedatt--email--otp--deliveries--cooldown))
- `count` (Number) Allowed maximum number of OTP deliveries.

<a id="nestedatt--email--otp--deliveries--cooldown"></a>
### Nested Schema for `email.otp.deliveries.count`

Required:

- `duration` (Number) Cooldown duration.
    - If `cooldown.time_unit` is `MINUTES`, the allowed range is `0 - 30`.
    - If `cooldown.time_unit` is `SECONDS`, the allowed range is `0 - 1800`.
    - Defaults to `30 SECONDS`.
- `time_unit` (String) Time unit of the cooldown duration configuration.  Options are `MINUTES`, `SECONDS`.  Defaults to `SECONDS`.



<a id="nestedatt--email--otp--lifetime"></a>
### Nested Schema for `email.otp.lifetime`

Required:

- `duration` (Number) Lifetime of the OTP delivered via email.
    - If `lifetime.time_unit` is `MINUTES`, the allowed range is `1 - 30`.
    - If `lifetime.time_unit` is `SECONDS`, the allowed range is `60 - 1800`.
    - Defaults to `10 MINUTES`.
- `time_unit` (String) Time unit of the OTP (Email) duration lifetime.  Options are `MINUTES`, `SECONDS`.  Defaults to `MINUTES`.


<a id="nestedatt--email--otp--notification"></a>
### Nested Schema for `email.otp.notification`

Optional:

- `variant_name` (String) Name of the template variant to use to pass a one-time passcode (OTP).

Read-Only:

- `template_name` (String) Name of the template to use to pass a one-time passcode (OTP). The default value of `email_phone_verification` is static. Use the `notification.variant_name` property to define an alternate template.




<a id="nestedatt--facial_comparison"></a>
### Nested Schema for `facial_comparison`

Required:

- `threshold` (String) Facial Comparison threshold requirements.  Options are `HIGH`, `LOW`, `MEDIUM`.  Defaults to `MEDIUM`.
- `verify` (String) Controls Facial Comparison verification requirements.  Options are `DISABLED`, `OPTIONAL`, `REQUIRED`.  Defaults to `DISABLED`.


<a id="nestedatt--government_id"></a>
### Nested Schema for `government_id`

Required:

- `verify` (String) Controls Government ID verification requirements.  Options are `DISABLED`, `OPTIONAL`, `REQUIRED`.  Defaults to `DISABLED`.


<a id="nestedatt--liveness"></a>
### Nested Schema for `liveness`

Required:

- `threshold` (String) Liveness Check threshold requirements.  Options are `HIGH`, `LOW`, `MEDIUM`.  Defaults to `MEDIUM`.
- `verify` (String) Controls Liveness Check verification requirements.  Options are `DISABLED`, `OPTIONAL`, `REQUIRED`.  Defaults to `DISABLED`.


<a id="nestedatt--phone"></a>
### Nested Schema for `phone`

Required:

- `verify` (String) Controls the verification requirements for an Email or Phone verification.  Options are `DISABLED`, `OPTIONAL`, `REQUIRED`.  Defaults to `DISABLED`.

Optional:

- `create_mfa_device` (Boolean) When enabled, PingOne Verify registers the mobile phone with PingOne MFA as a verified MFA device.
- `otp` (Attributes) SMS/Voice/Email one-time password (OTP) configuration. (see [below for nested schema](#nestedatt--phone--otp))

<a id="nestedatt--phone--otp"></a>
### Nested Schema for `phone.otp`

Required:

- `attempts` (Attributes) OTP attempts configuration. (see [below for nested schema](#nestedatt--phone--otp--attempts))
- `deliveries` (Attributes) OTP delivery configuration. (see [below for nested schema](#nestedatt--phone--otp--deliveries))
- `lifetime` (Attributes) The length of time for which the OTP is valid. (see [below for nested schema](#nestedatt--phone--otp--lifetime))

Optional:

- `notification` (Attributes) OTP notification template configuration. (see [below for nested schema](#nestedatt--phone--otp--notification))

<a id="nestedatt--phone--otp--attempts"></a>
### Nested Schema for `phone.otp.attempts`

Required:

- `count` (Number) Allowed maximum number of OTP failures.


<a id="nestedatt--phone--otp--deliveries"></a>
### Nested Schema for `phone.otp.deliveries`

Required:

- `cooldown` (Attributes) Cooldown (waiting period between OTP attempts) configuration. (see [below for nested schema](#nestedatt--phone--otp--deliveries--cooldown))
- `count` (Number) Allowed maximum number of OTP deliveries.

<a id="nestedatt--phone--otp--deliveries--cooldown"></a>
### Nested Schema for `phone.otp.deliveries.count`

Required:

- `duration` (Number) Cooldown duration.
    - If `cooldown.time_unit` is `MINUTES`, the allowed range is `0 - 30`.
    - If `cooldown.time_unit` is `SECONDS`, the allowed range is `0 - 1800`.
    - Defaults to `30 SECONDS`.
- `time_unit` (String) Time unit of the cooldown duration configuration.  Options are `MINUTES`, `SECONDS`.  Defaults to `SECONDS`.



<a id="nestedatt--phone--otp--lifetime"></a>
### Nested Schema for `phone.otp.lifetime`

Required:

- `duration` (Number) Lifetime of the OTP delivered via phone (SMS).
    - If `lifetime.time_unit` is `MINUTES`, the allowed range is `1 - 30`.
    - If `lifetime.time_unit` is `SECONDS`, the allowed range is `60 - 1800`.
    - Defaults to `5 MINUTES`.
- `time_unit` (String) Time unit of the OTP (SMS) duration lifetime.  Options are `MINUTES`, `SECONDS`.  Defaults to `MINUTES`.


<a id="nestedatt--phone--otp--notification"></a>
### Nested Schema for `phone.otp.notification`

Optional:

- `variant_name` (String) Name of the template variant to use to pass a one-time passcode (OTP).

Read-Only:

- `template_name` (String) Name of the template to use to pass a one-time passcode (OTP). The default value of `email_phone_verification` is static. Use the `notification.variant_name` property to define an alternate template.




<a id="nestedatt--transaction"></a>
### Nested Schema for `transaction`

Optional:

- `data_collection` (Attributes) Object for data collection timeout definition. (see [below for nested schema](#nestedatt--transaction--data_collection))
- `data_collection_only` (Boolean) When `true`, collects documents specified in the policy without determining their validity; defaults to `false`.
- `timeout` (Attributes) Object for transaction timeout. (see [below for nested schema](#nestedatt--transaction--timeout))

<a id="nestedatt--transaction--data_collection"></a>
### Nested Schema for `transaction.data_collection`

Required:

- `timeout` (Attributes) Object for data collection timeout. (see [below for nested schema](#nestedatt--transaction--data_collection--timeout))

<a id="nestedatt--transaction--data_collection--timeout"></a>
### Nested Schema for `transaction.data_collection.timeout`

Required:

- `duration` (Number) Length of time before the data collection transaction expires.
    - If `transaction.data_collection.timeout.time_unit` is `MINUTES`, the allowed range is `0 - 30`.
    - If `transaction.data_collection.timeout.time_unit` is `SECONDS`, the allowed range is `0 - 1800`.
    - Defaults to `15 MINUTES`.

    ~> When setting or changing timeouts in the transaction configuration object, `transaction.data_collection.timeout.duration` must be less than or equal to `transaction.timeout.duration`.
- `time_unit` (String) Time unit of data collection timeout.  Options are `MINUTES`, `SECONDS`.  Defaults to `MINUTES`.



<a id="nestedatt--transaction--timeout"></a>
### Nested Schema for `transaction.timeout`

Required:

- `duration` (Number) Length of time before the transaction expires.
    - If `transaction.timeout.time_unit` is `MINUTES`, the allowed range is `0 - 30`.
    - If `transaction.timeout.time_unit` is `SECONDS`, the allowed range is `0 - 1800`.
    - Defaults to `30 MINUTES`.
- `time_unit` (String) Time unit of transaction timeout.  Options are `MINUTES`, `SECONDS`.  Defaults to `MINUTES`.



<a id="nestedatt--voice"></a>
### Nested Schema for `voice`

Required:

- `comparison_threshold` (String) Comparison threshold requirements.  Options are `HIGH`, `LOW`, `MEDIUM`.  Defaults to `MEDIUM`.
- `enrollment` (Boolean) Controls if the transaction performs voice enrollment (`TRUE`) or voice verification (`FALSE`).
- `liveness_threshold` (String) Liveness threshold requirements.  Options are `HIGH`, `LOW`, `MEDIUM`.  Defaults to `MEDIUM`.
- `verify` (String) Controls the verification requirements for a Voice verification.  Options are `DISABLED`, `OPTIONAL`, `REQUIRED`.  Defaults to `DISABLED`.

Optional:

- `reference_data` (Attributes) Object for configuration of voice recording reference data. (see [below for nested schema](#nestedatt--voice--reference_data))
- `text_dependent` (Attributes) Object for configuration of text dependent voice verification. (see [below for nested schema](#nestedatt--voice--text_dependent))

<a id="nestedatt--voice--reference_data"></a>
### Nested Schema for `voice.reference_data`

Optional:

- `retain_original_recordings` (Boolean) Controls if the service stores the original voice recordings.
- `update_on_reenrollment` (Boolean) Controls updates to user's voice reference data (voice recordings) upon user re-enrollment. If `TRUE`, new data adds to existing data. If `FALSE`, new data replaces existing data.
- `update_on_verification` (Boolean) Controls updates to user's voice reference data (voice recordings) upon user verification. If `TRUE`, new data adds to existing data. If `FALSE`, new voice recordings are not retained as reference data.


<a id="nestedatt--voice--text_dependent"></a>
### Nested Schema for `voice.text_dependent`

Required:

- `samples` (Number) Number of voice samples to collect. The allowed range is `3 - 5`
- `voice_phrase_id` (String) For a customer-defined phrase, the identifier (UUID) of the voice phrase to use. For pre-defined phrases, a string value.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
$ terraform import pingone_verify_policy.example <environment_id>/<verify_policy_id>
```
