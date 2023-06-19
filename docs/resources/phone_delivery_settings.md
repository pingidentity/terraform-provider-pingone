---
page_title: "pingone_phone_delivery_settings Resource - terraform-provider-pingone"
subcategory: "Platform"
description: |-
  Resource to create and manage SMS/voice delivery settings in a PingOne environment.
---

# pingone_phone_delivery_settings (Resource)

Resource to create and manage SMS/voice delivery settings in a PingOne environment.

## Example Usage - Custom Twilio

```terraform
resource "pingone_phone_delivery_settings" "my_awesome_custom_provider" {
  environment_id = pingone_environment.my_environment.id

  provider_type = "CUSTOM_TWILIO"

  auth_token = var.twilio_auth_token
  sid        = var.twilio_sid
}
```

## Example Usage - Custom Syniverse

```terraform
resource "pingone_phone_delivery_settings" "my_awesome_custom_provider" {
  environment_id = pingone_environment.my_environment.id

  provider_type = "CUSTOM_SYNIVERSE"

  auth_token = var.syniverse_auth_token
}
```

## Example Usage - Custom Provider

```terraform
resource "pingone_phone_delivery_settings" "my_awesome_custom_provider" {
  environment_id = pingone_environment.my_environment.id

  provider_type = "CUSTOM_PROVIDER"

  provider_custom = {
    name = "My awesome custom notifications provider"

    authentication = {
      method   = "BASIC"
      username = var.custom_provider_username
      password = var.custom_provider_password
    }

    requests = [
      {
        delivery_method = "SMS"
        method          = "GET"
        url             = "https://api.my-sms-gateway.com/send-sms.json?to=$${to}&from=$${from}&message=$${message}"
      },
      {
        delivery_method = "Voice"
        method          = "POST"
        url             = "https://api.my-voice-gateway.com/send-voice"

        headers = {
          "Content-Type" = "application/json"
        }

        body = jsonencode({
          to      = "$${to}"
          from    = "$${from}"
          message = "$${message}"
        })

        after_tag  = "</Say> <Pause length=\"1\"/>"
        before_tag = "<Say>"
      }
    ]
  }
}
```

## Example Usage - PingOne Twilio

```terraform
resource "pingone_phone_delivery_settings" "my_awesome_custom_provider" {
  environment_id = pingone_environment.my_environment.id

  provider_type = "PINGONE_TWILIO"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment to configure SMS/voice settings for.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `provider_type` (String) A string that specifies the type of the phone delivery service.  Options are `CUSTOM_PROVIDER`, `CUSTOM_SYNIVERSE`, `CUSTOM_TWILIO`, `PINGONE_TWILIO`.

### Optional

- `provider_custom` (Attributes) Required when the `provider` parameter is set to `CUSTOM_PROVIDER`.  A nested attribute with attributes that describe custom phone delivery settings. (see [below for nested schema](#nestedatt--provider_custom))
- `provider_custom_syniverse` (Attributes) Required when the `provider` parameter is set to `CUSTOM_SYNIVERSE`.  A nested attribute with attributes that describe phone delivery settings for a custom syniverse account. (see [below for nested schema](#nestedatt--provider_custom_syniverse))
- `provider_custom_twilio` (Attributes) Required when the `provider` parameter is set to `CUSTOM_TWILIO`.  A nested attribute with attributes that describe phone delivery settings for a custom Twilio account. (see [below for nested schema](#nestedatt--provider_custom_twilio))

### Read-Only

- `created_at` (String) A string that specifies the time the resource was created.
- `id` (String) The ID of this resource.
- `updated_at` (String) A string that specifies the time the resource was last updated.

<a id="nestedatt--provider_custom"></a>
### Nested Schema for `provider_custom`

Required:

- `authentication` (Attributes) A single object that provides authentication settings for authenticating to the custom service API. (see [below for nested schema](#nestedatt--provider_custom--authentication))
- `name` (String) The string that specifies the name of the custom provider used to identify in the PingOne platform.
- `requests` (Attributes Set) One or more objects that describe the outbound custom notification requests. (see [below for nested schema](#nestedatt--provider_custom--requests))

Optional:

- `numbers` (Attributes Set) One or more objects that describe the numbers to use for phone delivery. (see [below for nested schema](#nestedatt--provider_custom--numbers))

<a id="nestedatt--provider_custom--authentication"></a>
### Nested Schema for `provider_custom.authentication`

Required:

- `method` (String) The custom provider account's authentication method.  Options are `BASIC` (`username` and `password` parameters are required to be set), `BEARER` (`token` parameter is required to be set).

Optional:

- `password` (String) A string that specifies the password for the custom provider account. Required when `method` is `BASIC`
- `token` (String) A string that specifies the authentication token to use for the custom provider account. Required when `method` is `BEARER`
- `username` (String) A string that specifies the username for the custom provider account. Required when `method` is `BASIC`


<a id="nestedatt--provider_custom--requests"></a>
### Nested Schema for `provider_custom.requests`

Required:

- `delivery_method` (String) A string that specifies the notification's delivery method.  Options are `SMS`, `VOICE`.
- `method` (String) A string that specifies the type of HTTP request method.  Options are `GET`, `POST`.
- `url` (String) The provider's remote gateway or customer gateway URL.  For requests using the `POST` method, use the provider's remote gateway URL.  For requests using the `GET` method, use the provider's remote gateway URL, including the `${to}` and `${message}` mandatory variables, and the optional `${from}` variable, for example: `https://api.transmitsms.com/send-sms.json?to=${to}&from=${from}&message=${message}`

Optional:

- `after_tag` (String) For voice OTP notifications only.  A string that specifies a closing tag which is commonly used by custom providers for defining a pause between each number in the OTP number string.  Example value: `</Say> <Pause length="1"/>`
- `before_tag` (String) For voice OTP notifications only.  A string that specifies an opening tag which is commonly used by custom providers for defining a pause between each number in the OTP number string.  Possible value: `<Say>`.
- `body` (String) Optional when the `method` is `POST`.  A string that specifies the notification's request body. The body should include the `${to}` and `${message}` mandatory variables. For some vendors, the optional `${from}` variable may also be required. For example `messageType=ARN&message=${message}&phoneNumber=${to}&sender=${from}`.  In addition, you can use [dynamic variables](https://apidocs.pingidentity.com/pingone/platform/v1/api/#notifications-templates-dynamic-variables) and the following optional variables:
    - `${voice}` - the type of voice configured for notifications
    - `${locale}` - locale
    - `${otp}` - OTP
    - `${user.username}` - user's username
    - `${user.name.given}` - user's given name
    - `${user.name.family}` - user's family name
- `headers` (Map of String) A map of strings that specifies the notification's request headers, matching the format of the request body. The header should include only one of the following if the `method` is set to `POST`:
    - `content-type` = `application/x-www-form-urlencoded` (where the `body` should be form encoded)
    - `content-type` = `application/json` (where the `body` should be JSON encoded)
- `phone_number_format` (String) A string that specifies the phone number format.  Options are `FULL` (The phone number format with a leading `+` sign, in the E.164 standard format.  For example: `+14155552671`), `NUMBER_ONLY` (The phone number format without a leading `+` sign, in the E.164 standard format.  For example: `14155552671`).  Defaults to `FULL`.


<a id="nestedatt--provider_custom--numbers"></a>
### Nested Schema for `provider_custom.numbers`

Required:

- `capabilities` (Set of String) A collection of the types of phone delivery service capabilities.  Options are `SMS`, `VOICE`.
- `number` (String) A string that specifies the phone number, toll-free number or short code.
- `supported_countries` (Set of String) Specifies the `number`'s supported countries for notification recipients, depending on the phone number type.  If an SMS template has an alphanumeric `sender` ID and also has short code, the `sender` ID will be used for destination countries that support both alphanumeric senders and short codes. For Unites States and Canada that don't support alphanumeric sender IDs, a short code will be used if both an alphanumeric sender and a short code are specified.
    - `SHORT_CODE`: A collection containing a single 2-character ISO country code, for example, `US`, `GB`, `CA`.
    If the custom provider is of `type` `CUSTOM_PROVIDER`, this attribute must not be empty or null.
    For other custom provider types, if this attribute is null (empty is not supported), the specified short code `number` can only be used to dispatch notifications to United States recipient numbers.
    - `TOLL_FREE`: A collection of valid 2-character country ISO codes, for example, `US`, `GB`, `CA`.
    If the custom provider is of `type` `CUSTOM_PROVIDER`, this attribute must not be empty or null.
    For other custom provider types, if this attribute is null (empty is not supported), the specified toll-free `number` can only be used to dispatch notifications to United States recipient numbers.
    - `PHONE_NUMBER`: this attribute cannot be specified.
- `type` (String) A string that specifies the type of phone number.  Options are `PHONE_NUMBER`, `SHORT_CODE`, `TOLL_FREE`.

Optional:

- `available` (Boolean) A boolean that specifies whether the number is currently available in the provider account.
- `selected` (Boolean) A boolean that specifies whether the number is currently available in the provider account.



<a id="nestedatt--provider_custom_syniverse"></a>
### Nested Schema for `provider_custom_syniverse`

Required:

- `auth_token` (String, Sensitive) The secret key of the Syniverse account.  This field is immutable and will trigger a replace plan if changed.


<a id="nestedatt--provider_custom_twilio"></a>
### Nested Schema for `provider_custom_twilio`

Required:

- `auth_token` (String, Sensitive) The secret key of the Twilio account.  This field is immutable and will trigger a replace plan if changed.
- `sid` (String) The public ID of the Twilio account.  This field is immutable and will trigger a replace plan if changed.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
$ terraform import pingone_phone_delivery_settings.example <environment_id>/<phone_delivery_settings_id>
```