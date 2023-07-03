---
layout: ""
page_title: "Upgrade MFA Policies for FIDO2 with Passkeys"
description: |-
  The guide describes how to upgrade MFA Policies, configured with Terraform, from using old FIDO2 policies to the new FIDO2 policies that include Passkeys functionality.
---

# Upgrade MFA Policies for FIDO2 with Passkeys

The guide describes how to upgrade MFA Policies, configured with Terraform, from using old FIDO2 policies to the new FIDO2 policies that include Passkeys functionality.

## Applicable Environments and Versions

The guide applies when all of the following are true:
* The target PingOne environment was created prior to 19th June 2023
* The target PingOne environment has MFA device policies managed by the PingOne Terraform provider version `< 1.0.0`
* The target PingOne environment has MFA device policies that specify the **Platform** or **Security Key** MFA device methods in their configuration (`security_key` and `platform` parameters in the `pingone_mfa_policy` resource)

No action is required if any of the following conditions are true:
* The target PingOne environment was created after 19th June 2023
* Existing MFA device policies specify the **FIDO2** device method (`fido2` parameter in the `pingone_mfa_policy` resource)

## Before and After Upgrade

The upgrade process is irreversable.  Once MFA device policies have been upgraded, they cannot be returned back to the pre-migration state.

### Before migration
* The `pingone_mfa_policy` resource can be configured with the `security_key` and `platform` device types.
* The `fido2` device type cannot be configured on the `pingone_mfa_policy` resource.
* Old style FIDO policies (configured with the `pingone_mfa_fido_policy` resource) can still be used.
* The new FIDO2 policies (configured with the `pingone_mfa_fido2_policy` resource) cannot be used.

Before example:
```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_policy" {
  environment_id = pingone_environment.my_environment.id
  name           = "My awesome MFA policy"

  ...

  security_key {
    enabled        = true
    fido_policy_id = pingone_mfa_fido_policy.my_awesome_fido_policy.id
  }

  platform {
    enabled        = true
    fido_policy_id = pingone_mfa_fido_policy.my_awesome_fido_policy.id
  }

  ...
}
```

### After migration
* The `pingone_mfa_policy` resource cannot be configured with the `security_key` and `platform` device types.
* The `fido2` device type should be configured on the `pingone_mfa_policy` resource when configuring FIDO2 policy.
* Old style FIDO policies (configured with the `pingone_mfa_fido_policy` resource) cannot be used.
* The new FIDO2 policies (configured with the `pingone_mfa_fido2_policy` resource) should be used.

Before example:
```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_policy" {
  environment_id = pingone_environment.my_environment.id
  name           = "My awesome MFA policy"

  ...

  fido2 {
    enabled         = true
    fido2_policy_id = pingone_mfa_fido2_policy.my_awesome_fido_policy.id
  }

  ...
}
```

## Upgrade MFA Device Policies

The following steps provide a guide as to upgrading MFA device policies in an environment.

TBC