---
layout: ""
page_title: "PingOne Terraform Provider Version 1 Upgrade Guide"
description: |-
  Version 1.0.0 of the PingOne Terraform provider is a major release that introduces breaking changes to existing HCL.  This guide describes the changes that are required to upgrade v0.* PingOne Terraform provider releases to v1.0.0 onwards.
---

# PingOne Terraform Provider Version 1 Upgrade Guide

Version 1.0.0 of the PingOne Terraform provider is a major release that introduces breaking changes to existing HCL. This guide describes the changes that are required to upgrade v0.* PingOne Terraform provider releases to v1.0.0 onwards.

## Provider Configuration

### `force_delete_production_type` optional parameter removed

This parameter was previously deprecated and has been removed.

### `global_options.environment.production_type_force_delete` optional parameter removed

This parameter has been removed to mitigate the potential for accidental data loss.  In order to delete environments that are of type `PRODUCTION`, this must be done manually through the web console.  Where environments need to be removed in Terraform, ensure that they do not contain production data, and set their type as `SANDBOX`.

### `region` (with `PINGONE_REGION` environment variable) parameter removed

The `region` parameter (with the `PINGONE_REGION` environment variable) has been removed and replaced with `region_code` (with the `PINGONE_REGION_CODE` environment variable).  The following lists the mapping between the legacy `region` and new `region_code` values:

| Tenant Type                                  | Legacy `region` value | Replacement `region_code` value |
| -------------------------------------------- | --------------------- | ------------------------------- |
| Asia-Pacific with `.asia` top level domain   | `AsiaPacific`         | `AP`                            |
| Asia-Pacific with `.com.au` top level domain | N/a                   | `AU`                            |
| Canada with `.ca` top level domain           | `Canada`              | `CA`                            |
| Europe with `.eu` top level domain           | `Europe`              | `EU`                            |
| North America with `.com` top level domain   | `NorthAmerica`        | `NA`                            |

## Resource: pingone_application

### `access_control_group_options` parameter data type change

The `access_control_group_options` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_group" "my_awesome_group" {
  # ... other configuration parameters
}

resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  access_control_group_options {
    type = "ANY_GROUP"

    groups = [
      pingone_group.my_awesome_group.id
    ]
  }
}
```

New configuration example:

```terraform
resource "pingone_group" "my_awesome_group" {
  # ... other configuration parameters
}

resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  access_control_group_options = {
    type = "ANY_GROUP"

    groups = [
      pingone_group.my_awesome_group.id
    ]
  }
}
```

### `external_link_options` parameter data type change

The `external_link_options` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  external_link_options {
    # ... other configuration parameters
  }
}
```

New configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  external_link_options = {
    # ... other configuration parameters
  }
}
```

### `icon` parameter data type change

The `icon` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_image" "my_awesome_image" {
  # ... other configuration parameters
}

resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  icon {
    id   = pingone_image.my_awesome_image.id
    href = pingone_image.my_awesome_image.uploaded_image[0].href
  }
}
```

New configuration example:

```terraform
resource "pingone_image" "my_awesome_image" {
  # ... other configuration parameters
}

resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  icon = {
    id   = pingone_image.my_awesome_image.id
    href = pingone_image.my_awesome_image.uploaded_image.href
  }
}
```

### `oidc_options` parameter data type change

The `oidc_options` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options {
    # ... other configuration parameters
  }
}
```

New configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options = {
    # ... other configuration parameters
  }
}
```

### `oidc_options.bundle_id` optional parameter removed

This parameter was previously deprecated and has been removed.  Use the `oidc_options.mobile_app.bundle_id` parameter going forward.

### `oidc_options.certificate_based_authentication` parameter data type change

The `oidc_options.certificate_based_authentication` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_key" "my_awesome_key" {
  # ... other configuration parameters

  usage_type = "ISSUANCE"
}

resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options {
    # ... other configuration parameters

    certificate_based_authentication {
      key_id = pingone_key.my_awesome_key.id
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_key" "my_awesome_key" {
  # ... other configuration parameters

  usage_type = "ISSUANCE"
}

resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options = {
    # ... other configuration parameters

    certificate_based_authentication = {
      key_id = pingone_key.my_awesome_key.id
    }
  }
}
```

### `oidc_options.client_id` computed attribute removed

The `oidc_options.client_id` attribute has been removed from the `pingone_application` resource, as it is a duplicate of the application's own ID.

Previous configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options {
    # ... other configuration parameters
  }
}

locals {
  my_awesome_application_client_id = pingone_application.my_awesome_application.oidc_options[0].client_id
}
```

New configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options = {
    # ... other configuration parameters
  }
}

locals {
  my_awesome_application_client_id = pingone_application.my_awesome_application.id
}
```

### `oidc_options.client_secret` computed attribute removed

The `oidc_options.client_secret` attribute has been removed from the `pingone_application` resource, and is now found in the `pingone_application_secret` resource and/or data source.  Using the `pingone_application_secret` resource and data source has the benefit of being able to track the state of, and manage, previous secrets when performing application secret rotation.

Previous configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options {
    # ... other configuration parameters
  }
}

locals {
  my_awesome_application_client_id     = pingone_application.my_awesome_application.oidc_options[0].client_id
  my_awesome_application_client_secret = pingone_application.my_awesome_application.oidc_options[0].client_secret
}
```

New configuration example (using the `pingone_application_secret` resource):

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options = {
    # ... other configuration parameters
  }
}

resource "pingone_application_secret" "my_awesome_application" {
  # ... other configuration parameters

  application_id = pingone_application.my_awesome_application.id
}

locals {
  my_awesome_application_client_id     = pingone_application.my_awesome_application.id
  my_awesome_application_client_secret = pingone_application_secret.my_awesome_application.secret
}
```

New configuration example (using the `pingone_application_secret` data source):

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options = {
    # ... other configuration parameters
  }
}

data "pingone_application_secret" "my_awesome_application" {
  # ... other configuration parameters

  application_id = pingone_application.my_awesome_application.id
}

locals {
  my_awesome_application_client_id     = pingone_application.my_awesome_application.id
  my_awesome_application_client_secret = data.pingone_application_secret.my_awesome_application.secret
}
```

### `oidc_options.cors_settings` parameter data type change

The `oidc_options.cors_settings` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options {
    # ... other configuration parameters

    cors_settings {
      behavior = "ALLOW_SPECIFIC_ORIGINS"
      origins = [
        "http://localhost",
        "https://localhost",
        "http://auth.bxretail.org",
        "https://auth.bxretail.org",
        "http://*.bxretail.org",
        "https://*.bxretail.org",
        "http://192.168.1.1",
        "https://192.168.1.1",
      ]
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options = {
    # ... other configuration parameters

    cors_settings = {
      behavior = "ALLOW_SPECIFIC_ORIGINS"
      origins = [
        "http://localhost",
        "https://localhost",
        "http://auth.bxretail.org",
        "https://auth.bxretail.org",
        "http://*.bxretail.org",
        "https://*.bxretail.org",
        "http://192.168.1.1",
        "https://192.168.1.1",
      ]
    }
  }
}
```

### `oidc_options.mobile_app` parameter data type change

The `oidc_options.mobile_app` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options {
    # ... other configuration parameters

    mobile_app {
      bundle_id    = "org.bxretail.bundle"
      package_name = "org.bxretail.package"

      passcode_refresh_seconds = 45

      universal_app_link = "https://bxretail.org"
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options = {
    # ... other configuration parameters

    mobile_app = {
      bundle_id    = "org.bxretail.bundle"
      package_name = "org.bxretail.package"

      passcode_refresh_seconds = 45

      universal_app_link = "https://bxretail.org"
    }
  }
}
```

### `oidc_options.mobile_app.integrity_detection` parameter data type change

The `oidc_options.mobile_app.integrity_detection` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options {
    # ... other configuration parameters

    mobile_app {
      # ... other configuration parameters

      integrity_detection {
        enabled = true

        excluded_platforms = ["IOS"]

        cache_duration {
          amount = 30
          units  = "HOURS"
        }

        google_play {
          verification_type = "INTERNAL"
          decryption_key    = var.google_play_decryption_key
          verification_key  = var.google_play_verification_key
        }
      }
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options = {
    # ... other configuration parameters

    mobile_app = {
      # ... other configuration parameters

      integrity_detection = {
        enabled = true

        excluded_platforms = ["IOS"]

        cache_duration = {
          amount = 30
          units  = "HOURS"
        }

        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = var.google_play_decryption_key
          verification_key  = var.google_play_verification_key
        }
      }
    }
  }
}
```

### `oidc_options.mobile_app.integrity_detection.cache_duration` parameter data type change

The `oidc_options.mobile_app.integrity_detection.cache_duration` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options {
    # ... other configuration parameters

    mobile_app {
      # ... other configuration parameters

      integrity_detection {
        # ... other configuration parameters

        cache_duration {
          amount = 30
          units  = "HOURS"
        }
      }
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options = {
    # ... other configuration parameters

    mobile_app = {
      # ... other configuration parameters

      integrity_detection = {
        # ... other configuration parameters

        cache_duration = {
          amount = 30
          units  = "HOURS"
        }
      }
    }
  }
}
```

### `oidc_options.mobile_app.integrity_detection.google_play` parameter data type change

The `oidc_options.mobile_app.integrity_detection.google_play` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options {
    # ... other configuration parameters

    mobile_app {
      # ... other configuration parameters

      integrity_detection {
        # ... other configuration parameters

        google_play {
          verification_type = "INTERNAL"
          decryption_key    = var.google_play_decryption_key
          verification_key  = var.google_play_verification_key
        }
      }
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options = {
    # ... other configuration parameters

    mobile_app = {
      # ... other configuration parameters

      integrity_detection = {
        # ... other configuration parameters

        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = var.google_play_decryption_key
          verification_key  = var.google_play_verification_key
        }
      }
    }
  }
}
```
### `oidc_options.package_name` optional parameter removed

This parameter was previously deprecated and has been removed.  Use the `oidc_options.mobile_app.package_name` parameter going forward.

### `saml_options` parameter data type change

The `saml_options` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_application" "my_awesome_app" {
  # ... other configuration parameters

  saml_options {
    acs_urls           = ["https:/bxretail.org"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:bxretail"
  }
}
```

New configuration example:

```terraform
resource "pingone_key" "my_awesome_key" {
  # ... other configuration parameters

  usage_type = "SIGNING"
}

resource "pingone_application" "my_awesome_app" {
  # ... other configuration parameters

  saml_options = {
    acs_urls           = ["https:/bxretail.org"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:bxretail"

    idp_signing_key = {
      key_id    = pingone_key.my_awesome_key.id
      algorithm = pingone_key.my_awesome_key.signature_algorithm
    }
  }
}
```

### `saml_options.cors_settings` parameter data type change

The `saml_options.cors_settings` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  saml_options {
    # ... other configuration parameters

    cors_settings {
      behavior = "ALLOW_SPECIFIC_ORIGINS"
      origins = [
        "http://localhost",
        "https://localhost",
        "http://auth.bxretail.org",
        "https://auth.bxretail.org",
        "http://*.bxretail.org",
        "https://*.bxretail.org",
        "http://192.168.1.1",
        "https://192.168.1.1",
      ]
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  saml_options = {
    # ... other configuration parameters

    cors_settings = {
      behavior = "ALLOW_SPECIFIC_ORIGINS"
      origins = [
        "http://localhost",
        "https://localhost",
        "http://auth.bxretail.org",
        "https://auth.bxretail.org",
        "http://*.bxretail.org",
        "https://*.bxretail.org",
        "http://192.168.1.1",
        "https://192.168.1.1",
      ]
    }
  }
}
```

### `saml_options.idp_signing_key` parameter changed

This parameter was previously optional and has now been made a required parameter.

Previous configuration example:

```terraform
resource "pingone_application" "my_awesome_app" {
  # ... other configuration parameters

  saml_options {
    acs_urls           = ["https:/bxretail.org"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:bxretail"
  }
}
```

New configuration example:

```terraform
resource "pingone_key" "my_awesome_key" {
  # ... other configuration parameters

  usage_type = "SIGNING"
}

resource "pingone_application" "my_awesome_app" {
  # ... other configuration parameters

  saml_options = {
    acs_urls           = ["https:/bxretail.org"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:bxretail"

    idp_signing_key = {
      key_id    = pingone_key.my_awesome_key.id
      algorithm = pingone_key.my_awesome_key.signature_algorithm
    }
  }
}
```

### `saml_options.idp_signing_key` parameter data type change

The `saml_options.idp_signing_key` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_key" "my_awesome_key" {
  # ... other configuration parameters

  usage_type = "SIGNING"
}

resource "pingone_application" "my_awesome_app" {
  # ... other configuration parameters

  saml_options {
    acs_urls           = ["https:/bxretail.org"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:bxretail"

    idp_signing_key {
      key_id    = pingone_key.my_awesome_key.id
      algorithm = pingone_key.my_awesome_key.signature_algorithm
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_key" "my_awesome_key" {
  # ... other configuration parameters

  usage_type = "SIGNING"
}

resource "pingone_application" "my_awesome_app" {
  # ... other configuration parameters

  saml_options = {
    acs_urls           = ["https:/bxretail.org"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:bxretail"

    idp_signing_key = {
      key_id    = pingone_key.my_awesome_key.id
      algorithm = pingone_key.my_awesome_key.signature_algorithm
    }
  }
}
```

### `saml_options.sp_verification` parameter data type change

The `saml_options.sp_verification` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_certificate" "my_awesome_certificate" {
  # ... other configuration parameters

  usage_type = "SIGNING"
}

resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  saml_options {
    # ... other configuration parameters

    sp_verification {
      authn_request_signed = true
      certificate_ids = [
        pingone_certificate.my_awesome_certificate.id,
      ]
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_certificate" "my_awesome_certificate" {
  # ... other configuration parameters

  usage_type = "SIGNING"
}

resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  saml_options = {
    # ... other configuration parameters

    sp_verification = {
      authn_request_signed = true
      certificate_ids = [
        pingone_certificate.my_awesome_certificate.id,
      ]
    }
  }
}
```

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

## Resource: pingone_custom_domain_verify

### `timeouts` optional parameter data type changed

The `timeouts` parameter data type is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
resource "pingone_custom_domain_verify" "my_awesome_domain" {
  # ... other configuration parameters

  timeouts {
    create = "10m"
  }
}
```

New configuration example:

```terraform
resource "pingone_custom_domain_verify" "my_awesome_domain" {
  # ... other configuration parameters

  timeouts = {
    create = "10m"
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

## Resource: pingone_gateway

### `radius_client` parameter rename and data type change

The `radius_client` parameter has been renamed to `radius_clients` and is now a set of objects data type and no longer a block set type.

Previous configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  type = "RADIUS"

  # ... other configuration parameters

  radius_client {
    ip = "127.0.0.1"
  }

  radius_client {
    ip = "127.0.0.2"
  }
}
```

New configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  type = "RADIUS"

  # ... other configuration parameters

  radius_clients = [
    {
      ip = "127.0.0.1"
    },
    {
      ip = "127.0.0.2"
    }
  ]
}
```

### `user_type` parameter rename and data type change

The `user_type` parameter has been renamed to `user_types` and is now a map of objects data type and no longer a block set of objects type.  The map key of the new data type is the name of the user type (previously the `user_type.name` parameter).

Previous configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  type = "LDAP"

  # ... other configuration parameters

  user_type {
    name               = "User Type 1"
    password_authority = "PING_ONE"

    # ... other configuration parameters
  }

  user_type {
    name               = "User Type 2"
    password_authority = "LDAP"

    # ... other configuration parameters
  }
}
```

New configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  type = "LDAP"

  # ... other configuration parameters

  user_types = {
    "User Type 1" = {
      password_authority = "PING_ONE"

      # ... other configuration parameters
    },

    "User Type 2" = {
      password_authority = "LDAP"
      
      # ... other configuration parameters
    }
  }
}
```

### `user_type.push_password_changes_to_ldap` parameter rename

The `user_type.push_password_changes_to_ldap` parameter has been renamed to `user_types.allow_password_changes`.

Previous configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  type = "LDAP"

  # ... other configuration parameters

  user_type {
    name               = "User Type 1"
    password_authority = "LDAP"

    push_password_changes_to_ldap = true

    # ... other configuration parameters
  }
}
```

New configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  type = "LDAP"

  # ... other configuration parameters

  user_types = {
    "User Type 1" = {
      password_authority = "LDAP"
      
      allow_password_changes = true

      # ... other configuration parameters
    }
  }
}
```

### `user_type.user_migration` parameter rename and data type change

The `user_type.user_migration` parameter has been renamed to `user_types.new_user_lookup` and is now a single object data type and no longer a block set type.

Previous configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  type = "LDAP"

  # ... other configuration parameters

  user_type {
    name               = "User Type 1"
    password_authority = "LDAP"

    user_migration {
      lookup_filter_pattern = "(|(sAMAccountName=$${identifier})(UserPrincipalName=$${identifier}))"

      # ... other configuration parameters
    }

    # ... other configuration parameters
  }
}
```

New configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  type = "LDAP"

  # ... other configuration parameters

  user_types = {
    "User Type 1" = {
      password_authority = "LDAP"
      
      new_user_lookup = {
        ldap_filter_pattern = "(|(sAMAccountName=$${identifier})(UserPrincipalName=$${identifier}))"

        # ... other configuration parameters
      }

      # ... other configuration parameters
    }
  }
}
```

### `user_type.user_migration.attribute_mapping` parameter rename and data type change

The `user_type.user_migration.attribute_mapping` parameter has been renamed to `user_types.new_user_lookup.attribute_mappings` and is now a set of objects data type and no longer a block set type.

Previous configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  type = "LDAP"

  # ... other configuration parameters

  user_type {
    name               = "User Type 1"
    password_authority = "LDAP"

    user_migration {
      lookup_filter_pattern = "(|(sAMAccountName=$${identifier})(UserPrincipalName=$${identifier}))"

      attribute_mapping {
        name  = "username"
        value = "$${ldapAttributes.sAMAccountName}"
      }

      attribute_mapping {
        name  = "email"
        value = "$${ldapAttributes.mail}"
      }

      # ... other configuration parameters
    }

    # ... other configuration parameters
  }
}
```

New configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  type = "LDAP"

  # ... other configuration parameters

  user_types = {
    "User Type 1" = {
      password_authority = "LDAP"
      
      new_user_lookup = {
        ldap_filter_pattern = "(|(sAMAccountName=$${identifier})(UserPrincipalName=$${identifier}))"

        attribute_mappings = [
          {
            name  = "username"
            value = "$${ldapAttributes.sAMAccountName}"
          },
          {
            name  = "email"
            value = "$${ldapAttributes.mail}"
          }
        ]

        # ... other configuration parameters
      }

      # ... other configuration parameters
    }
  }
}
```

### `user_type.user_migration.lookup_filter_pattern` parameter rename

The `user_type.user_migration.lookup_filter_pattern` parameter has been renamed to `user_types.new_user_lookup.ldap_filter_pattern`.

Previous configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  type = "LDAP"

  # ... other configuration parameters

  user_type {
    name               = "User Type 1"
    password_authority = "LDAP"

    user_migration {
      lookup_filter_pattern = "(|(sAMAccountName=$${identifier})(UserPrincipalName=$${identifier}))"

      # ... other configuration parameters
    }

    # ... other configuration parameters
  }
}
```

New configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  type = "LDAP"

  # ... other configuration parameters

  user_types = {
    "User Type 1" = {
      password_authority = "LDAP"
      
      new_user_lookup = {
        ldap_filter_pattern = "(|(sAMAccountName=$${identifier})(UserPrincipalName=$${identifier}))"

        # ... other configuration parameters
      }

      # ... other configuration parameters
    }
  }
}
```

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

## Resource: pingone_mfa_policy (now pingone_mfa_device_policy)

Review the [Upgrade MFA Policies to use FIDO2 with Passkeys](./upgrade-mfa-policy-for-fido2) to ensure all MFA Policies are upgraded in the PingOne tenant prior to upgrading the PingOne provider version to `v1.0.0`.

### Resource renamed to `pingone_mfa_device_policy`

The `pingone_mfa_policy` resource has been renamed to `pingone_mfa_device_policy` to better align with the console and API experience.

### `device_selection` moved

This parameter `device_selection` has been moved to `authentication.device_selection`.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  device_selection = "DEFAULT_TO_FIRST"
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  authentication = {
    device_selection = "DEFAULT_TO_FIRST"
  }
}
```

### `email` schema type change

This parameter `email` was previously a block data type, and is now a single nested object type.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  email {
    # ... other configuration parameters
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  email = {
    # ... other configuration parameters
  }
}
```

### `email.otp_failure_cooldown_duration`, `email.otp_failure_cooldown_timeunit` and `email.otp_failure_cooldown_count` moved

The parameters `email.otp_failure_cooldown_duration`, `email.otp_failure_cooldown_timeunit` and `email.otp_failure_cooldown_count` have been moved to a new object, `email.otp.failure.*`

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  email {
    # ... other configuration parameters

    otp_failure_cooldown_duration = 2
    otp_failure_cooldown_timeunit = "MINUTES"
    otp_failure_count             = 3

  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  email = {
    # ... other configuration parameters

    otp = {
      failure = {
        count = 3

        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
  }
}
```

### `email.otp_lifetime_duration` and `email.otp_lifetime_timeunit` moved

The parameters `email.otp_lifetime_duration` and `email.otp_lifetime_timeunit` have been moved to a new object, `email.otp.lifetime.*`

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  email {
    # ... other configuration parameters

    otp_lifetime_duration = 2
    otp_lifetime_timeunit = "MINUTES"

  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  email = {
    # ... other configuration parameters

    otp = {
      lifetime = {
        duration  = 2
        time_unit = "MINUTES"
      }
    }
  }
}
```

### `fido2` schema type change

This parameter `fido2` was previously a block data type, and is now a single nested object type.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  fido2 {
    # ... other configuration parameters
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  fido2 = {
    # ... other configuration parameters
  }
}
```

### `platform` optional parameter removed

This parameter was previously deprecated and has been removed.  Use the `fido2` parameter going forward.

### `security_key` optional parameter removed

This parameter was previously deprecated and has been removed.  Use the `fido2` parameter going forward.

### `mobile` schema type change

This parameter `mobile` was previously a block data type, and is now a single nested object type.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile {
    # ... other configuration parameters
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile = {
    # ... other configuration parameters
  }
}
```

### `mobile.application` rename and schema type change

This parameter `mobile.application` has been renamed to `mobile.applications` and was previously a block set data type, and is now a map of objects data type.  The `id` parameter of the application is now the map key for the application object as shown in the example.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile {
    # ... other configuration parameters

    application {
      # ... other configuration parameters

      id = pingone_application.my_mobile_application.id

      push_enabled = true
      otp_enabled  = true
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile = {
    # ... other configuration parameters

    applications = {
      (pingone_application.my_mobile_application.id) = {
        # ... other configuration parameters

        push = {
          enabled = true
        }

        otp = {
          enabled  = true
        }
      }
    }
  }
}
```

### `mobile.application.auto_enrollment_enabled` moved

This parameters `mobile.application.auto_enrollment_enabled` have been moved to `mobile.applications.*.auto_enrollment.enabled`.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile {
    # ... other configuration parameters

    application {
      # ... other configuration parameters

      auto_enrollment_enabled = true
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile = {
    # ... other configuration parameters

    applications = {
      (pingone_application.my_mobile_application.id) = {
        # ... other configuration parameters

        auto_enrollment = {
          enabled = true
        }
      }
    }
  }
}
```

### `mobile.application.device_authorization_enabled` moved

This parameter `mobile.application.device_authorization_enabled` has been moved to `mobile.applications.*.device_authorization.enabled`.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile {
    # ... other configuration parameters

    application {
      # ... other configuration parameters

      device_authorization_enabled = true
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile = {
    # ... other configuration parameters

    applications = {
      (pingone_application.my_mobile_application.id) = {
        # ... other configuration parameters

        device_authorization = {
          enabled = true
        }
      }
    }
  }
}
```

### `mobile.application.device_authorization_extra_verification` moved

This parameter `mobile.application.device_authorization_extra_verification` has been moved to `mobile.applications.*.device_authorization.extra_verification`.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile {
    # ... other configuration parameters

    application {
      # ... other configuration parameters

      device_authorization_extra_verification = "permissive"
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile = {
    # ... other configuration parameters

    applications = {
      (pingone_application.my_mobile_application.id) = {
        # ... other configuration parameters

        device_authorization = {
          enabled            = true
          extra_verification = "permissive"
        }
      }
    }
  }
}
```

### `mobile.application.otp_enabled` moved

This parameter `mobile.application.otp_enabled` has been moved to `mobile.applications.*.otp.enabled`.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile {
    # ... other configuration parameters

    application {
      # ... other configuration parameters

      id = pingone_application.my_mobile_application.id

      push_enabled = true
      otp_enabled  = true
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile = {
    # ... other configuration parameters

    applications = {
      (pingone_application.my_mobile_application.id) = {
        # ... other configuration parameters

        push = {
          enabled = true
        }

        otp = {
          enabled  = true
        }
      }
    }
  }
}
```

### `mobile.application.pairing_key_lifetime_duration` and `mobile.application.pairing_key_lifetime_timeunit` moved

This parameters `mobile.application.pairing_key_lifetime_duration` and `mobile.application.pairing_key_lifetime_timeunit` have been moved to `mobile.applications.*.pairing_key_lifetime`.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile {
    # ... other configuration parameters

    application {
      # ... other configuration parameters

      pairing_key_lifetime_duration = "10"
      pairing_key_lifetime_timeunit = "MINUTES"
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile = {
    # ... other configuration parameters

    applications = {
      (pingone_application.my_mobile_application.id) = {
        # ... other configuration parameters

        pairing_key_lifetime = {
          duration = "10"
          time_unit = "MINUTES"
        }
      }
    }
  }
}
```

### `mobile.application.push_enabled` moved

This parameter `mobile.application.push_enabled` has been moved to `mobile.applications.*.push.enabled`.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile {
    # ... other configuration parameters

    application {
      # ... other configuration parameters

      id = pingone_application.my_mobile_application.id

      push_enabled = true
      otp_enabled  = true
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile = {
    # ... other configuration parameters

    applications = {
      (pingone_application.my_mobile_application.id) = {
        # ... other configuration parameters

        push = {
          enabled = true
        }

        otp = {
          enabled  = true
        }
      }
    }
  }
}
```

### `mobile.application.push_limit_count`, `mobile.application.push_limit_lock_duration`, `mobile.application.push_limit_lock_duration_timeunit`, `mobile.application.push_limit_time_period_duration` and `mobile.application.push_limit_time_period_timeunit` moved

This parameters `mobile.application.push_limit_count`, `mobile.application.push_limit_lock_duration`, `mobile.application.push_limit_lock_duration_timeunit`, `mobile.application.push_limit_time_period_duration` and `mobile.application.push_limit_time_period_timeunit` have been moved to `mobile.applications.*.push_limit`.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile {
    # ... other configuration parameters

    application {
      # ... other configuration parameters

      push_limit_count                  = 5
      push_limit_lock_duration          = 30
      push_limit_lock_duration_timeunit = "MINUTES"
      push_limit_time_period_duration   = 10
      push_limit_time_period_timeunit   = "MINUTES"
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile = {
    # ... other configuration parameters

    applications = {
      (pingone_application.my_mobile_application.id) = {
        # ... other configuration parameters

        push_limit = {
          count = 5

          lock_duration = {
            duration  = 30
            time_unit = "MINUTES"
          }

          time_period = {
            duration  = 10
            time_unit = "MINUTES"
          }
        }
      }
    }
  }
}
```

### `mobile.application.push_timeout_duration` and `mobile.application.push_timeout_timeunit` moved

This parameters `mobile.application.push_timeout_duration` and `mobile.application.push_timeout_timeunit` have been moved to `mobile.applications.*.push_timeout`.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile {
    # ... other configuration parameters

    application {
      # ... other configuration parameters

      push_timeout_duration = 40
      push_timeout_timeunit = "SECONDS"
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile = {
    # ... other configuration parameters

    applications = {
      (pingone_application.my_mobile_application.id) = {
        # ... other configuration parameters

        push_timeout = {
          duration  = 40
          time_unit = "SECONDS"
        }
      }
    }
  }
}
```

### `mobile.otp_failure_cooldown_duration`, `mobile.otp_failure_cooldown_timeunit` and `mobile.otp_failure_cooldown_count` moved

The parameters `mobile.otp_failure_cooldown_duration`, `mobile.otp_failure_cooldown_timeunit` and `mobile.otp_failure_cooldown_count` have been moved to a new object, `mobile.otp.failure.*`

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile {
    # ... other configuration parameters

    otp_failure_cooldown_duration = 2
    otp_failure_cooldown_timeunit = "MINUTES"
    otp_failure_count             = 3

  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  mobile = {
    # ... other configuration parameters

    otp = {
      failure = {
        count = 3

        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
  }
}
```

### `sms` schema type change

This parameter `sms` was previously a block data type, and is now a single nested object type.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  sms {
    # ... other configuration parameters
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  sms = {
    # ... other configuration parameters
  }
}
```

### `sms.otp_failure_cooldown_duration`, `sms.otp_failure_cooldown_timeunit` and `sms.otp_failure_cooldown_count` moved

The parameters `sms.otp_failure_cooldown_duration`, `sms.otp_failure_cooldown_timeunit` and `sms.otp_failure_cooldown_count` have been moved to a new object, `sms.otp.failure.*`

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  sms {
    # ... other configuration parameters

    otp_failure_cooldown_duration = 2
    otp_failure_cooldown_timeunit = "MINUTES"
    otp_failure_count             = 3

  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  sms = {
    # ... other configuration parameters

    otp = {
      failure = {
        count = 3

        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
  }
}
```

### `sms.otp_lifetime_duration` and `sms.otp_lifetime_timeunit` moved

The parameters `sms.otp_lifetime_duration` and `sms.otp_lifetime_timeunit` have been moved to a new object, `sms.otp.lifetime.*`

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  sms {
    # ... other configuration parameters

    otp_lifetime_duration = 2
    otp_lifetime_timeunit = "MINUTES"

  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  sms = {
    # ... other configuration parameters

    otp = {
      lifetime = {
        duration  = 2
        time_unit = "MINUTES"
      }
    }
  }
}
```

### `totp` schema type change

This parameter `totp` was previously a block data type, and is now a single nested object type.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  totp {
    # ... other configuration parameters
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  totp = {
    # ... other configuration parameters
  }
}
```

### `totp.otp_failure_cooldown_duration`, `totp.otp_failure_cooldown_timeunit` and `totp.otp_failure_cooldown_count` moved

The parameters `totp.otp_failure_cooldown_duration`, `totp.otp_failure_cooldown_timeunit` and `totp.otp_failure_cooldown_count` have been moved to a new object, `totp.otp.failure.*`

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  totp {
    # ... other configuration parameters

    otp_failure_cooldown_duration = 2
    otp_failure_cooldown_timeunit = "MINUTES"
    otp_failure_count             = 3

  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  totp = {
    # ... other configuration parameters

    otp = {
      failure = {
        count = 3

        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
  }
}
```

### `voice` schema type change

This parameter `voice` was previously a block data type, and is now a single nested object type.

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  voice {
    # ... other configuration parameters
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  voice = {
    # ... other configuration parameters
  }
}
```

### `voice.otp_failure_cooldown_duration`, `voice.otp_failure_cooldown_timeunit` and `voice.otp_failure_cooldown_count` moved

The parameters `voice.otp_failure_cooldown_duration`, `voice.otp_failure_cooldown_timeunit` and `voice.otp_failure_cooldown_count` have been moved to a new object, `voice.otp.failure.*`

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  voice {
    # ... other configuration parameters

    otp_failure_cooldown_duration = 2
    otp_failure_cooldown_timeunit = "MINUTES"
    otp_failure_count             = 3

  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  voice = {
    # ... other configuration parameters

    otp = {
      failure = {
        count = 3

        cool_down = {
          duration  = 2
          time_unit = "MINUTES"
        }
      }
    }
  }
}
```

### `voice.otp_lifetime_duration` and `voice.otp_lifetime_timeunit` moved

The parameters `voice.otp_lifetime_duration` and `voice.otp_lifetime_timeunit` have been moved to a new object, `voice.otp.lifetime.*`

Previous configuration example:

```terraform
resource "pingone_mfa_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  voice {
    # ... other configuration parameters

    otp_lifetime_duration = 2
    otp_lifetime_timeunit = "MINUTES"

  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_device_policy" "my_awesome_mfa_device_policy" {
  # ... other configuration parameters

  voice = {
    # ... other configuration parameters

    otp = {
      lifetime = {
        duration  = 2
        time_unit = "MINUTES"
      }
    }
  }
}
```

## Resource: pingone_mfa_settings

### `authentication` optional parameter removed

This parameter was previously deprecated and has been removed.  Device authentication parameters have moved to the `pingone_mfa_device_policy` resource.

### `lockout` schema type change

This parameter `lockout` was previously a block data type, and is now a single nested object type.

Previous configuration example:

```terraform
resource "pingone_mfa_settings" "my_awesome_mfa_settings" {
  # ... other configuration parameters

  lockout {
    failure_count    = 5
    duration_seconds = 600
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_settings" "my_awesome_mfa_settings" {
  # ... other configuration parameters

  lockout = {
    failure_count    = 5
    duration_seconds = 600
  }
}
```

### `pairing` schema type change

This parameter `pairing` was previously a block data type, and is now a single nested object type.

Previous configuration example:

```terraform
resource "pingone_mfa_settings" "my_awesome_mfa_settings" {
  # ... other configuration parameters

  pairing {
    max_allowed_devices = 5
    pairing_key_format  = "ALPHANUMERIC"
  }
}
```

New configuration example:

```terraform
resource "pingone_mfa_settings" "my_awesome_mfa_settings" {
  # ... other configuration parameters

  pairing = {
    max_allowed_devices = 5
    pairing_key_format  = "ALPHANUMERIC"
  }
}
```

### `phone_extensions_enabled` parameter moved

This parameter `phone_extensions_enabled` has moved to a nested object type at `phone_extensions.enabled`.

Previous configuration example:

```terraform
resource "pingone_mfa_settings" "my_awesome_mfa_settings" {
  # ... other configuration parameters

  phone_extensions_enabled = true
}
```

New configuration example:

```terraform
resource "pingone_mfa_settings" "my_awesome_mfa_settings" {
  # ... other configuration parameters

  phone_extensions = {
    enabled = true
  }
}
```

## Resource: pingone_notification_policy

### `quota` schema type change

This parameter `quota` was previously a list block data type, and is now a set of nested objects type.

Previous configuration example:

```terraform
resource "pingone_notification_policy" "my_awesome_notification_policy" {
  # ... other configuration parameters

  quota {
    type             = "USER"
    delivery_methods = ["SMS", "Voice"]
    total            = 30
  }
  
  quota {
    type             = "USER"
    delivery_methods = ["Email"]
    total            = 30
  }
}
```

New configuration example:

```terraform
resource "pingone_notification_policy" "my_awesome_notification_policy" {
  # ... other configuration parameters

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

## Resource: pingone_notification_settings_email

### `from` schema type change

This parameter `from` was previously a block data type, and is now a single nested object type.

Previous configuration example:

```terraform
resource "pingone_notification_settings_email" "my_awesome_email_settings" {
  # ... other configuration parameters

  from {
    email_address = "noreply@bxretail.org"
  }
}
```

New configuration example:

```terraform
resource "pingone_notification_settings_email" "my_awesome_email_settings" {
  # ... other configuration parameters

  from = {
    email_address = "noreply@bxretail.org"
  }
}
```

### `reply_to` schema type change

This parameter `reply_to` was previously a block data type, and is now a single nested object type.

Previous configuration example:

```terraform
resource "pingone_notification_settings_email" "my_awesome_email_settings" {
  # ... other configuration parameters

  from {
    email_address = "noreply@bxretail.org"
  }

  reply_to {
    email_address = "customerservices@bxretail.org"
    name          = "BXRetail Customer Services"
  }
}
```

New configuration example:

```terraform
resource "pingone_notification_settings_email" "my_awesome_email_settings" {
  # ... other configuration parameters

  from = {
    email_address = "noreply@bxretail.org"
  }

  reply_to = {
    email_address = "customerservices@bxretail.org"
    name          = "BXRetail Customer Services"
  }
}
```

## Resource: pingone_password_policy

### `account_lockout` parameter rename and data type change

The `account_lockout` parameter has been renamed to `lockout` and is now a nested object type and no longer a block list type.

Previous configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  account_lockout {
    duration_seconds = 900
    fail_count       = 5
  }
}
```

New configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  lockout = {
    duration_seconds = 900
    failure_count    = 5
  }
}
```

### `account_lockout.fail_count` parameter renamed

The `account_lockout.fail_count` parameter has been renamed to `lockout.failure_count`.

Previous configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  account_lockout {
    # ... other configuration parameters

    fail_count = 5
  }
}
```

New configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  lockout = {
    # ... other configuration parameters

    failure_count = 5
  }
}
```

### `bypass_policy` parameter removed

The `bypass_policy` parameter has no effect and has been removed.

### `environment_default` parameter renamed

The `environment_default` parameter has been renamed to `default`.

Previous configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  environment_default = true
}
```

New configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  default = true
}
```

### `exclude_commonly_used_passwords` parameter renamed

The `exclude_commonly_used_passwords` parameter has been renamed to `excludes_commonly_used_passwords`.

Previous configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  exclude_commonly_used_passwords = true
}
```

New configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  excludes_commonly_used_passwords = true
}
```

### `exclude_profile_data` parameter renamed

The `exclude_profile_data` parameter has been renamed to `excludes_profile_data`.

Previous configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  exclude_profile_data = true
}
```

New configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  excludes_profile_data = true
}
```

### `password_age.max` parameter moved

The `password_age.max` parameter has been moved to `password_age_max`.

Previous configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  password_age {
    # ... other configuration parameters

    max = 30
  }
}
```

New configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  password_age_max = 30
}
```

### `password_age.min` parameter moved

The `password_age.min` parameter has been moved to `password_age_min`.

Previous configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  password_age {
    # ... other configuration parameters

    min = 1
  }
}
```

New configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  password_age_min = 1
}
```

### `password_history` parameter rename and data type change

The `password_history` parameter has been renamed to `history` and is now a nested object type and no longer a block list type.

Previous configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  password_history {
    prior_password_count = 6
    retention_days       = 365
  }
}
```

New configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  history {
    count          = 6
    retention_days = 365
  }
}
```

### `password_history.prior_password_count` parameter renamed

The `password_history.prior_password_count` parameter has been renamed to `history.count`.

Previous configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  password_history {
    # ... other configuration parameters

    prior_password_count = 6
  }
}
```

New configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  history {
    # ... other configuration parameters

    count = 6
  }
}
```

### `password_length` parameter rename and data type change

The `password_length` parameter has been renamed to `length` and is now a nested object type and no longer a block list type.

Previous configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  password_length {
    min = 8
    max = 255
  }
}
```

New configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  length = {
    min = 8
    max = 255
  }
}
```

## Resource: pingone_resource

### `client_secret` computed attribute removed

The `client_secret` attribute has been removed from the `pingone_resource` resource, and is now found in the `pingone_resource_secret` resource and/or data source.  Using the `pingone_resource_secret` resource and data source has the benefit of being able to track the state of, and manage, previous secrets when performing resource secret rotation.

Previous configuration example:

```terraform
resource "pingone_resource" "my_awesome_custom_resource" {
  # ... other configuration parameters
}

locals {
  my_awesome_resource_client_id     = pingone_resource.my_awesome_custom_resource.id
  my_awesome_resource_client_secret = pingone_resource.my_awesome_custom_resource.client_secret
}
```

New configuration example (using the `pingone_resource_secret` resource):

```terraform
resource "pingone_resource" "my_awesome_custom_resource" {
  # ... other configuration parameters
}

resource "pingone_resource_secret" "my_awesome_custom_resource" {
  # ... other configuration parameters

  resource_id = pingone_resource.my_awesome_custom_resource.id
}

locals {
  my_awesome_resource_client_id     = pingone_resource.my_awesome_custom_resource.id
  my_awesome_resource_client_secret = pingone_resource_secret.my_awesome_custom_resource.secret
}
```

New configuration example (using the `pingone_resource_secret` data source):

```terraform
resource "pingone_resource" "my_awesome_custom_resource" {
  # ... other configuration parameters
}

data "pingone_resource_secret" "my_awesome_custom_resource" {
  # ... other configuration parameters

  resource_id = pingone_resource.my_awesome_custom_resource.id
}

locals {
  my_awesome_resource_client_id     = pingone_resource.my_awesome_custom_resource.id
  my_awesome_resource_client_secret = data.pingone_resource_secret.my_awesome_custom_resource.secret
}
```

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

### `oidc_options.client_id` computed attribute removed

The `oidc_options.client_id` attribute has been removed from the `pingone_application` data source, as it is a duplicate of the application's own ID.

Previous configuration example:

```terraform
data "pingone_application" "my_awesome_application" {
  # ... other configuration parameters
}

locals {
  my_awesome_application_client_id = data.pingone_application.my_awesome_application.oidc_options[0].client_id
}
```

New configuration example:

```terraform
data "pingone_application" "my_awesome_application" {
  # ... other configuration parameters
}

locals {
  my_awesome_application_client_id = data.pingone_application.my_awesome_application.id
}
```

### `oidc_options.client_secret` computed attribute removed

The `oidc_options.client_secret` attribute has been removed from the `pingone_application` data source, and is now found in the `pingone_application_secret` resource and/or data source.  Using the `pingone_application_secret` resource and data source has the benefit of being able to track the state of, and manage, previous secrets when performing application secret rotation.

Previous configuration example:

```terraform
data "pingone_application" "my_awesome_application" {
  # ... other configuration parameters
}

locals {
  my_awesome_application_client_id     = data.pingone_application.my_awesome_application.oidc_options[0].client_id
  my_awesome_application_client_secret = data.pingone_application.my_awesome_application.oidc_options[0].client_secret
}
```

New configuration example (using the `pingone_application_secret` resource):

```terraform
data "pingone_application" "my_awesome_application" {
  # ... other configuration parameters
}

resource "pingone_application_secret" "my_awesome_application" {
  # ... other configuration parameters

  application_id = data.pingone_application.my_awesome_application.id
}

locals {
  my_awesome_application_client_id     = data.pingone_application.my_awesome_application.id
  my_awesome_application_client_secret = pingone_application_secret.my_awesome_application.secret
}
```

New configuration example (using the `pingone_application_secret` data source):

```terraform
data "pingone_application" "my_awesome_application" {
  # ... other configuration parameters
}

data "pingone_application_secret" "my_awesome_application" {
  # ... other configuration parameters

  application_id = data.pingone_application.my_awesome_application.id
}

locals {
  my_awesome_application_client_id     = data.pingone_application.my_awesome_application.id
  my_awesome_application_client_secret = data.pingone_application_secret.my_awesome_application.secret
}
```

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

## Data Source: pingone_flow_policy

### `davinci_application` computed attribute data type change

The `davinci_application` computed attribute is now a nested object type and no longer a list block type.

### `trigger` computed attribute data type change

The `trigger` computed attribute is now a nested object type and no longer a list block type.

## Data Source: pingone_gateway

### `radius_client` computed attribute rename and data type change

The `radius_client` computed attribute has been renamed to `radius_clients` and is now a set of objects data type and no longer a block set type.

### `user_type` computed attribute rename and data type change

The `user_type` computed attribute has been renamed to `user_types` and is now a map of objects data type and no longer a block set of objects type.  The map key of the new data type is the name of the user type (previously the `user_type.name` parameter).

### `user_type.push_password_changes_to_ldap` computed attribute rename

The `user_type.push_password_changes_to_ldap` computed attribute has been renamed to `user_types.allow_password_changes`.

### `user_type.user_migration` computed attribute rename and data type change

The `user_type.user_migration` computed attribute has been renamed to `user_types.new_user_lookup` and is now a single object data type and no longer a block set type.

### `user_type.user_migration.attribute_mapping` computed attribute rename and data type change

The `user_type.user_migration.attribute_mapping` computed attribute has been renamed to `user_types.new_user_lookup.attribute_mappings` and is now a set of objects data type and no longer a block set type.

### `user_type.user_migration.lookup_filter_pattern` computed attribute rename

The `user_type.user_migration.lookup_filter_pattern` computed attribute has been renamed to `user_types.new_user_lookup.ldap_filter_pattern`.

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

## Data Source: pingone_licenses

### `data_filter` optional parameter renamed and data type changed

This parameter has been renamed to `data_filters` and the data type changed.  The `data_filters` parameter is now a nested object type and no longer a block type.

Previous configuration example:

```terraform
data "pingone_licenses" "example_by_data_filter" {
  # ... other configuration parameters
  
  data_filter {
    name   = "status"
    values = ["ACTIVE"]
  }
}
```

New configuration example:

```terraform
data "pingone_licenses" "example_by_data_filter" {
  # ... other configuration parameters
  
  data_filters = [
    {
      name   = "status"
      values = ["ACTIVE"]
    }
  ]
}
```

## Data Source: pingone_mfa_policies (now pingone_mfa_device_policies)

### Data Source renamed to `pingone_mfa_device_policies`

The `pingone_mfa_policies` data source has been renamed to `pingone_mfa_device_policies` to better align with the console and API experience.

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

## Data Source: pingone_password_policy

### `account_lockout` computed attribute rename and data type change

The `account_lockout` computed attribute has been renamed to `lockout` and is now a nested object type and no longer a block list type.

### `account_lockout.fail_count` computed attribute renamed

The `account_lockout.fail_count` computed attribute has been renamed to `lockout.failure_count`.

### `bypass_policy` computed attribute removed

The `bypass_policy` computed attribute has no effect and has been removed.

### `environment_default` computed attribute renamed

The `environment_default` computed attribute has been renamed to `default`.

### `exclude_commonly_used_passwords` computed attribute renamed

The `exclude_commonly_used_passwords` computed attribute has been renamed to `excludes_commonly_used_passwords`.

### `exclude_profile_data` computed attribute renamed

The `exclude_profile_data` computed attribute has been renamed to `excludes_profile_data`.

### `password_age.max` computed attribute moved

The `password_age.max` computed attribute has been moved to `password_age_max`.

### `password_age.min` computed attribute moved

The `password_age.min` computed attribute has been moved to `password_age_min`.

### `password_history` computed attribute rename and data type change

The `password_history` computed attribute has been renamed to `history` and is now a nested object type and no longer a block list type.

### `password_history.prior_password_count` computed attribute renamed

The `password_history.prior_password_count` computed attribute has been renamed to `history.count`.

### `password_length` computed attribute rename and data type change

The `password_length` computed attribute has been renamed to `length` and is now a nested object type and no longer a block list type.

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

## Data Source: pingone_resource

### `client_secret` computed attribute removed

The `client_secret` attribute has been removed from the `pingone_resource` data source, and is now found in the `pingone_resource_secret` resource and/or data source.  Using the `pingone_resource_secret` resource and data source has the benefit of being able to track the state of, and manage, previous secrets when performing resource secret rotation.

Previous configuration example:

```terraform
data "pingone_resource" "my_awesome_custom_resource" {
  # ... other configuration parameters
}

locals {
  my_awesome_resource_client_id     = data.pingone_resource.my_awesome_custom_resource.id
  my_awesome_resource_client_secret = data.pingone_resource.my_awesome_custom_resource.client_secret
}
```

New configuration example (using the `pingone_resource_secret` resource):

```terraform
data "pingone_resource" "my_awesome_custom_resource" {
  # ... other configuration parameters
}

resource "pingone_resource_secret" "my_awesome_custom_resource" {
  # ... other configuration parameters

  resource_id = data.pingone_resource.my_awesome_custom_resource.id
}

locals {
  my_awesome_resource_client_id     = data.pingone_resource.my_awesome_custom_resource.id
  my_awesome_resource_client_secret = pingone_resource_secret.my_awesome_custom_resource.secret
}
```

New configuration example (using the `pingone_resource_secret` data source):

```terraform
data "pingone_resource" "my_awesome_custom_resource" {
  # ... other configuration parameters
}

data "pingone_resource_secret" "my_awesome_custom_resource" {
  # ... other configuration parameters

  resource_id = data.pingone_resource.my_awesome_custom_resource.id
}

locals {
  my_awesome_resource_client_id     = data.pingone_resource.my_awesome_custom_resource.id
  my_awesome_resource_client_secret = data.pingone_resource_secret.my_awesome_custom_resource.secret
}
```

## Data Source: pingone_trusted_email_domain_dkim

### `id` computed attribute removed

The unnecessary `id` computed attribute has been removed.

### `region` computed attribute renamed

The `region` computed attribute has been renamed to `regions`.

### `region.token` computed attribute renamed

The `region.token` computed attribute has been renamed to `regions.tokens`.

## Data Source: pingone_trusted_email_domain_ownership

### `id` computed attribute removed

The unnecessary `id` computed attribute has been removed.

### `region` computed attribute renamed

The `region` computed attribute has been renamed to `regions`.

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