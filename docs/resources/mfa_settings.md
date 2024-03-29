---
page_title: "pingone_mfa_settings Resource - terraform-provider-pingone"
subcategory: "MFA"
description: |-
  Resource to create and manage a PingOne Environment's MFA Settings
---

# pingone_mfa_settings (Resource)

Resource to create and manage a PingOne Environment's MFA Settings

~> Only one `pingone_mfa_settings` resource should be configured for an environment.  If multiple `pingone_mfa_settings` resource definitions exist in HCL code, these are likely to conflict with each other on apply.

## Example Usage

```terraform
resource "pingone_mfa_settings" "mfa_settings" {
  environment_id = pingone_environment.my_environment.id

  pairing {
    max_allowed_devices = 5
    pairing_key_format  = "ALPHANUMERIC"
  }

  lockout {
    failure_count    = 5
    duration_seconds = 600
  }

  phone_extensions_enabled = true

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment to create the sign on policy in.
- `pairing` (Block List, Min: 1, Max: 1) An object that contains pairing settings. (see [below for nested schema](#nestedblock--pairing))

### Optional

- `authentication` (Block List, Max: 1, Deprecated) **This property is deprecated.**  Device selection settings should now be configured on the device policy, the `pingone_mfa_policy` resource. An object that contains the device selection settings. (see [below for nested schema](#nestedblock--authentication))
- `lockout` (Block List, Max: 1) An object that contains lockout settings. (see [below for nested schema](#nestedblock--lockout))
- `phone_extensions_enabled` (Boolean) A boolean when set to `true` allows one-time passwords to be delivered via voice to phone numbers that include extensions. Set to `false` to disable support for extensions. Defaults to `false`.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--pairing"></a>
### Nested Schema for `pairing`

Required:

- `pairing_key_format` (String) String that controls the type of pairing key issued. The valid values are `NUMERIC` (12-digit key) and `ALPHANUMERIC` (16-character alphanumeric key).

Optional:

- `max_allowed_devices` (Number) An integer that defines the maximum number of MFA devices each user can have. This can be any number up to 15. The default value is 5. Defaults to `5`.


<a id="nestedblock--authentication"></a>
### Nested Schema for `authentication`

Required:

- `device_selection` (String, Deprecated) **This property is deprecated.**  Device selection settings should now be configured on the device policy, the `pingone_mfa_policy` resource.  A string that defines the device selection method. Options are `DEFAULT_TO_FIRST` (this is the default setting for new environments) and `PROMPT_TO_SELECT`.


<a id="nestedblock--lockout"></a>
### Nested Schema for `lockout`

Required:

- `failure_count` (Number) An integer that defines the maximum number of incorrect authentication attempts before the account is locked.

Optional:

- `duration_seconds` (Number) An integer that defines the number of seconds to keep the account in a locked state.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
$ terraform import pingone_mfa_settings.example <environment_id>
```
