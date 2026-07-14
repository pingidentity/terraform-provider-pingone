---
page_title: "pingone_mfa_device_policy Resource - terraform-provider-pingone"
subcategory: "MFA"
description: |-
  Resource to create and manage MFA device policies for a PingOne environment.
---

# pingone_mfa_device_policy (Resource)

Resource to create and manage MFA device policies for a PingOne environment.

## Example Usage - Basic Policy

The following example enables the FIDO2 and TOTP Authenticator factors. 

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  environment_id = pingone_environment.my_environment.id
  name           = "My awesome MFA device policy"

  mobile = {
    enabled = false
  }

  totp = {
    enabled = true
  }

  fido2 = {
    enabled = true
  }

  sms = {
    enabled = false
  }

  voice = {
    enabled = false
  }

  email = {
    enabled = false
  }
}
```

## Example Usage - Mobile Authenticator

The following example configures and enables the Mobile Authenticator (using PingOne MFA SDK) and Authenticator TOTP factors. 

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_mobile_application" {
  environment_id = pingone_environment.my_environment.id
  name           = "Mobile App"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"

    mobile_app = {
      bundle_id    = "org.bxretail.mobileapp"
      package_name = "org.bxretail.mobileapp"

      passcode_refresh_seconds = 45

      integrity_detection = {
        enabled = true

        cache_duration = {
          amount = 30
          units  = "HOURS"
        }

        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = var.google_play_integrity_api_decryption_key
          verification_key  = var.google_play_integrity_api_verification_key
        }
      }
    }
  }
}

resource "pingone_mfa_application_push_credential" "example_fcm" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_mobile_application.id

  fcm = {
    google_service_account_credentials = var.google_service_account_credentials_json
  }
}

resource "pingone_mfa_application_push_credential" "example_apns" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_mobile_application.id

  apns = {
    key               = var.apns_key
    team_id           = var.apns_team_id
    token_signing_key = var.apns_token_signing_key
  }
}

resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  environment_id = pingone_environment.my_environment.id
  name           = "My awesome MFA device policy"

  depends_on = [
    pingone_mfa_application_push_credential.example_fcm,
    pingone_mfa_application_push_credential.example_apns,
  ]

  mobile = {
    enabled = true

    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 5
          time_unit = "MINUTES"
        }
      }
    }

    applications = {
      (pingone_application.my_mobile_application.id) = {

        push = {
          enabled = true
        }

        otp = {
          enabled = true
        }

        device_authorization = {
          enabled            = true
          extra_verification = "restrictive"
        }

        auto_enrollment = {
          enabled = true
        }

        integrity_detection = "restrictive"
      }
    }
  }

  totp = {
    enabled = true
  }

  fido2 = {
    enabled = true
  }

  sms = {
    enabled = false
  }

  voice = {
    enabled = false
  }

  email = {
    enabled = false
  }
}
```

## Example Usage - PingID Device Methods

The following example sets `policy_type` to `PING_ONE_ID` and configures the PingID-specific `desktop`, `yubikey`, and `oath_token` device methods, together with the PingID-only mobile application sub-fields `biometrics_enabled`, `ip_pairing_configuration`, and `new_request_duration_configuration`.

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  environment_id = pingone_environment.my_environment.id
  policy_type    = "PING_ONE_ID"

  name = "My Awesome PingID Device Policy"

  sms = {
    enabled = true
  }

  voice = {
    enabled = true
  }

  email = {
    enabled = true
  }

  mobile = {
    enabled = true
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }

    applications = {
      (pingone_application.my_native_app.id) = {
        biometrics_enabled = true
        ip_pairing_configuration = {
          any_ip_address = true
        }
        new_request_duration_configuration = {
          device_timeout = {
            duration  = 25
            time_unit = "SECONDS"
          }
          total_timeout = {
            duration  = 40
            time_unit = "SECONDS"
          }
        }
        otp = {
          enabled = true
        }
        push = {
          enabled = true
        }
        pairing_disabled = false
      }
    }
  }

  totp = {
    enabled = true
  }

  fido2 = {
    enabled = true
  }

  desktop = {
    enabled = true
    otp = {
      failure = {
        count = 5
        cool_down = {
          duration  = 3
          time_unit = "MINUTES"
        }
      }
    }
    pairing_disabled = false
    pairing_key_lifetime = {
      duration  = 10
      time_unit = "MINUTES"
    }
  }

  yubikey = {
    enabled = true
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
    pairing_disabled = false
    pairing_key_lifetime = {
      duration  = 10
      time_unit = "MINUTES"
    }
  }

  oath_token = {
    enabled = true
    otp = {
      failure = {
        count = 3
        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
    pairing_disabled = false
    pairing_key_lifetime = {
      duration  = 10
      time_unit = "MINUTES"
    }
  }
}

resource "pingone_application" "my_native_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Native App"

  enabled = true

  oidc_options = {
    type                       = "NATIVE_APP"
    grant_types                = ["CLIENT_CREDENTIALS"]
    token_endpoint_auth_method = "CLIENT_SECRET_BASIC"
  }
}

resource "pingone_environment" "my_environment" {
  # ...
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `email` (Attributes) A single object that allows configuration of email OTP device authentication policy settings. (see [below for nested schema](#nestedatt--email))
- `environment_id` (String) The ID of the environment that contains the MFA device policy to manage.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `mobile` (Attributes) A single object that allows configuration of mobile push/OTP device authentication policy settings.  This factor requires embedding the PingOne MFA SDK into a customer facing mobile application, and configuring as a Native application using the `pingone_application` resource. (see [below for nested schema](#nestedatt--mobile))
- `name` (String) A string that specifies the MFA policy's unique name within the environment.
- `sms` (Attributes) A single object that allows configuration of SMS OTP device authentication policy settings. (see [below for nested schema](#nestedatt--sms))
- `totp` (Attributes) (see [below for nested schema](#nestedatt--totp))
- `voice` (Attributes) A single object that allows configuration of voice OTP device authentication policy settings. (see [below for nested schema](#nestedatt--voice))

### Optional

- `authentication` (Attributes) A single object that allows configuration of authentication settings in the device policy. (see [below for nested schema](#nestedatt--authentication))
- `desktop` (Attributes) A single object that allows configuration of PingID desktop device authentication policy settings. Only applicable when `policy_type` is `PING_ONE_ID`. (see [below for nested schema](#nestedatt--desktop))
- `fido2` (Attributes) A single object that allows configuration of FIDO2 device authentication policy settings. (see [below for nested schema](#nestedatt--fido2))
- `ignore_user_lock` (Boolean) A boolean that, when set to `true`, allows PingOne to skip the account lock check during MFA authentication.  Defaults to `false`.
- `new_device_notification` (String) A string that defines whether a user should be notified if a new authentication method has been added to their account.  Options are `EMAIL_THEN_SMS`, `NONE`, `SMS_THEN_EMAIL`.  Defaults to `NONE`.
- `notifications_policy` (Attributes) A single object that specifies the notification policy to use for this MFA device policy. If not specified, the default notification policy for the environment will be used. (see [below for nested schema](#nestedatt--notifications_policy))
- `oath_token` (Attributes) A single object that allows configuration of OATH token device authentication policy settings. (see [below for nested schema](#nestedatt--oath_token))
- `policy_type` (String) A string that specifies the type of MFA device policy.  Options are `PING_ONE_ID`, `PING_ONE_MFA`.  Defaults to `PING_ONE_MFA`.
- `remember_me` (Attributes) A single object that specifies 'remember me' settings so that users do not have to authenticate when accessing applications from a device they have used already. (see [below for nested schema](#nestedatt--remember_me))
- `whats_app` (Attributes) A single object that allows configuration of WhatsApp OTP device authentication policy settings. To set `enabled = true`, WhatsApp sender settings must already be configured in PingOne. (see [below for nested schema](#nestedatt--whats_app))
- `yubikey` (Attributes) A single object that allows configuration of PingID Yubikey device authentication policy settings. Only applicable when `policy_type` is `PING_ONE_ID`. (see [below for nested schema](#nestedatt--yubikey))

### Read-Only

- `default` (Boolean) A boolean that specifies whether this MFA device policy is enforced as the default within the environment. When set to `true`, all other MFA device policies are `false`.  Defaults to `false`.
- `id` (String) The ID of this resource.

<a id="nestedatt--email"></a>
### Nested Schema for `email`

Required:

- `enabled` (Boolean) A boolean that specifies whether the email OTP method is enabled or disabled in the policy.

Optional:

- `otp` (Attributes) A single object that allows configuration of email OTP settings. (see [below for nested schema](#nestedatt--email--otp))
- `pairing_disabled` (Boolean) A boolean that, when set to `true`, prevents users from pairing new devices with the email OTP method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.  Defaults to `false`.
- `prompt_for_nickname_on_pairing` (Boolean) A boolean that, when set to `true`, prompts users to provide nicknames for devices during pairing.

<a id="nestedatt--email--otp"></a>
### Nested Schema for `email.otp`

Optional:

- `failure` (Attributes) A single object that allows configuration of email OTP failure settings. (see [below for nested schema](#nestedatt--email--otp--failure))
- `lifetime` (Attributes) A single object that allows configuration of email OTP lifetime settings. (see [below for nested schema](#nestedatt--email--otp--lifetime))
- `otp_length` (Number) An integer that specifies the length of the OTP that is shown to users.  Minimum length is `6` digits and maximum is `10` digits.  Defaults to `6`.

<a id="nestedatt--email--otp--failure"></a>
### Nested Schema for `email.otp.failure`

Required:

- `cool_down` (Attributes) A single object that allows configuration of email OTP failure cool down settings. (see [below for nested schema](#nestedatt--email--otp--failure--cool_down))
- `count` (Number) An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. Minimum is `1` and maximum is `7`.  Defaults to `3`.

<a id="nestedatt--email--otp--failure--cool_down"></a>
### Nested Schema for `email.otp.failure.cool_down`

Required:

- `duration` (Number) An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.



<a id="nestedatt--email--otp--lifetime"></a>
### Nested Schema for `email.otp.lifetime`

Required:

- `duration` (Number) An integer that defines the duration (number of time units) that the passcode is valid before it expires.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.




<a id="nestedatt--mobile"></a>
### Nested Schema for `mobile`

Required:

- `enabled` (Boolean) A boolean that specifies whether the mobile device method is enabled or disabled in the policy.

Optional:

- `applications` (Attributes Map) A map of objects that specifies settings for a configured Mobile Application.  The ID of the application should be configured as the map key. (see [below for nested schema](#nestedatt--mobile--applications))
- `otp` (Attributes) A single object that specifies OTP settings for mobile applications in the policy. (see [below for nested schema](#nestedatt--mobile--otp))
- `prompt_for_nickname_on_pairing` (Boolean) A boolean that, when set to `true`, prompts users to provide nicknames for devices during pairing.

<a id="nestedatt--mobile--applications"></a>
### Nested Schema for `mobile.applications`

Optional:

- `auto_enrollment` (Attributes) A single object that specifies auto enrollment settings for the application in the policy. (see [below for nested schema](#nestedatt--mobile--applications--auto_enrollment))
- `biometrics_enabled` (Boolean) A boolean that specifies whether biometric authentication methods (such as fingerprint or facial recognition) are enabled for MFA. Only applicable for PING_ONE_ID policies.
- `device_authorization` (Attributes) A single object that specifies device authorization settings for the application in the policy. (see [below for nested schema](#nestedatt--mobile--applications--device_authorization))
- `integrity_detection` (String) Controls how authentication or registration attempts should proceed if a device integrity check does not receive a response.  Options are `permissive` (if you want to allow the process to continue if a device integrity check does not receive a response), `restrictive` (if you want to block the user if a device integrity check does not receive a response).
- `ip_pairing_configuration` (Attributes) A single object that allows you to restrict device pairing to specific IP addresses. Only applicable for PING_ONE_ID policies. (see [below for nested schema](#nestedatt--mobile--applications--ip_pairing_configuration))
- `new_request_duration_configuration` (Attributes) A single object that configures timeout settings for authentication request notifications. Only applicable for PING_ONE_ID policies. (see [below for nested schema](#nestedatt--mobile--applications--new_request_duration_configuration))
- `otp` (Attributes) A single object that specifies OTP settings for the application in the policy. (see [below for nested schema](#nestedatt--mobile--applications--otp))
- `pairing_disabled` (Boolean) A boolean that, when set to `true`, prevents users from pairing new devices with the relevant application. You can use this option if you want to phase out an existing mobile application but want to allow users to continue using the application for authentication for existing devices.
- `pairing_key_lifetime` (Attributes) A single object that specifies pairing key lifetime settings for the application in the policy. (see [below for nested schema](#nestedatt--mobile--applications--pairing_key_lifetime))
- `push` (Attributes) A single object that specifies push settings for the application in the policy. (see [below for nested schema](#nestedatt--mobile--applications--push))
- `push_limit` (Attributes) A single object that specifies push limit settings for the application in the policy. (see [below for nested schema](#nestedatt--mobile--applications--push_limit))
- `push_timeout` (Attributes) A single object that specifies push timeout settings for the application in the policy. (see [below for nested schema](#nestedatt--mobile--applications--push_timeout))

<a id="nestedatt--mobile--applications--auto_enrollment"></a>
### Nested Schema for `mobile.applications.auto_enrollment`

Required:

- `enabled` (Boolean) A boolean that, when set to `true` if you want the application to allow Auto Enrollment. Auto Enrollment means that the user can authenticate for the first time from an unpaired device, and the successful authentication will result in the pairing of the device for MFA.


<a id="nestedatt--mobile--applications--device_authorization"></a>
### Nested Schema for `mobile.applications.device_authorization`

Required:

- `enabled` (Boolean) Specifies the enabled or disabled state of automatic MFA for native devices paired with the user, for the specified application.

Optional:

- `extra_verification` (String) Specifies the level of further verification when device authorization is enabled. The PingOne platform performs an extra verification check by sending a "silent" push notification to the customer native application, and receives a confirmation in return.  By default, the PingOne platform does not perform the extra verification check.  Options are `permissive` (the PingOne platform performs the extra verification check. Upon timeout or failure to get a response from the native app, the MFA step is treated as successfully completed), `restrictive` (the PingOne platform performs the extra verification check. Upon timeout or failure to get a response from the native app, the MFA step is treated as failed).


<a id="nestedatt--mobile--applications--ip_pairing_configuration"></a>
### Nested Schema for `mobile.applications.ip_pairing_configuration`

Optional:

- `any_ip_address` (Boolean) A boolean that, when set to `false`, restricts device pairing to specific IP addresses defined in `only_these_ip_addresses`.  Defaults to `true`.
- `only_these_ip_addresses` (Set of String) A list of IP addresses or address ranges from which users can pair their devices. This parameter is required when `any_ip_address` is set to `false`. Each item in the array must be in CIDR notation, for example, `192.168.1.1/32` or `10.0.0.0/8`.


<a id="nestedatt--mobile--applications--new_request_duration_configuration"></a>
### Nested Schema for `mobile.applications.new_request_duration_configuration`

Required:

- `device_timeout` (Attributes) A single object that specifies the maximum time a notification can remain pending before it is displayed to the user. Value must be between `15` and `75` seconds.  Defaults to `25`. (see [below for nested schema](#nestedatt--mobile--applications--new_request_duration_configuration--device_timeout))
- `total_timeout` (Attributes) A single object that specifies the total time an authentication request notification has to be handled by the user before timing out. The `total_timeout.duration` must exceed `device_timeout.duration` by at least 15 seconds.  Value must be between `30` and `90` seconds.  Defaults to `40`. (see [below for nested schema](#nestedatt--mobile--applications--new_request_duration_configuration--total_timeout))

<a id="nestedatt--mobile--applications--new_request_duration_configuration--device_timeout"></a>
### Nested Schema for `mobile.applications.new_request_duration_configuration.device_timeout`

Optional:

- `duration` (Number) An integer that specifies the timeout duration in seconds.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Defaults to `SECONDS`.  Options are `SECONDS`.


<a id="nestedatt--mobile--applications--new_request_duration_configuration--total_timeout"></a>
### Nested Schema for `mobile.applications.new_request_duration_configuration.total_timeout`

Optional:

- `duration` (Number) An integer that specifies the timeout duration in seconds.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Defaults to `SECONDS`.  Options are `SECONDS`.



<a id="nestedatt--mobile--applications--otp"></a>
### Nested Schema for `mobile.applications.otp`

Required:

- `enabled` (Boolean) A boolean that specifies whether OTP authentication is enabled or disabled for the application in the policy.


<a id="nestedatt--mobile--applications--pairing_key_lifetime"></a>
### Nested Schema for `mobile.applications.pairing_key_lifetime`

Required:

- `duration` (Number) An integer that defines the amount of time an issued pairing key can be used until it expires. Minimum is 1 minute and maximum is 48 hours. If this parameter is not provided, the duration is set to 10 minutes.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `HOURS`, `MINUTES`.


<a id="nestedatt--mobile--applications--push"></a>
### Nested Schema for `mobile.applications.push`

Required:

- `enabled` (Boolean) A boolean that specifies whether push notification is enabled or disabled for the application in the policy.

Optional:

- `number_matching` (Attributes) A single object that configures number matching for push notifications. (see [below for nested schema](#nestedatt--mobile--applications--push--number_matching))

<a id="nestedatt--mobile--applications--push--number_matching"></a>
### Nested Schema for `mobile.applications.push.number_matching`

Required:

- `enabled` (Boolean) A boolean that, when set to `true`, requires the authenticating user to select a number that was displayed to them on the accessing device.



<a id="nestedatt--mobile--applications--push_limit"></a>
### Nested Schema for `mobile.applications.push_limit`

Optional:

- `count` (Number) An integer that specifies the number of consecutive push notifications that can be ignored or rejected by a user within a defined period before push notifications are blocked for the application. The minimum value is "1" and the maximum value is "50". If this parameter is not provided, the default value is "5".
- `lock_duration` (Attributes) A single object that specifies push limit lock duration settings for the application in the policy. (see [below for nested schema](#nestedatt--mobile--applications--push_limit--lock_duration))
- `time_period` (Attributes) A single object that specifies push limit time period settings for the application in the policy. (see [below for nested schema](#nestedatt--mobile--applications--push_limit--time_period))

<a id="nestedatt--mobile--applications--push_limit--lock_duration"></a>
### Nested Schema for `mobile.applications.push_limit.lock_duration`

Required:

- `duration` (Number) An integer that defines the length of time that push notifications should be blocked for the application if the defined limit has been reached. The minimum value is `1` minute and the maximum value is `120` minutes. If this parameter is not provided, the default value is `30` minutes.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.  Defaults to `MINUTES`.


<a id="nestedatt--mobile--applications--push_limit--time_period"></a>
### Nested Schema for `mobile.applications.push_limit.time_period`

Required:

- `duration` (Number) An integer that defines the length of time that push notifications should be blocked for the application if the defined limit has been reached. The minimum value is `1` minute and the maximum value is `120` minutes. If this parameter is not provided, the default value is `30` minutes.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.  Defaults to `MINUTES`.



<a id="nestedatt--mobile--applications--push_timeout"></a>
### Nested Schema for `mobile.applications.push_timeout`

Required:

- `duration` (Number) An integer that defines the length of time that push notifications should be blocked for the application if the defined limit has been reached. The minimum value is `1` minute and the maximum value is `120` minutes. If this parameter is not provided, the default value is `30` minutes.

Optional:

- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Defaults to `SECONDS`.  Options are `SECONDS`.



<a id="nestedatt--mobile--otp"></a>
### Nested Schema for `mobile.otp`

Required:

- `failure` (Attributes) A single object that specifies OTP failure settings for mobile applications in the policy. (see [below for nested schema](#nestedatt--mobile--otp--failure))

<a id="nestedatt--mobile--otp--failure"></a>
### Nested Schema for `mobile.otp.failure`

Required:

- `cool_down` (Attributes) A single object that specifies OTP failure cool down settings for mobile applications in the policy. (see [below for nested schema](#nestedatt--mobile--otp--failure--cool_down))

Optional:

- `count` (Number) An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. The minimum value is `1`, maximum is `7`, and the default is `3`.

<a id="nestedatt--mobile--otp--failure--cool_down"></a>
### Nested Schema for `mobile.otp.failure.cool_down`

Required:

- `duration` (Number) An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures. The minimum value is `2`, maximum is `30`, and the default is `2`. Note that when using the "onetime authentication" feature, the user is not blocked after the maximum number of failures even if you specified a block duration.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.  Defaults to `MINUTES`.





<a id="nestedatt--sms"></a>
### Nested Schema for `sms`

Required:

- `enabled` (Boolean) A boolean that specifies whether the SMS OTP method is enabled or disabled in the policy.

Optional:

- `otp` (Attributes) A single object that allows configuration of SMS OTP settings. (see [below for nested schema](#nestedatt--sms--otp))
- `pairing_disabled` (Boolean) A boolean that, when set to `true`, prevents users from pairing new devices with the SMS OTP method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.  Defaults to `false`.
- `prompt_for_nickname_on_pairing` (Boolean) A boolean that, when set to `true`, prompts users to provide nicknames for devices during pairing.

<a id="nestedatt--sms--otp"></a>
### Nested Schema for `sms.otp`

Optional:

- `failure` (Attributes) A single object that allows configuration of SMS OTP failure settings. (see [below for nested schema](#nestedatt--sms--otp--failure))
- `lifetime` (Attributes) A single object that allows configuration of SMS OTP lifetime settings. (see [below for nested schema](#nestedatt--sms--otp--lifetime))
- `otp_length` (Number) An integer that specifies the length of the OTP that is shown to users.  Minimum length is `6` digits and maximum is `10` digits.  Defaults to `6`.

<a id="nestedatt--sms--otp--failure"></a>
### Nested Schema for `sms.otp.failure`

Required:

- `cool_down` (Attributes) A single object that allows configuration of SMS OTP failure cool down settings. (see [below for nested schema](#nestedatt--sms--otp--failure--cool_down))
- `count` (Number) An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. Minimum is `1` and maximum is `7`.  Defaults to `3`.

<a id="nestedatt--sms--otp--failure--cool_down"></a>
### Nested Schema for `sms.otp.failure.cool_down`

Required:

- `duration` (Number) An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.



<a id="nestedatt--sms--otp--lifetime"></a>
### Nested Schema for `sms.otp.lifetime`

Required:

- `duration` (Number) An integer that defines the duration (number of time units) that the passcode is valid before it expires.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.




<a id="nestedatt--totp"></a>
### Nested Schema for `totp`

Required:

- `enabled` (Boolean) A boolean that specifies whether the TOTP method is enabled or disabled in the policy.

Optional:

- `otp` (Attributes) A single object that allows configuration of TOTP OTP settings. (see [below for nested schema](#nestedatt--totp--otp))
- `pairing_disabled` (Boolean) A boolean that, when set to `true`, prevents users from pairing new devices with the TOTP method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.  Defaults to `false`.
- `passcode_grace_period` (Number) An integer that specifies the TOTP passcode grace period in 30-second windows. The minimum value is `1` and the maximum value is `10`.  Defaults to `5`.
- `prompt_for_nickname_on_pairing` (Boolean) A boolean that, when set to `true`, prompts users to provide nicknames for devices during pairing.
- `uri_parameters` (Map of String) A map of string key:value pairs that specifies `otpauth` URI parameters. For example, if you provide a value for the `issuer` parameter, then authenticators that support that parameter will display the text you specify together with the OTP (in addition to the username). This can help users recognize which application the OTP is for. If you intend on using the same MFA policy for multiple applications, choose a name that reflects the group of applications.

<a id="nestedatt--totp--otp"></a>
### Nested Schema for `totp.otp`

Optional:

- `failure` (Attributes) A single object that allows configuration of TOTP OTP failure settings. (see [below for nested schema](#nestedatt--totp--otp--failure))

<a id="nestedatt--totp--otp--failure"></a>
### Nested Schema for `totp.otp.failure`

Required:

- `count` (Number) An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked.

Optional:

- `cool_down` (Attributes) A single object that allows configuration of TOTP OTP failure cool down settings. (see [below for nested schema](#nestedatt--totp--otp--failure--cool_down))

<a id="nestedatt--totp--otp--failure--cool_down"></a>
### Nested Schema for `totp.otp.failure.cool_down`

Required:

- `duration` (Number) An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.  Defaults to `MINUTES`.





<a id="nestedatt--voice"></a>
### Nested Schema for `voice`

Required:

- `enabled` (Boolean) A boolean that specifies whether the voice OTP method is enabled or disabled in the policy.

Optional:

- `otp` (Attributes) A single object that allows configuration of voice OTP settings. (see [below for nested schema](#nestedatt--voice--otp))
- `pairing_disabled` (Boolean) A boolean that, when set to `true`, prevents users from pairing new devices with the voice OTP method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.  Defaults to `false`.
- `prompt_for_nickname_on_pairing` (Boolean) A boolean that, when set to `true`, prompts users to provide nicknames for devices during pairing.

<a id="nestedatt--voice--otp"></a>
### Nested Schema for `voice.otp`

Optional:

- `failure` (Attributes) A single object that allows configuration of voice OTP failure settings. (see [below for nested schema](#nestedatt--voice--otp--failure))
- `lifetime` (Attributes) A single object that allows configuration of voice OTP lifetime settings. (see [below for nested schema](#nestedatt--voice--otp--lifetime))
- `otp_length` (Number) An integer that specifies the length of the OTP that is shown to users.  Minimum length is `6` digits and maximum is `10` digits.  Defaults to `6`.

<a id="nestedatt--voice--otp--failure"></a>
### Nested Schema for `voice.otp.failure`

Required:

- `cool_down` (Attributes) A single object that allows configuration of voice OTP failure cool down settings. (see [below for nested schema](#nestedatt--voice--otp--failure--cool_down))
- `count` (Number) An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. Minimum is `1` and maximum is `7`.  Defaults to `3`.

<a id="nestedatt--voice--otp--failure--cool_down"></a>
### Nested Schema for `voice.otp.failure.cool_down`

Required:

- `duration` (Number) An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.



<a id="nestedatt--voice--otp--lifetime"></a>
### Nested Schema for `voice.otp.lifetime`

Required:

- `duration` (Number) An integer that defines the duration (number of time units) that the passcode is valid before it expires.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.




<a id="nestedatt--authentication"></a>
### Nested Schema for `authentication`

Optional:

- `device_selection` (String) A string that defines the device selection method.  Options are `ALWAYS_DISPLAY_DEVICES`, `DEFAULT_TO_FIRST`, `PROMPT_TO_SELECT`.  Defaults to `DEFAULT_TO_FIRST`.


<a id="nestedatt--desktop"></a>
### Nested Schema for `desktop`

Required:

- `enabled` (Boolean) A boolean that specifies whether the desktop device method is enabled or disabled in the policy.

Optional:

- `otp` (Attributes) A single object that specifies OTP failure settings for desktop devices. (see [below for nested schema](#nestedatt--desktop--otp))
- `pairing_disabled` (Boolean) A boolean that, when set to `true`, prevents users from pairing new desktop devices.
- `pairing_key_lifetime` (Attributes) A single object that specifies pairing key lifetime settings for desktop devices. (see [below for nested schema](#nestedatt--desktop--pairing_key_lifetime))
- `prompt_for_nickname_on_pairing` (Boolean) A boolean that, when set to `true`, prompts users to provide nicknames for devices during pairing.

<a id="nestedatt--desktop--otp"></a>
### Nested Schema for `desktop.otp`

Optional:

- `failure` (Attributes) A single object that allows configuration of OTP failure settings. (see [below for nested schema](#nestedatt--desktop--otp--failure))

<a id="nestedatt--desktop--otp--failure"></a>
### Nested Schema for `desktop.otp.failure`

Optional:

- `cool_down` (Attributes) A single object that specifies OTP failure cool down settings. (see [below for nested schema](#nestedatt--desktop--otp--failure--cool_down))
- `count` (Number) An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. Must be between 1 and 7.

<a id="nestedatt--desktop--otp--failure--cool_down"></a>
### Nested Schema for `desktop.otp.failure.cool_down`

Required:

- `duration` (Number) An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures. Must be between `1` seconds and `30` minutes.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.  Defaults to `MINUTES`.




<a id="nestedatt--desktop--pairing_key_lifetime"></a>
### Nested Schema for `desktop.pairing_key_lifetime`

Required:

- `duration` (Number) An integer that defines the amount of time an issued pairing key can be used until it expires. Must be between 1 minutes and 48 hours.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `HOURS`, `MINUTES`.



<a id="nestedatt--fido2"></a>
### Nested Schema for `fido2`

Required:

- `enabled` (Boolean) A boolean that specifies whether the FIDO2 method is enabled or disabled in the policy.

Optional:

- `failure` (Attributes) A single object that allows configuration of FIDO2 authentication failure settings. (see [below for nested schema](#nestedatt--fido2--failure))
- `fido2_policy_id` (String) A string that specifies the resource UUID that represents the FIDO2 policy in PingOne. This property can be null / left undefined. When null, the environment's default FIDO2 Policy is used.  Must be a valid PingOne resource ID.
- `pairing_disabled` (Boolean) A boolean that, when set to `true`, prevents users from pairing new devices with the FIDO2 method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.  Defaults to `false`.
- `prompt_for_nickname_on_pairing` (Boolean) A boolean that, when set to `true`, prompts users to provide nicknames for devices during pairing.

<a id="nestedatt--fido2--failure"></a>
### Nested Schema for `fido2.failure`

Optional:

- `cool_down` (Attributes) A single object that allows configuration of FIDO2 authentication failure cool down settings. (see [below for nested schema](#nestedatt--fido2--failure--cool_down))
- `count` (Number) An integer that defines the maximum number of times that authentication can fail before the user is blocked. The minimum value is `1` and the maximum value is `7`.  Defaults to `3`.

<a id="nestedatt--fido2--failure--cool_down"></a>
### Nested Schema for `fido2.failure.cool_down`

Required:

- `duration` (Number) An integer that defines the length of time that the user is blocked after reaching the maximum number of failures. The minimum value is `2` minutes and the maximum value is `30` minutes.  Defaults to `2`.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.  Defaults to `MINUTES`.




<a id="nestedatt--notifications_policy"></a>
### Nested Schema for `notifications_policy`

Required:

- `id` (String) A string that specifies the ID of the notification policy to use.


<a id="nestedatt--oath_token"></a>
### Nested Schema for `oath_token`

Optional:

- `enabled` (Boolean) A boolean that specifies whether the OATH token device method is enabled or disabled in the policy.
- `otp` (Attributes) A single object that specifies OTP failure settings for OATH token devices. (see [below for nested schema](#nestedatt--oath_token--otp))
- `pairing_disabled` (Boolean) A boolean that, when set to `true`, prevents users from pairing new OATH token devices.
- `pairing_key_lifetime` (Attributes) A single object that specifies pairing key lifetime settings for OATH token devices. (see [below for nested schema](#nestedatt--oath_token--pairing_key_lifetime))
- `prompt_for_nickname_on_pairing` (Boolean) A boolean that, when set to `true`, prompts users to provide nicknames for devices during pairing.

<a id="nestedatt--oath_token--otp"></a>
### Nested Schema for `oath_token.otp`

Optional:

- `failure` (Attributes) A single object that allows configuration of OTP failure settings. (see [below for nested schema](#nestedatt--oath_token--otp--failure))

<a id="nestedatt--oath_token--otp--failure"></a>
### Nested Schema for `oath_token.otp.failure`

Optional:

- `cool_down` (Attributes) A single object that specifies OTP failure cool down settings. (see [below for nested schema](#nestedatt--oath_token--otp--failure--cool_down))
- `count` (Number) An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. Must be between `1` and `7`.

<a id="nestedatt--oath_token--otp--failure--cool_down"></a>
### Nested Schema for `oath_token.otp.failure.cool_down`

Required:

- `duration` (Number) An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures. Must be between `1` seconds and `30` minutes.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.  Defaults to `MINUTES`.




<a id="nestedatt--oath_token--pairing_key_lifetime"></a>
### Nested Schema for `oath_token.pairing_key_lifetime`

Required:

- `duration` (Number) An integer that defines the amount of time an issued pairing key can be used until it expires. Must be between `1` minutes and `48` hours.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `HOURS`, `MINUTES`.



<a id="nestedatt--remember_me"></a>
### Nested Schema for `remember_me`

Required:

- `web` (Attributes) A single object that contains the 'remember me' settings for accessing applications from a browser. (see [below for nested schema](#nestedatt--remember_me--web))

<a id="nestedatt--remember_me--web"></a>
### Nested Schema for `remember_me.web`

Required:

- `enabled` (Boolean) A boolean that, when set to `true`, enables the 'remember me' option in the MFA policy.
- `life_time` (Attributes) A single object that defines the period during which users will not have to authenticate if they are accessing applications from a device they have used before. The 'remember me' period can be anywhere from `1` minute to `90` days. (see [below for nested schema](#nestedatt--remember_me--web--life_time))

<a id="nestedatt--remember_me--web--life_time"></a>
### Nested Schema for `remember_me.web.life_time`

Required:

- `duration` (Number) An integer that, used in conjunction with `time_unit`, defines the 'remember me' period.
- `time_unit` (String) A string that specifies the time unit to use for the 'remember me' period.  Options are `DAYS`, `HOURS`, `MINUTES`.




<a id="nestedatt--whats_app"></a>
### Nested Schema for `whats_app`

Required:

- `enabled` (Boolean) A boolean that specifies whether the WhatsApp OTP method is enabled or disabled in the policy.
- `otp` (Attributes) A single object that allows configuration of WhatsApp OTP settings. (see [below for nested schema](#nestedatt--whats_app--otp))

Optional:

- `pairing_disabled` (Boolean) A boolean that, when set to `true`, prevents users from pairing new devices with the WhatsApp OTP method, though keeping it active in the policy for existing users. You can use this option if you want to phase out an existing authentication method but want to allow users to continue using the method for authentication for existing devices.
- `prompt_for_nickname_on_pairing` (Boolean) A boolean that, when set to `true`, prompts users to provide nicknames for devices during pairing.

<a id="nestedatt--whats_app--otp"></a>
### Nested Schema for `whats_app.otp`

Required:

- `failure` (Attributes) A single object that allows configuration of WhatsApp OTP failure settings. (see [below for nested schema](#nestedatt--whats_app--otp--failure))
- `lifetime` (Attributes) A single object that allows configuration of WhatsApp OTP lifetime settings. (see [below for nested schema](#nestedatt--whats_app--otp--lifetime))

Optional:

- `otp_length` (Number) An integer that specifies the length of the OTP that is shown to users.  Minimum length is `6` digits and maximum is `10` digits.

<a id="nestedatt--whats_app--otp--failure"></a>
### Nested Schema for `whats_app.otp.failure`

Required:

- `cool_down` (Attributes) A single object that allows configuration of WhatsApp OTP failure cool down settings. (see [below for nested schema](#nestedatt--whats_app--otp--failure--cool_down))
- `count` (Number) An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. Minimum is `1` and maximum is `7`.

<a id="nestedatt--whats_app--otp--failure--cool_down"></a>
### Nested Schema for `whats_app.otp.failure.cool_down`

Required:

- `duration` (Number) An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.



<a id="nestedatt--whats_app--otp--lifetime"></a>
### Nested Schema for `whats_app.otp.lifetime`

Required:

- `duration` (Number) An integer that defines the duration (number of time units) that the passcode is valid before it expires.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.




<a id="nestedatt--yubikey"></a>
### Nested Schema for `yubikey`

Required:

- `enabled` (Boolean) A boolean that specifies whether the Yubikey device method is enabled or disabled in the policy.

Optional:

- `otp` (Attributes) A single object that specifies OTP failure settings for Yubikey devices. (see [below for nested schema](#nestedatt--yubikey--otp))
- `pairing_disabled` (Boolean) A boolean that, when set to `true`, prevents users from pairing new Yubikey devices.
- `pairing_key_lifetime` (Attributes) A single object that specifies pairing key lifetime settings for Yubikey devices. (see [below for nested schema](#nestedatt--yubikey--pairing_key_lifetime))
- `prompt_for_nickname_on_pairing` (Boolean) A boolean that, when set to `true`, prompts users to provide nicknames for devices during pairing.

<a id="nestedatt--yubikey--otp"></a>
### Nested Schema for `yubikey.otp`

Optional:

- `failure` (Attributes) A single object that allows configuration of OTP failure settings. (see [below for nested schema](#nestedatt--yubikey--otp--failure))

<a id="nestedatt--yubikey--otp--failure"></a>
### Nested Schema for `yubikey.otp.failure`

Optional:

- `cool_down` (Attributes) A single object that specifies OTP failure cool down settings. (see [below for nested schema](#nestedatt--yubikey--otp--failure--cool_down))
- `count` (Number) An integer that defines the maximum number of times that the OTP entry can fail for a user, before they are blocked. Must be between 1 and 7.

<a id="nestedatt--yubikey--otp--failure--cool_down"></a>
### Nested Schema for `yubikey.otp.failure.cool_down`

Required:

- `duration` (Number) An integer that defines the duration (number of time units) the user is blocked after reaching the maximum number of passcode failures. Must be between `1` seconds and `30` minutes.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `MINUTES`, `SECONDS`.  Defaults to `MINUTES`.




<a id="nestedatt--yubikey--pairing_key_lifetime"></a>
### Nested Schema for `yubikey.pairing_key_lifetime`

Required:

- `duration` (Number) An integer that defines the amount of time an issued pairing key can be used until it expires. Must be between 1 minutes and 48 hours.
- `time_unit` (String) A string that specifies the type of time unit for `duration`.  Options are `HOURS`, `MINUTES`.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_mfa_device_policy.example <environment_id>/<mfa_device_policy_id>
```
