---
page_title: "pingone_notification_policy Resource - terraform-provider-pingone"
subcategory: "Platform"
description: |-
  Resource to create and manage notification policies in a PingOne environment.
---

# pingone_notification_policy (Resource)

Resource to create and manage notification policies in a PingOne environment.

## Example Usage - Unlimited Quota

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_policy" "unlimited" {
  environment_id = pingone_environment.my_environment.id

  name = "Unlimited Quotas SMS Voice and Email"
}
```

## Example Usage - Environment Quota

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_policy" "environment" {
  environment_id = pingone_environment.my_environment.id

  name = "Environment Quota SMS Voice and Email"

  quota = [
    {
      type             = "ENVIRONMENT"
      delivery_methods = ["SMS", "Voice"]
      total            = 100
    },
    {
      type             = "ENVIRONMENT"
      delivery_methods = ["Email"]
      total            = 100
    }
  ]
}
```

## Example Usage - User Quota

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_notification_policy" "user" {
  environment_id = pingone_environment.my_environment.id

  name = "User Quota SMS Voice and Email"

  quota = [
    {
      type             = "USER"
      delivery_methods = ["SMS", "Voice"]
      total            = 30
    },
    {
      type             = "USER"
      delivery_methods = ["Email"]
      total            = 30
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment to associate the notification policy with.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `name` (String) The name to use for the notification policy.  Must be unique among the notification policies in the environment.

### Optional

- `country_limit` (Attributes) A single object to limit the countries where you can send SMS and voice notifications. (see [below for nested schema](#nestedatt--country_limit))
- `quota` (Attributes Set) A set of objects that define the SMS/Voice limits.  A maximum of two quota objects can be defined, one for SMS and/or Voice quota, and one for Email quota. (see [below for nested schema](#nestedatt--quota))

### Read-Only

- `default` (Boolean) A boolean to provide an indication of whether this policy is the default notification policy for the environment. If the parameter is not provided, the value used is `false`.
- `id` (String) The ID of this resource.

<a id="nestedatt--country_limit"></a>
### Nested Schema for `country_limit`

Required:

- `type` (String) A string that specifies the kind of limitation being defined.  Options are `ALLOWED` (allows notifications only for the countries specified in the `countries` parameter), `DENIED` (denies notifications only for the countries specified in the `countries` parameter), `NONE` (no limitation is defined).

Optional:

- `countries` (Set of String) The countries where the specified methods should be allowed or denied. Use two-letter country codes from ISO 3166-1.  Required when `type` is not `NONE`.
- `delivery_methods` (Set of String) The delivery methods that the defined limitation should be applied to. Content of the array can be `SMS`, `Voice`, or both. If the parameter is not provided, the default is `SMS` and `Voice`.


<a id="nestedatt--quota"></a>
### Nested Schema for `quota`

Required:

- `type` (String) A string to specify whether the limit defined is per-user or per environment.  Options are `ENVIRONMENT`, `USER`.

Optional:

- `delivery_methods` (Set of String) The delivery methods for which the limit is being defined.  This limits defined in this block are configured as two groups, Voice/SMS, or Email.  Email cannot be configured with Voice and/or SMS limits.  Options are `Email` (configuration of Email limits but can not be set alongside `SMS` or `Voice`), `SMS` (configuration of SMS limits and can be set alongside `Voice`, but not `Email`), `Voice` (configuration of Voice limits and can be set alongside `SMS`, but not `Email`).  Defaults to `["SMS", "Voice"]`.
- `total` (Number) The maximum number of notifications allowed per day.  Cannot be set with `used` and `unused`.
- `unused` (Number) The maximum number of notifications that can be received and not responded to each day. Must be configured with `used` and cannot be configured with `total`.
- `used` (Number) The maximum number of notifications that can be received and responded to each day. Must be configured with `unused` and cannot be configured with `total`.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_notification_policy.example <environment_id>/<notification_policy_id>
```
