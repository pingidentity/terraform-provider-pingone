---
layout: ""
page_title: "PingOne Terraform Provider Version 1 Upgrade Guide"
description: |-
  Version 1.0.0 of the PingOne Terraform provider is a major release that introduces breaking changes to existing HCL.  This guide describes the changes that are required to upgrade v0.* PingOne Terraform provider releases to v1.0.0 onwards.
---

# PingOne Terraform Provider Version 1 Upgrade Guide

Version 1.0.0 of the PingOne Terraform provider is a major release that introduces breaking changes to existing HCL. This guide describes the changes that are required to upgrade v0.* PingOne Terraform provider releases to v1.0.0 onwards.

## Provider Configuration


## Resource: pingone_application

### `oidc_options.bundle_id` optional parameter removed

This parameter was previously deprecated and has been removed.  Use the `oidc_options.mobile_app.bundle_id` parameter going forward.

### `oidc_options.package_name` optional parameter removed

This parameter was previously deprecated and has been removed.  Use the `oidc_options.mobile_app.package_name` parameter going forward.

### `saml_options.idp_signing_key_id` optional parameter removed

This parameter was previously deprecated and has been removed.  Use the `saml_options.idp_signing_key` parameter going forward.  When using the `saml_options.idp_signing_key` object parameter, the `saml_options.idp_signing_key.algorithm` now also needs to be defined.

Previous configuration example:

```terraform
resource "pingone_application" "my_awesome_saml_app" {
  # ... other configuration parameters

  saml_options {
    # ... other configuration parameters

    idp_signing_key_id = pingone_key.my_awesome_key.id
  }
}
```

New configuration example:

```terraform
resource "pingone_application" "my_awesome_saml_app" {
  # ... other configuration parameters

  saml_options = {
    # ... other configuration parameters

    idp_signing_key = {
      key_id    = pingone_key.my_awesome_key.id
      algorithm = pingone_key.my_awesome_key.signature_algorithm
    }
  }
}
```

## Resource: pingone_application_resource_grant

### `resource_id` parameter changed

This parameter was previously deprecated and has now been made read only.  Use the `resource_name` parameter going forward.

### `resource_name` parameter changed

This parameter was previously optional and has now been made a required parameter.

### `scopes` parameter changed

This parameter was previously deprecated and has now been made read only.  Use the `scope_names` parameter going forward.

### `scope_names` parameter changed

This parameter was previously optional and has now been made a required parameter.

## Resource: pingone_branding_settings

### `logo_image` parameter data type change

The `logo_image` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_branding_settings" "branding" {
  # ... other configuration parameters

  logo_image {
    id   = pingone_image.company_logo.id
    href = pingone_image.company_logo.uploaded_image[0].href
  }
}
```

New configuration example:

```terraform
resource "pingone_branding_settings" "branding" {
  # ... other configuration parameters

  logo_image = {
    id   = pingone_image.company_logo.id
    href = pingone_image.company_logo.uploaded_image.href
  }
}
```

## Resource: pingone_branding_theme

### `background_image` parameter data type change

The `background_image` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_branding_theme" "my_awesome_theme" {
  # ... other configuration parameters

  background_image {
    id   = pingone_image.company_image.id
    href = pingone_image.company_image.uploaded_image[0].href
  }
}
```

New configuration example:

```terraform
resource "pingone_branding_theme" "my_awesome_theme" {
  # ... other configuration parameters

  background_image = {
    id   = pingone_image.company_image.id
    href = pingone_image.company_image.uploaded_image.href
  }
}
```

### `logo` parameter data type change

The `logo` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_branding_theme" "my_awesome_theme" {
  # ... other configuration parameters

  logo {
    id   = pingone_image.company_logo.id
    href = pingone_image.company_logo.uploaded_image[0].href
  }
}
```

New configuration example:

```terraform
resource "pingone_branding_theme" "my_awesome_theme" {
  # ... other configuration parameters

  logo = {
    id   = pingone_image.company_logo.id
    href = pingone_image.company_logo.uploaded_image.href
  }
}
```

## Resource: pingone_environment

### `default_population` optional parameter removed

This parameter was previously deprecated and has been removed.  Default populations are managed with the `pingone_population_default` resource.

### `default_population_id` computed attribute removed

This attribute was previously deprecated and has been removed.  Default populations are managed with the `pingone_population_default` resource.

### `timeouts` block removed

This parameter block is no longer needed and has been removed.

## Resource: pingone_identity_provider

### `icon` parameter data type change

The `icon` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  icon {
    id   = pingone_image.identity_provider_icon.id
    href = pingone_image.identity_provider_icon.uploaded_image[0].href
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  icon = {
    id   = pingone_image.identity_provider_icon.id
    href = pingone_image.identity_provider_icon.uploaded_image.href
  }
}
```

### `login_button_icon` parameter data type change

The `login_button_icon` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  login_button_icon {
    id   = pingone_image.identity_provider_login_button_icon.id
    href = pingone_image.identity_provider_login_button_icon.uploaded_image[0].href
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  login_button_icon = {
    id   = pingone_image.identity_provider_login_button_icon.id
    href = pingone_image.identity_provider_login_button_icon.uploaded_image.href
  }
}
```

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

  logo = {
    id   = pingone_image.company_logo.id
    href = pingone_image.company_logo.uploaded_image.href
  }

  background_image = {
    id   = pingone_image.theme_background.id
    href = pingone_image.theme_background.uploaded_image.href
  }
}
```

## Resource: pingone_mfa_application_push_credential

### `apns` schema type change

This parameter `apns` was previously a block data type, and is now a single nested object type.

Previous configuration example:

```terraform
resource "pingone_mfa_application_push_credential" "example_apns" {
  # ... other configuration parameters

  apns {
    key               = var.apns_key
    team_id           = var.apns_team_id
    token_signing_key = var.apns_token_signing_key
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_application_push_credential" "example_apns" {
  # ... other configuration parameters

  apns = {
    key               = var.apns_key
    team_id           = var.apns_team_id
    token_signing_key = var.apns_token_signing_key
  }
}
```

### `fcm` schema type change

This parameter `fcm` was previously a block data type, and is now a single nested object type.

Previous configuration example:

```terraform
resource "pingone_mfa_application_push_credential" "example_fcm" {
  # ... other configuration parameters

  fcm {
    google_service_account_credentials = var.google_service_account_credentials_json
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_application_push_credential" "example_fcm" {
  # ... other configuration parameters

  fcm = {
    google_service_account_credentials = var.google_service_account_credentials_json
  }
}
```

### `fcm.key` optional parameter removed

This parameter was previously deprecated and has been removed.  Use the `fcm.google_service_account_credentials` parameter going forward.

### `hms` schema type change

This parameter `hms` was previously a block data type, and is now a single nested object type.

Previous configuration example:

```terraform
resource "pingone_mfa_application_push_credential" "example_hms" {
  # ... other configuration parameters

  hms {
    client_id     = var.hms_client_id
    client_secret = var.hms_client_secret
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_application_push_credential" "example_hms" {
  # ... other configuration parameters

  hms = {
    client_id     = var.hms_client_id
    client_secret = var.hms_client_secret
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

## Resource: pingone_mfa_settings

### `authentication` optional parameter removed

This parameter was previously deprecated and has been removed.  Device authentication parameters have moved to the `pingone_mfa_device_policy` resource.

## Resource: pingone_resource_attribute

### `resource_id` parameter changed

This parameter was previously deprecated and has now been made read only.  Use the `resource_name` parameter going forward.

### `resource_name` parameter changed

This parameter was previously optional and has now been made a required field.

## Resource: pingone_schema_attribute

### `schema_id` parameter changed

This parameter was previously deprecated and has now been made read only.  Use the optional `schema_name` parameter going forward.

## Resource: pingone_user

### `status` optional parameter removed

This parameter was previously deprecated and has been removed.  Use the `enabled` parameter going forward.

## Data Source: pingone_organization

### `base_url_agreement_management` computed attribute removed

This parameter was previously deprecated and has been removed.  Consider using the [PingOne Utilities module](https://registry.terraform.io/modules/pingidentity/utils/pingone/latest) going forward.

### `base_url_api` computed attribute removed

This parameter was previously deprecated and has been removed.  Consider using the [PingOne Utilities module](https://registry.terraform.io/modules/pingidentity/utils/pingone/latest) going forward.

### `base_url_apps` computed attribute removed

This parameter was previously deprecated and has been removed.  Consider using the [PingOne Utilities module](https://registry.terraform.io/modules/pingidentity/utils/pingone/latest) going forward.

### `base_url_auth` computed attribute removed

This parameter was previously deprecated and has been removed.  Consider using the [PingOne Utilities module](https://registry.terraform.io/modules/pingidentity/utils/pingone/latest) going forward.

### `base_url_console` computed attribute removed

This parameter was previously deprecated and has been removed.  Consider using the [PingOne Utilities module](https://registry.terraform.io/modules/pingidentity/utils/pingone/latest) going forward.

### `base_url_orchestrate` computed attribute removed

This parameter was previously deprecated and has been removed.  Consider using the [PingOne Utilities module](https://registry.terraform.io/modules/pingidentity/utils/pingone/latest) going forward.

## Data Source: pingone_user

### `status` computed attribute removed

This attribute was previously deprecated and has been removed.  Use the `enabled` attribute going forward.
