---
page_title: "pingone_application_secret Data Source - terraform-provider-pingone"
subcategory: "SSO"
description: |-
  Datasource to retrieve the currently active secret, and the active previous secret for a PingOne application in an environment.
---

# pingone_application_secret (Data Source)

Datasource to retrieve the currently active secret, and the active previous secret for a PingOne application in an environment.

## Example Usage

```terraform
data "pingone_application_secret" "my_awesome_oidc_application" {
  environment_id = var.environment_id
  application_id = var.application_id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `application_id` (String) PingOne application identifier (UUID) for which to retrieve the application secret.  The application must be an OpenID Connect type.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `environment_id` (String) PingOne environment identifier (UUID) in which the application exists.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.

### Read-Only

- `previous` (Attributes) An object that specifies the previous secret, when it expires, and when it was last used. (see [below for nested schema](#nestedatt--previous))
- `secret` (String, Sensitive) The application secret ID used to authenticate to the authorization server. The secret has a minimum length of 64 characters per SHA-512 requirements when using the HS512 algorithm to sign ID tokens using the secret as the key.

<a id="nestedatt--previous"></a>
### Nested Schema for `previous`

Read-Only:

- `expires_at` (String) A timestamp that specifies how long this secret is saved (and can be used) before it expires. Supported time range is 1 minute to 30 days.
- `last_used` (String) A timestamp that specifies when the previous secret was last used.
- `secret` (String, Sensitive) A string that specifies the previous application secret. This property is returned in the response if the previous secret is not expired.
