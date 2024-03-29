---
page_title: "pingone_license Data Source - terraform-provider-pingone"
subcategory: "Platform"
description: |-
  Datasource to read detailed PingOne license data, selected by the license ID.
---

# pingone_license (Data Source)

Datasource to read detailed PingOne license data, selected by the license ID.

## Example Usage

```terraform
data "pingone_license" "my_license" {
  organization_id = var.organization_id
  license_id      = var.license_id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `license_id` (String) A string that specifies the license resource’s unique identifier.
- `organization_id` (String) A string that specifies the organization resource’s unique identifier associated with the license.

### Read-Only

- `advanced_services` (List of Object) A block that describes features related to **advanced services**. (see [below for nested schema](#nestedatt--advanced_services))
- `assigned_environments_count` (Number) An integer that specifies the total number of environments associated with this license.
- `authorize` (List of Object) A block that describes features related to the **authorize** services. (see [below for nested schema](#nestedatt--authorize))
- `begins_at` (String) The date and time this license begins.
- `credentials` (List of Object) A block that describes features related to the **credentials** services. (see [below for nested schema](#nestedatt--credentials))
- `environments` (List of Object) A block that describes features related to the **environments** in the organization. (see [below for nested schema](#nestedatt--environments))
- `expires_at` (String) The date and time this license expires. `TRIAL` licenses stop access to PingOne services at expiration. All other licenses trigger an event to send a notification when the license expires but do not block services.
- `fraud` (List of Object) A block that describes features related to the **fraud** services. (see [below for nested schema](#nestedatt--fraud))
- `gateways` (List of Object) A block that describes features related to the **gateway** services. (see [below for nested schema](#nestedatt--gateways))
- `id` (String) The ID of this resource.
- `intelligence` (List of Object) A block that describes features related to the **intelligence** services. (see [below for nested schema](#nestedatt--intelligence))
- `mfa` (List of Object) A block that describes features related to the **mfa** service. (see [below for nested schema](#nestedatt--mfa))
- `name` (String) A string that specifies a descriptive name for the license.
- `orchestrate` (List of Object) A block that describes features related to the **identity orchestration** services. (see [below for nested schema](#nestedatt--orchestrate))
- `package` (String) A string that specifies the license template on which this license is based. Options are `TRIAL`, `STANDARD`, `PREMIUM`, `MFA`, `RISK`, `MFARISK`, and `GLOBAL`.
- `replaced_by_license_id` (String) A string that specifies the license ID of the license that replaces this license.
- `replaces_license_id` (String) A string that specifies the license ID of the license that is replaced by this license.
- `status` (String) A string that specifies the status of the license. Options are `ACTIVE`, `EXPIRED`, and `FUTURE`.
- `terminates_at` (String) An attribute that designates the exact date and time when this license terminates access to PingOne services.
- `users` (List of Object) A block that describes features related to the **users** in the organization. (see [below for nested schema](#nestedatt--users))
- `verify` (List of Object) A block that describes features related to the **verify** services. (see [below for nested schema](#nestedatt--verify))

<a id="nestedatt--advanced_services"></a>
### Nested Schema for `advanced_services`

Read-Only:

- `pingid` (List of Object) (see [below for nested schema](#nestedobjatt--advanced_services--pingid))

<a id="nestedobjatt--advanced_services--pingid"></a>
### Nested Schema for `advanced_services.pingid`

Read-Only:

- `included` (Boolean)
- `type` (String)



<a id="nestedatt--authorize"></a>
### Nested Schema for `authorize`

Read-Only:

- `allow_api_access_management` (Boolean)
- `allow_dynamic_authorization` (Boolean)


<a id="nestedatt--credentials"></a>
### Nested Schema for `credentials`

Read-Only:

- `allow_credentials` (Boolean)


<a id="nestedatt--environments"></a>
### Nested Schema for `environments`

Read-Only:

- `allow_add_resources` (Boolean)
- `allow_connections` (Boolean)
- `allow_custom_domain` (Boolean)
- `allow_custom_schema` (Boolean)
- `allow_production` (Boolean)
- `max` (Number)
- `regions` (Set of String)


<a id="nestedatt--fraud"></a>
### Nested Schema for `fraud`

Read-Only:

- `allow_account_protection` (Boolean)
- `allow_bot_malicious_device_detection` (Boolean)


<a id="nestedatt--gateways"></a>
### Nested Schema for `gateways`

Read-Only:

- `allow_kerberos_gateway` (Boolean)
- `allow_ldap_gateway` (Boolean)
- `allow_radius_gateway` (Boolean)


<a id="nestedatt--intelligence"></a>
### Nested Schema for `intelligence`

Read-Only:

- `allow_advanced_predictors` (Boolean)
- `allow_anonymous_network_detection` (Boolean)
- `allow_data_consent` (Boolean)
- `allow_geo_velocity` (Boolean)
- `allow_reputation` (Boolean)
- `allow_risk` (Boolean)


<a id="nestedatt--mfa"></a>
### Nested Schema for `mfa`

Read-Only:

- `allow_email_otp` (Boolean)
- `allow_fido2_devices` (Boolean)
- `allow_notification_outside_whitelist` (Boolean)
- `allow_push_notification` (Boolean)
- `allow_sms_otp` (Boolean)
- `allow_totp` (Boolean)
- `allow_voice_otp` (Boolean)


<a id="nestedatt--orchestrate"></a>
### Nested Schema for `orchestrate`

Read-Only:

- `allow_orchestration` (Boolean)


<a id="nestedatt--users"></a>
### Nested Schema for `users`

Read-Only:

- `allow_identity_providers` (Boolean)
- `allow_inbound_provisioning` (Boolean)
- `allow_my_account` (Boolean)
- `allow_password_management_notifications` (Boolean)
- `allow_password_only_authentication` (Boolean)
- `allow_password_policy` (Boolean)
- `allow_provisioning` (Boolean)
- `allow_role_assignment` (Boolean)
- `allow_update_self` (Boolean)
- `allow_verification_flow` (Boolean)
- `annual_active_included` (Number)
- `entitled_to_support` (Boolean)
- `max` (Number)
- `max_hard_limit` (Number)
- `monthly_active_included` (Number)


<a id="nestedatt--verify"></a>
### Nested Schema for `verify`

Read-Only:

- `allow_document_match` (Boolean)
- `allow_face_match` (Boolean)
- `allow_manual_id_inspection` (Boolean)
- `allow_push_notifications` (Boolean)
