---
layout: ""
page_title: "PingOne Terraform Provider Version 1 Upgrade Guide"
description: |-
  Version 1.0.0 of the PingOne Terraform provider is a major release that introduces breaking changes to existing HCL.  This guide describes the changes that are required to upgrade v0.* PingOne Terraform provider releases to v1.0.0 onwards.
---

# PingOne Terraform Provider Version 1 Upgrade Guide

Version 1.0.0 of the PingOne Terraform provider is a major release that introduces breaking changes to existing HCL. This guide describes the changes that are required to upgrade v0.* PingOne Terraform provider releases to v1.0.0 onwards.

## Provider Configuration


## Resource: pingone_image

### `uploaded_image` computed attribute data type change

The `uploaded_image` computed attribute is now a nested object type and no longer a block type.  Where the image data is referred to in other resources (such as `pingone_application` or `pingone_branding_theme`), the variable address needs to change.

Previous configuration example:

```terraform
resource "pingone_image" "company_logo" {
  environment_id = pingone_environment.my_environment.id

  image_file_base64 = filebase64("../path/to/image.jpg")
}

resource "pingone_image" "theme_background" {
  environment_id = pingone_environment.my_environment.id

  image_file_base64 = filebase64("../path/to/background-image.jpg")
}

resource "pingone_branding_theme" "my_awesome_theme" {
  # ...

  logo {
    id   = pingone_image.company_logo.id
    href = pingone_image.company_logo.uploaded_image[0].href
  }

  background_image {
    id   = pingone_image.theme_background.id
    href = pingone_image.theme_background.uploaded_image[0].href
  }
}
```

New configuration example:

```terraform
resource "pingone_image" "company_logo" {
  environment_id = pingone_environment.my_environment.id

  image_file_base64 = filebase64("../path/to/image.jpg")
}

resource "pingone_image" "theme_background" {
  environment_id = pingone_environment.my_environment.id

  image_file_base64 = filebase64("../path/to/background-image.jpg")
}

resource "pingone_branding_theme" "my_awesome_theme" {
  # ...

  logo {
    id   = pingone_image.company_logo.id
    href = pingone_image.company_logo.uploaded_image.href
  }

  background_image {
    id   = pingone_image.theme_background.id
    href = pingone_image.theme_background.uploaded_image.href
  }
}
```

## Resource: pingone_mfa_fido_policy

This resource was previously deprecated and has been removed.  Use the `pingone_mfa_fido2_policy` resource going forward.

## Resource: pingone_mfa_policies

This resource was previously deprecated and has been removed.  Review the [Upgrade MFA Policies to use FIDO2 with Passkeys](./upgrade-mfa-policy-for-fido2) to ensure all MFA Policies are upgraded in the PingOne tenant prior to upgrading the PingOne provider version to `v1.0.0`.

## Resource: pingone_mfa_policy

Review the [Upgrade MFA Policies to use FIDO2 with Passkeys](./upgrade-mfa-policy-for-fido2) to ensure all MFA Policies are upgraded in the PingOne tenant prior to upgrading the PingOne provider version to `v1.0.0`.

### `platform` optional parameter removed

This parameter was previously deprecated and has been removed.  Use the `fido2` parameter going forward.

### `security_key` optional parameter removed

This parameter was previously deprecated and has been removed.  Use the `fido2` parameter going forward.

## Data Source: pingone_organization

### `base_url_agreement_management` computed attribute removed

This parameter was previously deprecated and has been removed.  Consider using the [PingOne Utilities module](https://registry.terraform.io/modules/pingidentity/utils/pingone/latest) going foward.

### `base_url_api` computed attribute removed

This parameter was previously deprecated and has been removed.  Consider using the [PingOne Utilities module](https://registry.terraform.io/modules/pingidentity/utils/pingone/latest) going foward.

### `base_url_apps` computed attribute removed

This parameter was previously deprecated and has been removed.  Consider using the [PingOne Utilities module](https://registry.terraform.io/modules/pingidentity/utils/pingone/latest) going foward.

### `base_url_auth` computed attribute removed

This parameter was previously deprecated and has been removed.  Consider using the [PingOne Utilities module](https://registry.terraform.io/modules/pingidentity/utils/pingone/latest) going foward.

### `base_url_console` computed attribute removed

This parameter was previously deprecated and has been removed.  Consider using the [PingOne Utilities module](https://registry.terraform.io/modules/pingidentity/utils/pingone/latest) going foward.

### `base_url_orchestrate` computed attribute removed

This parameter was previously deprecated and has been removed.  Consider using the [PingOne Utilities module](https://registry.terraform.io/modules/pingidentity/utils/pingone/latest) going foward.
