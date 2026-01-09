---
layout: ""
page_title: "Upgrade MFA Policies to use FIDO2 with Passkeys"
description: |-
  The guide describes how to upgrade MFA Policies, configured with Terraform, from using old FIDO2 policies to the new FIDO2 policies that include Passkeys functionality.
---

# Upgrade MFA Policies to use FIDO2 with Passkeys

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

### After migration
* The `pingone_mfa_policy` resource cannot be configured with the `security_key` and `platform` device types.
* The `fido2` device type should be configured on the `pingone_mfa_policy` resource when configuring FIDO2 policy.
* Old style FIDO policies (configured with the `pingone_mfa_fido_policy` resource) cannot be used.
* The new FIDO2 policies (configured with the `pingone_mfa_fido2_policy` resource) should be used.

## Upgrade Procedure

The following HCL provides an example of how to upgrade MFA device policies in an environment.

### Before Migration

Policies use the deprecated `security_key` and `platform` device types.

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

### Migration

All policies must be upgraded at the same time, using the `pingone_mfa_policies` resource.  The `data.pingone_mfa_policies` data source can be used to retrieve the IDs of all the MFA device policies in an environment, including the environment's default policy.  These IDs can then be provided to the `pingone_mfa_policies` resource as shown.

Note that the MFA device policy to upgrade must have `depends_on = [pingone_mfa_policies.mfa_policies]` set so that updates happen once the policy has been upgraded.

```terraform
data "pingone_mfa_policies" "example_all_mfa_policy_ids" {
  environment_id = pingone_environment.my_environment.id
}

resource "pingone_mfa_policies" "mfa_policies" {
  environment_id = pingone_environment.my_environment.id

  migrate_data = [
    for policy_id in data.pingone_mfa_policies.example_all_mfa_policy_ids.ids : {
      device_authentication_policy_id = policy_id
    }
  ]
}

resource "pingone_mfa_policy" "my_awesome_mfa_policy" {

  depends_on = [pingone_mfa_policies.mfa_policies]

  environment_id = pingone_environment.my_environment.id
  name           = "My awesome MFA policy"

  ...

  #   platform {
  #     enabled = true
  #   }

  #   security_key {
  #     enabled = true
  #   }

  fido2 {
    enabled = true
    fido2_policy_id = pingone_mfa_fido2_policy.my_awesome_fido2_policy.id
  }

  ...
}
```

```shell
$ terraform plan
$ terraform apply
```

The policies will only be upgraded once, but the **plan** and **apply** stages can be run multiple times as needed.

### After Migration

The `pingone_mfa_policies` resource can be safely removed.

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_policy" {
  environment_id = pingone_environment.my_environment.id
  name           = "My awesome MFA policy"

  ...

  fido2 {
    enabled = true
    fido2_policy_id = pingone_mfa_fido2_policy.my_awesome_fido2_policy.id
  }

  ...
}
```

```shell
$ terraform plan
$ terraform apply
```