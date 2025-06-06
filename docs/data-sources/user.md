---
page_title: "pingone_user Data Source - terraform-provider-pingone"
subcategory: "SSO"
description: |-
  Data source to read a single PingOne user's data in an environment for a given username, email address or user ID.
---

# pingone_user (Data Source)

Data source to read a single PingOne user's data in an environment for a given username, email address or user ID.

## Example Usage

```terraform
data "pingone_user" "example_by_username" {
  environment_id = var.environment_id

  username = "user123"
}

data "pingone_user" "example_by_email" {
  environment_id = var.environment_id

  email = "user123@bxretail.org"
}

data "pingone_user" "example_by_id" {
  environment_id = var.environment_id

  user_id = var.user_id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment that the user has been created in.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.

### Optional

- `email` (String) A string that specifies the user's email address. For more information about email address formatting, see section 3.4 of [RFC 2822, Internet Message Format](http://www.faqs.org/rfcs/rfc2822.html).  Exactly one of the following must be defined: `user_id`, `username`, `email`.
- `user_id` (String) A string that specifies the ID of the user.  Must be a valid PingOne resource ID.  Exactly one of the following must be defined: `user_id`, `username`, `email`.
- `username` (String) A string that specifies the user name, which is unique within an environment. The `username` must either be a well-formed email address or a string. The string can contain any letters, numbers, combining characters, math and currency symbols, dingbats and drawing characters, and invisible whitespace.  Exactly one of the following must be defined: `user_id`, `username`, `email`.

### Read-Only

- `account` (Attributes) A single object that specifies the user's account information. (see [below for nested schema](#nestedatt--account))
- `address` (Attributes) A single object that specifies the user's address information. (see [below for nested schema](#nestedatt--address))
- `email_verified` (Boolean) A boolean that specifies whether the user's email is verified.
- `enabled` (Boolean) A boolean that specifies whether the user is enabled. This attribute is set to `true` by default when a user is created.
- `external_id` (String) A string that specifies an identifier for the user resource as defined by the provisioning client. The external id attribute simplifies the correlation of the user in PingOne with the user's account in another system of record. The platform does not use this attribute directly in any way, but it is used by Ping Identity's Data Sync product.
- `id` (String) The ID of this resource.
- `identity_provider` (Attributes) A single object that specifies the user's identity provider information. (see [below for nested schema](#nestedatt--identity_provider))
- `locale` (String) A string that specifies the user's default location as a valid language tag as defined in [RFC 5646](https://www.rfc-editor.org/rfc/rfc5646.html). The following are example tags: `fr`, `en-US`, `es-419`, `az-Arab`, `man-Nkoo-GN`. This is used for purposes of localizing such items as currency, date time format, or numerical representations.
- `mfa_enabled` (Boolean) A boolean that specifies whether multi-factor authentication is enabled. This attribute is set to `false` by default when the user is created.
- `mobile_phone` (String) A string that specifies the user's native phone number. This might also match the `primary_phone` attribute.
- `name` (Attributes) A single object that specifies the user's name information. (see [below for nested schema](#nestedatt--name))
- `nickname` (String) A string that specifies the user's nickname.
- `password` (Attributes) A single object that specifies the user's password information. (see [below for nested schema](#nestedatt--password))
- `photo` (Attributes) A single object that describes the user's photo information. (see [below for nested schema](#nestedatt--photo))
- `population_id` (String) A PingOne resource identifier of the population resource associated with the user.
- `preferred_language` (String) A string that specifies the user's preferred written or spoken languages, as a valid language range that is the same as the HTTP `Accept-Language` header field (not including `Accept-Language:` prefix) and is specified in [Section 5.3.5 of RFC 7231](https://datatracker.ietf.org/doc/html/rfc7231#section-5.3.5). For example: `en-US`, `en-gb;q=0.8`, `en;q=0.7`.
- `primary_phone` (String) A string that specifies the user's primary phone number. This might also match the `mobile_phone` attribute.
- `timezone` (String) A string that specifies the user's time zone, conforming with the IANA Time Zone database format [RFC 6557](https://www.rfc-editor.org/rfc/rfc6557.html), also known as the "Olson" time zone database format [Olson-TZ](https://www.iana.org/time-zones). For example, `America/Los_Angeles`.
- `title` (String) A string that specifies the user's title, such as `Vice President`.
- `type` (String) A string that specifies the user's type.
- `user_lifecycle` (Attributes) A single object that specifies the user's identity lifecycle information. (see [below for nested schema](#nestedatt--user_lifecycle))
- `verify_status` (String) A string that indicates whether ID verification can be done for the user.  Options are `DISABLED`, `ENABLED`, `NOT_INITIATED`.  If the user verification status is `DISABLED`, a new verification status cannot be created for that user until the status is changed to `ENABLED`.

<a id="nestedatt--account"></a>
### Nested Schema for `account`

Read-Only:

- `can_authenticate` (Boolean) A boolean that specifies whether the user can authenticate. If the value is set to `false`, the account is locked or the user is disabled, and unless specified otherwise in administrative configuration, the user will be unable to authenticate.
- `locked_at` (String) The time the specified user account was locked. This property might be absent if the account is unlocked or if the account was locked out automatically by failed password attempts.
- `status` (String) A string that specifies the the account locked state.  Options are `LOCKED`, `OK`.


<a id="nestedatt--address"></a>
### Nested Schema for `address`

Read-Only:

- `country_code` (String) A string that specifies the country name component in [ISO 3166-1](https://www.iso.org/iso-3166-country-codes.html) "alpha-2" code format. For example, the country codes for the United States and Sweden are `US` and `SE`, respectively.
- `locality` (String) A string that specifies the city or locality component of the address.
- `postal_code` (String) A string that specifies the ZIP code or postal code component of the address.
- `region` (String) A string that specifies the state, province, or region component of the address.
- `street_address` (String) A string that specifies the full street address component, which may include house number, street name, P.O. box, and multi-line extended street address information.


<a id="nestedatt--identity_provider"></a>
### Nested Schema for `identity_provider`

Read-Only:

- `id` (String) A string that identifies the external identity provider used to authenticate the user. If not provided, PingOne is the identity provider. This attribute is required if the identity provider is authoritative for just-in-time user provisioning.
- `type` (String) A string that specifies the type of identity provider used to authenticate the user.  Options are `AMAZON`, `APPLE`, `FACEBOOK`, `GITHUB`, `GOOGLE`, `LINKEDIN`, `LINKEDIN_OIDC`, `MICROSOFT`, `OPENID_CONNECT`, `PAYPAL`, `PING_ONE`, `SAML`, `TWITTER`, `YAHOO`.  The default value of `PING_ONE` is set when a value for `id` was not provided when the user was originally created.


<a id="nestedatt--name"></a>
### Nested Schema for `name`

Read-Only:

- `family` (String) A string that specifies the family name of the user, or Last in most Western languages (for example, `Jensen` given the full name `Ms. Barbara J Jensen, III`).
- `formatted` (String) A string that specifies the fully formatted name of the user (for example `Ms. Barbara J Jensen, III`).
- `given` (String) A string that specifies the given name of the user, or First in most Western languages (for example, `Barbara` given the full name `Ms. Barbara J Jensen, III`).
- `honorific_prefix` (String) A string that specifies the honorific prefix(es) of the user, or title in most Western languages (for example, `Ms.` given the full name `Ms. Barbara Jane Jensen, III`).
- `honorific_suffix` (String) A string that specifies the honorific suffix(es) of the user, or suffix in most Western languages (for example, `III` given the full name `Ms. Barbara Jane Jensen, III`).
- `middle` (String) A string that specifies the middle name(s) of the user (for exmple, `Jane` given the full name `Ms. Barbara Jane Jensen, III`).


<a id="nestedatt--password"></a>
### Nested Schema for `password`

Read-Only:

- `external` (Attributes) A single object that maps the information relevant to the user's password, and its association to external directories. (see [below for nested schema](#nestedatt--password--external))

<a id="nestedatt--password--external"></a>
### Nested Schema for `password.external`

Read-Only:

- `gateway` (Attributes) A single object that contains the external gateway properties. When this is value is specified, the user's password is managed in an external directory. (see [below for nested schema](#nestedatt--password--external--gateway))

<a id="nestedatt--password--external--gateway"></a>
### Nested Schema for `password.external.gateway`

Read-Only:

- `correlation_attributes` (Map of String) A string map that maps the external LDAP directory attributes to PingOne attributes. PingOne uses these values to read the attributes from the external LDAP directory and map them to the corresponding PingOne attributes.
- `id` (String) A string that specifies the PingOne resource ID of the linked gateway that references the remote directory.
- `type` (String) A string that indicates one of the supported gateway types.  Options are `API_GATEWAY_INTEGRATION`, `LDAP`, `PING_FEDERATE`, `PING_INTELLIGENCE`, `RADIUS`.
- `user_type_id` (String) A string that specifies the PingOne resource ID of a user type in the list of user types for the LDAP gateway.




<a id="nestedatt--photo"></a>
### Nested Schema for `photo`

Read-Only:

- `href` (String) The URI that is a uniform resource locator (as defined in [Section 1.1.3 of RFC 3986](https://www.rfc-editor.org/rfc/rfc3986#section-1.3)) that points to a resource location representing the user's image.


<a id="nestedatt--user_lifecycle"></a>
### Nested Schema for `user_lifecycle`

Read-Only:

- `status` (String) A string that specifies the status of the account lifecycle.  Options are `ACCOUNT_OK`, `VERIFICATION_REQUIRED`.
