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

### `saml_options.sp_verification_certificate_ids` optional parameter removed

This parameter was previously deprecated and has been removed.  Use the `saml_options.sp_verification.certificate_ids` parameter going forward.

Previous configuration example:

```terraform
resource "pingone_application" "my_awesome_saml_app" {
  # ... other configuration parameters

  saml_options {
    # ... other configuration parameters

    sp_verification_certificate_ids = [pingone_certificate.my_awesome_certificate.id]
  }
}
```

New configuration example:

```terraform
resource "pingone_application" "my_awesome_saml_app" {
  # ... other configuration parameters

  saml_options = {
    # ... other configuration parameters

    sp_verification = {
      certificate_ids = [pingone_certificate.my_awesome_certificate.id]
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

### `service` block parameter renamed, data type changed and made a required parameter

The `service` parameter has been renamed to `services` and is now a required parameter.  The data type is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_environment" "my_environment" {
  # ... other configuration parameters
  
  service {
    type = "SSO"
  }

  service {
    type = "MFA"
  }

  service {
    type        = "PingFederate"
    console_url = "https://my-pingfederate-console.example.com/pingfederate"
  }
}
```

New configuration example:

```terraform
resource "pingone_environment" "my_environment" {
  # ... other configuration parameters
  
  services = [
    {
      type = "SSO"
    },
    {
      type = "MFA"
    },
    {
      type        = "PingFederate"
      console_url = "https://my-pingfederate-console.example.com/pingfederate"
    }
  ]
}
```

### `service.bookmark` block parameter renamed and data type changed

The `service.bookmark` parameter has been renamed to `services.bookmarks`.  The data type is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_environment" "my_environment" {
  # ... other configuration parameters
  
  service {
    type = "SSO"
    
    bookmark {
      name = "My awesome bookmark"
      url = "https://www.bxretail.org"
    }

    bookmark {
      name = "My second awesome bookmark"
      url = "https://www.bxretail.org/awesome"
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_environment" "my_environment" {
  # ... other configuration parameters
  
  services = [
    {
      type = "SSO"
    
      bookmarks = [
        {
          name = "My awesome bookmark"
          url = "https://www.bxretail.org"
        },
        {
          name = "My second awesome bookmark"
          url = "https://www.bxretail.org/awesome"
        }
      ]
    }
  ]
}
```

### `service.type` now made a required parameter

The `service.type` parameter has moved to `services.type` and is now a required parameter.

Previous configuration example:

```terraform
resource "pingone_environment" "my_environment" {
  # ... other configuration parameters
  
  service {}
}
```

New configuration example:

```terraform
resource "pingone_environment" "my_environment" {
  # ... other configuration parameters
  
  services = [
    {
      type = "SSO"
    }
  ]
}
```

### `timeouts` block removed

This parameter block is no longer needed and has been removed.

## Resource: pingone_identity_provider

### `amazon` parameter data type change

The `amazon` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  amazon {
    client_id     = var.amazon_client_id
    client_secret = var.amazon_client_secret
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  amazon = {
    client_id     = var.amazon_client_id
    client_secret = var.amazon_client_secret
  }
}
```

### `apple` parameter data type change

The `apple` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  apple {
    client_id                 = var.apple_client_id
    client_secret_signing_key = var.apple_client_secret_signing_key
    key_id                    = var.apple_key_id
    team_id                   = var.apple_team_id
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  apple = {
    client_id                 = var.apple_client_id
    client_secret_signing_key = var.apple_client_secret_signing_key
    key_id                    = var.apple_key_id
    team_id                   = var.apple_team_id
  }
}
```

### `facebook` parameter data type change

The `facebook` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  facebook {
    app_id     = var.facebook_app_id
    app_secret = var.facebook_app_secret
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  facebook = {
    app_id     = var.facebook_app_id
    app_secret = var.facebook_app_secret
  }
}
```

### `github` parameter data type change

The `github` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  github {
    client_id     = var.github_client_id
    client_secret = var.github_client_secret
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  github = {
    client_id     = var.github_client_id
    client_secret = var.github_client_secret
  }
}
```

### `google` parameter data type change

The `google` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  google {
    client_id     = var.google_client_id
    client_secret = var.google_client_secret
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  google = {
    client_id     = var.google_client_id
    client_secret = var.google_client_secret
  }
}
```

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

### `linkedin` parameter data type change

The `linkedin` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  linkedin {
    client_id     = var.linkedin_client_id
    client_secret = var.linkedin_client_secret
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  linkedin = {
    client_id     = var.linkedin_client_id
    client_secret = var.linkedin_client_secret
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

### `microsoft` parameter data type change

The `microsoft` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  microsoft {
    client_id     = var.microsoft_client_id
    client_secret = var.microsoft_client_secret
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  microsoft = {
    client_id     = var.microsoft_client_id
    client_secret = var.microsoft_client_secret
  }
}
```

### `openid_connect` parameter data type change

The `openid_connect` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  openid_connect {
    authorization_endpoint = var.openid_connect_idp_authorization_endpoint
    client_id              = var.openid_connect_idp_client_id
    client_secret          = var.openid_connect_idp_client_secret
    issuer                 = var.openid_connect_idp_issuer
    jwks_endpoint          = var.openid_connect_idp_jwks_endpoint
    scopes                 = var.openid_connect_idp_scopes
    token_endpoint         = var.openid_connect_idp_token_endpoint
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  openid_connect = {
    authorization_endpoint = var.openid_connect_idp_authorization_endpoint
    client_id              = var.openid_connect_idp_client_id
    client_secret          = var.openid_connect_idp_client_secret
    issuer                 = var.openid_connect_idp_issuer
    jwks_endpoint          = var.openid_connect_idp_jwks_endpoint
    scopes                 = var.openid_connect_idp_scopes
    token_endpoint         = var.openid_connect_idp_token_endpoint
  }
}
```

### `paypal` parameter data type change

The `paypal` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  paypal {
    client_environment     = var.paypal_client_environment
    client_id              = var.paypal_client_id
    client_secret          = var.paypal_client_secret
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  paypal = {
    client_environment     = var.paypal_client_environment
    client_id              = var.paypal_client_id
    client_secret          = var.paypal_client_secret
  }
}
```

### `saml` parameter data type change

The `saml` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  saml {
    idp_entity_id                    = var.saml_idp_entity_id
    idp_verification_certificate_ids = var.saml_idp_verification_certificate_ids 
    sp_entity_id                     = var.saml_idp_sp_entity_id
    sso_binding                      = var.saml_idp_sso_binding
    sso_endpoint                     = var.saml_idp_sso_endpoint
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  saml = {
    idp_entity_id                    = var.saml_idp_entity_id
    idp_verification_certificate_ids = var.saml_idp_verification_certificate_ids 
    sp_entity_id                     = var.saml_idp_sp_entity_id
    sso_binding                      = var.saml_idp_sso_binding
    sso_endpoint                     = var.saml_idp_sso_endpoint
  }
}
```

### `twitter` parameter data type change

The `twitter` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  twitter {
    client_id     = var.twitter_client_id
    client_secret = var.twitter_client_secret
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  twitter = {
    client_id     = var.twitter_client_id
    client_secret = var.twitter_client_secret
  }
}
```

### `yahoo` parameter data type change

The `yahoo` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  yahoo {
    client_id     = var.yahoo_client_id
    client_secret = var.yahoo_client_secret
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  yahoo = {
    client_id     = var.yahoo_client_id
    client_secret = var.yahoo_client_secret
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

## Resource: pingone_webhook

### `filter_options` parameter data type change

The `filter_options` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_webhook" "my_webhook" {
  # ... other configuration parameters
  
  filter_options {
    included_action_types = ["ACCOUNT.LINKED", "ACCOUNT.UNLINKED"]
  }
}
```

New configuration example:

```terraform
resource "pingone_webhook" "my_webhook" {
  # ... other configuration parameters
  
  filter_options = {
    included_action_types = ["ACCOUNT.LINKED", "ACCOUNT.UNLINKED"]
  }
}
```

## Data Source: pingone_application

### `access_control_group_options` computed attribute data type change

The `access_control_group_options` computed attribute is now a nested object type and no longer a list type.

### `external_link_options` computed attribute data type change

The `external_link_options` computed attribute is now a nested object type and no longer a list type.

### `icon` computed attribute data type change

The `icon` computed attribute is now a nested object type and no longer a list type.

### `oidc_options` computed attribute data type change

The `oidc_options` computed attribute is now a nested object type and no longer a list type.

### `oidc_options.certificate_based_authentication` computed attribute data type change

The `oidc_options.certificate_based_authentication` computed attribute is now a nested object type and no longer a list type.

### `oidc_options.cors_settings` computed attribute data type change

The `oidc_options.cors_settings` computed attribute is now a nested object type and no longer a list type.

### `oidc_options.mobile_app` computed attribute data type change

The `oidc_options.mobile_app` computed attribute is now a nested object type and no longer a list type.

### `oidc_options.mobile_app.integrity_detection` computed attribute data type change

The `oidc_options.mobile_app.integrity_detection` computed attribute is now a nested object type and no longer a list type.

### `oidc_options.mobile_app.integrity_detection.cache_duration` computed attribute data type change

The `oidc_options.mobile_app.integrity_detection.cache_duration` computed attribute is now a nested object type and no longer a list type.

### `oidc_options.mobile_app.integrity_detection.google_play` computed attribute data type change

The `oidc_options.mobile_app.integrity_detection.google_play` computed attribute is now a nested object type and no longer a list type.

### `saml_options` computed attribute data type change

The `saml_options` computed attribute is now a nested object type and no longer a list type.

### `saml_options.cors_settings` computed attribute data type change

The `saml_options.cors_settings` computed attribute is now a nested object type and no longer a list type.

### `saml_options.idp_signing_key` computed attribute data type change

The `saml_options.idp_signing_key` computed attribute is now a nested object type and no longer a list type.

### `saml_options.sp_verification` computed attribute data type change

The `saml_options.sp_verification` computed attribute is now a nested object type and no longer a list type.

### `saml_options.sp_verification_certificate_ids` computed attribute removed

This parameter was previously deprecated and has been removed.  Use the `saml_options.sp_verification.certificate_ids` attribute going forward.

## Data Source: pingone_environment

### `service` computed attribute rename and data type change

The `service` computed attribute has been renamed to `services` and is now a nested object type and no longer a block type.

### `service.bookmark` computed attribute rename and data type change

The `service.bookmark` computed attribute has been renamed to `services.bookmarks` and is now a nested object type and no longer a block type.

## Data Source: pingone_flow_policies

### `data_filter` optional parameter renamed and data type changed

This parameter has been renamed to `data_filters` and the data type changed.  The `data_filters` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
data "pingone_flow_policies" "example_by_data_filter" {
  # ... other configuration parameters
  
  data_filter {
    name   = "trigger.type"
    values = ["AUTHENTICATION"]
  }
}
```

New configuration example:

```terraform
data "pingone_flow_policies" "example_by_data_filter" {
  # ... other configuration parameters
  
  data_filters = [
    {
      name   = "trigger.type"
      values = ["AUTHENTICATION"]
    }
  ]
}
```

## Data Source: pingone_groups

### `data_filter` optional parameter renamed and data type changed

This parameter has been renamed to `data_filters` and the data type changed.  The `data_filters` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
data "pingone_groups" "example_by_data_filter" {
  # ... other configuration parameters
  
  data_filter {
    name   = "name"
    values = ["My first group", "My second group"]
  }
}
```

New configuration example:

```terraform
data "pingone_groups" "example_by_data_filter" {
  # ... other configuration parameters
  
  data_filters = [
    {
      name   = "name"
      values = ["My first group", "My second group"]
    }
  ]
}
```

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

## Data Source: pingone_populations

### `data_filter` optional parameter renamed and data type changed

This parameter has been renamed to `data_filters` and the data type changed.  The `data_filters` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
data "pingone_populations" "example_by_data_filter" {
  # ... other configuration parameters
  
  data_filter {
    name   = "name"
    values = ["My first population", "My second population"]
  }
}
```

New configuration example:

```terraform
data "pingone_populations" "example_by_data_filter" {
  # ... other configuration parameters
  
  data_filters = [
    {
      name   = "name"
      values = ["My first population", "My second population"]
    }
  ]
}
```

## Data Source: pingone_user

### `status` computed attribute removed

This attribute was previously deprecated and has been removed.  Use the `enabled` attribute going forward.

## Data Source: pingone_users

### `data_filter` optional parameter renamed and data type changed

This parameter has been renamed to `data_filters` and the data type changed.  The `data_filters` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
data "pingone_users" "example_by_data_filter" {
  # ... other configuration parameters
  
  data_filter {
    name = "memberOfGroups.id"
    values = [
      pingone_group.my_first_group.id,
      pingone_group.my_second_group.id
    ]
  }

  data_filter {
    name = "population.id"
    values = [
      pingone_population.my_population.id
    ]
  }
}
```

New configuration example:

```terraform
data "pingone_users" "example_by_data_filter" {
  # ... other configuration parameters
  
  data_filters = [
    {
      name = "memberOfGroups.id"
      values = [
        pingone_group.my_first_group.id,
        pingone_group.my_second_group.id
      ]
    },
    {
      name = "population.id"
      values = [
        pingone_population.my_population.id
      ]
    }
  ]
}
```