---
layout: ""
page_title: "Version 1 Upgrade Guide (from version 0)"
description: |-
  Version 1.0.0 of the PingOne Terraform provider is a major release that introduces breaking changes to existing HCL.  This guide describes the changes that are required to upgrade v0.* PingOne Terraform provider releases to v1.0.0 onwards.
---

# PingOne Terraform Provider Version 1 Upgrade Guide (from version 0)

Version `1.0` of the PingOne Terraform provider is a major release that introduces breaking changes to existing HCL. This guide describes the changes that are required to upgrade `v0.*` PingOne Terraform provider releases to `v1.*`.

## Why have schemas changed?

As part of ensuring the ongoing maintainability of the Terraform provider integration and to solve functional issues, some resource schemas have changed going to version `1` from version `0`.

The schemas may have changed in the following ways:

* Removal of previously deprecated fields
* Removal of previously deprecated resources/data-sources
* Re-alignment of the Terraform schema with the API schema
* Data type changes (which change how HCL is written)

The following sections detail the rationale for the above changes, and whether the changes are routine for a major version upgrade or one off changes that aren't expected in future major version changes.

### Removal of previously deprecated fields

Removal of deprecated fields are expected on each major release going forward.  Ping maintains a deprecation and release strategy according to [Terraform provider creation best practices](https://developer.hashicorp.com/terraform/plugin/best-practices/versioning).

### Re-alignment of the Terraform schema with the API schema for some resources

Some resources have differing schemas between the Terraform HCL schema and the underlying API schema.  Re-alignment to the API simplifies writing HCL where developers are familiar with the API schema but also reduces ongoing maintenance overheads on the provider going forward.

Some re-alignment changes have been necessary in this major release. Further re-alignment to the API will be handled through backward-compatible deprecation, allowing a gradual approach to cutover.

### Data type changes

The most common data type change is changing block type fields (`example { ... }`) to nested object fields (`example = { ... }` for single objects or `example = [ { ... }, { ... } ]` for lists/sets of objects).  These data type changes ensure that the provider remains aligned with Terraform's strategic integration Framework SDK, which ensures that customers immediately benefit from the latest Terraform features as they become available.

These data type changes are a one-time set of changes from `v0` to `v1` and are not expected in future major releases.

## Provider Configuration Changes

### Major Version Change

Customers can keep operating existing `v0.*` releases until ready to upgrade to `v1.*`.  Remaining on the latest `v0.*` release can be achieved using the following syntax:

```terraform
terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = "~> 0.29"
    }
  }
}

provider "pingone" {
  client_id      = var.client_id
  client_secret  = var.client_secret
  environment_id = var.environment_id
  region         = var.region
}
```

It is highly recommended to go through the guide and make updates to each impacted resource before changing the version, as there are backward-incompatible changes.  Once ready to upgrade, the version can be incremented as follows:

```terraform
terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = "~> 1.0"
    }
  }
}

provider "pingone" {
  client_id      = var.client_id
  client_secret  = var.client_secret
  environment_id = var.environment_id
  region_code    = var.region_code
}
```

Ping recommends using [Provider version control](https://terraform.pingidentity.com/best-practices/#use-provider-version-control), detailed in the [Terraform best practices guide](https://terraform.pingidentity.com/best-practices/).

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

### `oidc_options.allow_wildcards_in_redirect_uris` optional parameter renamed

The `oidc_options.allow_wildcards_in_redirect_uris` optional parameter has been renamed to `oidc_options.allow_wildcard_in_redirect_uris` to align with the API.

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options {
    # ... other configuration parameters

    allow_wildcards_in_redirect_uris = "NONE"
  }
}
```

New configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options = {
    # ... other configuration parameters

    allow_wildcard_in_redirect_uris = "NONE"
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
  my_awesome_application_client_id     = pingone_application.my_awesome_application.oidc_options.client_id
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
  my_awesome_application_client_id     = pingone_application.my_awesome_application.oidc_options.client_id
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

### `oidc_options.token_endpoint_authn_method` optional parameter renamed

The `oidc_options.token_endpoint_authn_method` optional parameter has been renamed to `oidc_options.token_endpoint_auth_method` to align with the API.

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options {
    # ... other configuration parameters

    token_endpoint_authn_method = "NONE"
  }
}
```

New configuration example:

```terraform
resource "pingone_application" "my_awesome_application" {
  # ... other configuration parameters

  oidc_options = {
    # ... other configuration parameters

    token_endpoint_auth_method = "NONE"
  }
}
```

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

This parameter was previously deprecated and has now been made read only.  Use the `resource_type` and `custom_resource_id` fields going forward.

Previous configuration example (OIDC):

```terraform
resource "pingone_application_resource_grant" "my_awesome_resource_grant" {
  # ... other configuration parameters

  resource_id = var.my_oidc_resource_id
}
```

New configuration example (OIDC):

```terraform
resource "pingone_application" "my_awesome_saml_app" {
  # ... other configuration parameters

  resource_type = "OPENID_CONNECT"
}
```

Previous configuration example (Custom resource):

```terraform
resource "pingone_application_resource_grant" "my_awesome_resource_grant" {
  # ... other configuration parameters

  resource_id = var.my_custom_resource_id
}
```

New configuration example (Custom resource):

```terraform
resource "pingone_application" "my_awesome_saml_app" {
  # ... other configuration parameters

  resource_type      = "CUSTOM"
  custom_resource_id = var.my_custom_resource_id
}
```

### `resource_name` parameter removed

This parameter was previously required and has now been removed.  Use the `resource_type` and `custom_resource_id` fields going forward.

Previous configuration example (OIDC):

```terraform
resource "pingone_application_resource_grant" "my_awesome_resource_grant" {
  # ... other configuration parameters

  resource_name = "openid"
}
```

New configuration example (OIDC):

```terraform
resource "pingone_application" "my_awesome_saml_app" {
  # ... other configuration parameters

  resource_type = "OPENID_CONNECT"
}
```

Previous configuration example (Custom resource):

```terraform
resource "pingone_application_resource_grant" "my_awesome_resource_grant" {
  # ... other configuration parameters

  resource_name = var.my_custom_resource_name
}
```

New configuration example (Custom resource):

```terraform
resource "pingone_application" "my_awesome_saml_app" {
  # ... other configuration parameters

  resource_type      = "CUSTOM"
  custom_resource_id = var.my_custom_resource_id
}
```

### `scopes` parameter changed

This parameter was previously deprecated and has now been made required.

### `scope_names` parameter changed

This parameter was previously optional and has now been removed.  Use the `scopes` field going forward.

Previous configuration example (OIDC):

```terraform
resource "pingone_application_resource_grant" "my_awesome_resource_grant" {
  # ... other configuration parameters

  resource_name = "openid"

  scope_names [
    "email",
    "profile",
  ]
}
```

New configuration example (OIDC):

```terraform
resource "pingone_application" "my_awesome_saml_app" {
  # ... other configuration parameters

  resource_type = "OPENID_CONNECT"

  scopes [
    var.email_scope_id,
    var.profile_scope_id,
  ]
}
```

Previous configuration example (Custom resource):

```terraform
resource "pingone_application_resource_grant" "my_awesome_resource_grant" {
  # ... other configuration parameters

  resource_name = var.my_custom_resource_name

  scope_names [
    "myscope1",
    "myscope2",
  ]
}
```

New configuration example (Custom resource):

```terraform
resource "pingone_application" "my_awesome_saml_app" {
  # ... other configuration parameters

  resource_type      = "CUSTOM"
  custom_resource_id = var.my_custom_resource_id

  scope_names [
    var.my_custom_scope1_id,
    var.my_custom_scope2_id,
  ]
}
```

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

### `region` parameter ENUM values changed

The `region` parameter's values now aligns with the API request/response payload values.  Possible values are `AP` (for Asia-Pacific `.asia` environments), `AU` (for Asia-Pacific `.com.au` environments), `CA` (for Canada `.ca` environments), `EU` (for Europe `.eu` environments) and `NA` (for North America `.com` environments).

Previous configuration example:

```terraform
resource "pingone_environment" "my_environment" {
  # ... other configuration parameters
  
  region = "NorthAmerica"
}
```

New configuration example:

```terraform
resource "pingone_environment" "my_environment" {
  # ... other configuration parameters
  
  region = "NA"
}
```

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

### `kerberos_service_account_upn` parameter moved

The `kerberos_service_account_upn` parameter has been moved to `kerberos.service_account_upn`

Previous configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  # ... other configuration parameters

  kerberos_service_account_upn              = var.kerberos_service_account_upn
  kerberos_service_account_password         = var.kerberos_service_account_password
  kerberos_retain_previous_credentials_mins = 20
}
```

New configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  # ... other configuration parameters

  kerberos = {
    service_account_upn              = var.kerberos_service_account_upn
    service_account_password         = var.kerberos_service_account_password
    retain_previous_credentials_mins = 20
  }
}
```

### `kerberos_service_account_password` parameter moved

The `kerberos_service_account_password` parameter has been moved to `kerberos.service_account_password`

Previous configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  # ... other configuration parameters

  kerberos_service_account_upn              = var.kerberos_service_account_upn
  kerberos_service_account_password         = var.kerberos_service_account_password
  kerberos_retain_previous_credentials_mins = 20
}
```

New configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  # ... other configuration parameters

  kerberos = {
    service_account_upn              = var.kerberos_service_account_upn
    service_account_password         = var.kerberos_service_account_password
    retain_previous_credentials_mins = 20
  }
}
```

### `kerberos_retain_previous_credentials_mins` parameter moved

The `kerberos_retain_previous_credentials_mins` parameter has been moved to `kerberos.retain_previous_credentials_mins`

Previous configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  # ... other configuration parameters

  kerberos_service_account_upn              = var.kerberos_service_account_upn
  kerberos_service_account_password         = var.kerberos_service_account_password
  kerberos_retain_previous_credentials_mins = 20
}
```

New configuration example:

```terraform
resource "pingone_gateway" "my_awesome_gateway" {
  # ... other configuration parameters

  kerberos = {
    service_account_upn              = var.kerberos_service_account_upn
    service_account_password         = var.kerberos_service_account_password
    retain_previous_credentials_mins = 20
  }
}
```

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

The `user_type` parameter has been renamed to `user_types` and is now a map of objects data type and no longer a block set of objects type, to ensure that differences of embedded objects can be correctly correlated between plan and state (ref: [Issue #753](https://github.com/pingidentity/terraform-provider-pingone/issues/753)).  The map key of the new data type is the name of the user type (previously the `user_type.name` parameter).

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

### `saml.idp_verification_certificate_ids` parameter moved

The `saml.idp_verification_certificate_ids` parameter has been moved to `saml.idp_verification.certificates.*.id`

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  saml {
    idp_verification_certificate_ids = [
      var.saml_idp_verification_certificate_id_1,
      var.saml_idp_verification_certificate_id_2,
    ]
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  saml = {
    idp_verification = {
      certificates = [
        {
          id = var.saml_idp_verification_certificate_id_1
        },
        {
          id = var.saml_idp_verification_certificate_id_2
        }
      ]
    }
  }
}
```

### `saml.sp_signing_key_id` parameter moved

The `saml.sp_signing_key_id` parameter has been moved to `saml.sp_signing.key.id`

Previous configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  saml {
    sp_signing_key_id = var.saml_sp_signing_key_id
  }
}
```

New configuration example:

```terraform
resource "pingone_identity_provider" "my_awesome_identity_provider" {
  # ... other configuration parameters

  saml = {
    sp_signing = {
      key = {
        id = var.saml_sp_signing_key_id
      }
    }
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

## Resource: pingone_notification_template_content

### `email` parameter data type change

The `email` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_notification_template_content" "email" {
  # ... other configuration parameters

  email {
    body    = "Please approve this transaction with passcode $${otp}."
    subject = "BX Retail Transaction Request"

    from {
      name    = "BX Retail"
      address = "noreply@bxretail.org"
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_notification_template_content" "email" {
  # ... other configuration parameters

  email = {
    body    = "Please approve this transaction with passcode $${otp}."
    subject = "BX Retail Transaction Request"

    from = {
      name    = "BX Retail"
      address = "noreply@bxretail.org"
    }
  }
}
```

### `email.from` parameter data type change

The `email.from` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_notification_template_content" "email" {
  # ... other configuration parameters

  email {
    body    = "Please approve this transaction with passcode $${otp}."
    subject = "BX Retail Transaction Request"

    from {
      name    = "BX Retail"
      address = "noreply@bxretail.org"
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_notification_template_content" "email" {
  # ... other configuration parameters

  email = {
    body    = "Please approve this transaction with passcode $${otp}."
    subject = "BX Retail Transaction Request"

    from = {
      name    = "BX Retail"
      address = "noreply@bxretail.org"
    }
  }
}
```

### `email.reply_to` parameter data type change

The `email.reply_to` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_notification_template_content" "email" {
  # ... other configuration parameters

  email {
    # ... other configuration parameters

    body    = "Please approve this transaction with passcode $${otp}."
    subject = "BX Retail Transaction Request"

    reply_to {
      name    = "BX Retail"
      address = "reply@bxretail.org"
    }
  }
}
```

New configuration example:

```terraform
resource "pingone_notification_template_content" "email" {
  # ... other configuration parameters

  email = {
    # ... other configuration parameters

    body    = "Please approve this transaction with passcode $${otp}."
    subject = "BX Retail Transaction Request"

    reply_to = {
      name    = "BX Retail"
      address = "reply@bxretail.org"
    }
  }
}
```

### `push` parameter data type change

The `push` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_notification_template_content" "email" {
  # ... other configuration parameters

  push {
    body  = "Please approve this transaction."
    title = "BX Retail Transaction Request"
  }
}
```

New configuration example:

```terraform
resource "pingone_notification_template_content" "email" {
  # ... other configuration parameters

  push = {
    body  = "Please approve this transaction."
    title = "BX Retail Transaction Request"
  }
}
```

### `sms` parameter data type change

The `sms` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_notification_template_content" "email" {
  # ... other configuration parameters

  sms {
    content = "Please approve this transaction with passcode $${otp}."
    sender  = "BX Retail"
  }
}
```

New configuration example:

```terraform
resource "pingone_notification_template_content" "email" {
  # ... other configuration parameters

  sms = {
    content = "Please approve this transaction with passcode $${otp}."
    sender  = "BX Retail"
  }
}
```

### `voice` parameter data type change

The `voice` parameter is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
resource "pingone_notification_template_content" "email" {
  # ... other configuration parameters

  voice {
    content = "Hello <pause1sec> your authentication code is <sayCharValue>$${otp}</sayCharValue><pause1sec><pause1sec><repeatMessage val=2>I repeat <pause1sec>your code is <sayCharValue>$${otp}</sayCharValue></repeatMessage>"
    type    = "Alice"
  }
}
```

New configuration example:

```terraform
resource "pingone_notification_template_content" "email" {
  # ... other configuration parameters

  voice = {
    content = "Hello <pause1sec> your authentication code is <sayCharValue>$${otp}</sayCharValue><pause1sec><pause1sec><repeatMessage val=2>I repeat <pause1sec>your code is <sayCharValue>$${otp}</sayCharValue></repeatMessage>"
    type    = "Alice"
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

### `min_characters` parameter data type change

The `min_characters` parameter is now a nested object type and no longer a block list type.

Previous configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  min_characters {
    alphabetical_uppercase = 1
    alphabetical_lowercase = 1
    numeric                = 1
    special_characters     = 1
  }
}
```

New configuration example:

```terraform
resource "pingone_password_policy" "my_awesome_password_policy" {
  # ... other configuration parameters

  min_characters = {
    alphabetical_uppercase = 1
    alphabetical_lowercase = 1
    numeric                = 1
    special_characters     = 1
  }
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

  history = {
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

This parameter was previously deprecated and has now been made read only.  Use the `resource_type` and `custom_resource_id` parameters going forward.

Previous configuration example (OIDC):

```terraform
resource "pingone_resource_attribute" "my_awesome_resource_attribute" {
  # ... other configuration parameters

  resource_id = var.my_oidc_resource_id
}
```

New configuration example (OIDC):

```terraform
resource "pingone_resource_attribute" "my_awesome_resource_attribute" {
  # ... other configuration parameters

  resource_type = "OPENID_CONNECT"
}
```

Previous configuration example (Custom resource):

```terraform
resource "pingone_resource_attribute" "my_awesome_resource_attribute" {
  # ... other configuration parameters

  resource_id = var.my_custom_resource_id
}
```

New configuration example (Custom resource):

```terraform
resource "pingone_resource_attribute" "my_awesome_resource_attribute" {
  # ... other configuration parameters

  resource_type      = "CUSTOM"
  custom_resource_id = var.my_custom_resource_id
}
```

### `resource_name` parameter removed

This parameter was previously optional and has now been removed.  Use the `resource_type` and `custom_resource_id` fields going forward.

Previous configuration example (OIDC):

```terraform
resource "pingone_resource_attribute" "my_awesome_resource_attribute" {
  # ... other configuration parameters

  resource_name = "openid"
}
```

New configuration example (OIDC):

```terraform
resource "pingone_resource_attribute" "my_awesome_resource_attribute" {
  # ... other configuration parameters

  resource_type = "OPENID_CONNECT"
}
```

Previous configuration example (Custom resource):

```terraform
resource "pingone_resource_attribute" "my_awesome_resource_attribute" {
  # ... other configuration parameters

  resource_name = var.my_custom_resource_name
}
```

New configuration example (Custom resource):

```terraform
resource "pingone_resource_attribute" "my_awesome_resource_attribute" {
  # ... other configuration parameters

  resource_type      = "CUSTOM"
  custom_resource_id = var.my_custom_resource_id
}
```

## Resource: pingone_role_assignment_user (now pingone_user_role_assignment)

### Resource renamed to `pingone_user_role_assignment`

The `pingone_role_assignment_user` resource has been renamed to `pingone_user_role_assignment` to better align with other resources of the same nature.

## Resource: pingone_schema_attribute

### `schema_id` parameter changed

This parameter was previously deprecated and has now been made read only.  The default environment schema will be implicitly selected.

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

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.access_control_group_options[0].type
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.access_control_group_options.type
}
```

### `external_link_options` computed attribute data type change

The `external_link_options` computed attribute is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.external_link_options[0].home_page_url
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.external_link_options.home_page_url
}
```

### `icon` computed attribute data type change

The `icon` computed attribute is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.icon[0].href
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.icon.href
}
```

### `oidc_options` computed attribute data type change

The `oidc_options` computed attribute is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options[0].type
  client_id = data.pingone_application.example.oidc_options[0].client_id
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options.type
  client_id = data.pingone_application.example.oidc_options.client_id
}
```

### `oidc_options.allow_wildcards_in_redirect_uris` computed attribute renamed

The `oidc_options.allow_wildcards_in_redirect_uris` computed attribute has been renamed to `oidc_options.allow_wildcard_in_redirect_uris` to align with the API.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options[0].allow_wildcards_in_redirect_uris
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options.allow_wildcard_in_redirect_uris
}
```

### `oidc_options.certificate_based_authentication` computed attribute data type change

The `oidc_options.certificate_based_authentication` computed attribute is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options[0].certificate_based_authentication[0].key_id
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options.certificate_based_authentication.key_id
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
  my_awesome_application_client_id     = data.pingone_application.my_awesome_application.oidc_options.client_id
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
  my_awesome_application_client_id     = data.pingone_application.my_awesome_application.oidc_options.client_id
  my_awesome_application_client_secret = data.pingone_application_secret.my_awesome_application.secret
}
```

### `oidc_options.cors_settings` computed attribute data type change

The `oidc_options.cors_settings` computed attribute is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options[0].cors_settings[0].origins
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options.cors_settings.origins
}
```

### `oidc_options.mobile_app` computed attribute data type change

The `oidc_options.mobile_app` computed attribute is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options[0].mobile_app[0].bundle_id
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options.mobile_app.bundle_id
}
```

### `oidc_options.mobile_app.integrity_detection` computed attribute data type change

The `oidc_options.mobile_app.integrity_detection` computed attribute is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options[0].mobile_app[0].integrity_detection[0].enabled
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options.mobile_app.integrity_detection.enabled
}
```

### `oidc_options.mobile_app.integrity_detection.cache_duration` computed attribute data type change

The `oidc_options.mobile_app.integrity_detection.cache_duration` computed attribute is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options[0].mobile_app[0].integrity_detection[0].cache_duration[0].amount
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options.mobile_app.integrity_detection.cache_duration.amount
}
```

### `oidc_options.mobile_app.integrity_detection.google_play` computed attribute data type change

The `oidc_options.mobile_app.integrity_detection.google_play` computed attribute is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options[0].mobile_app[0].integrity_detection[0].google_play[0].verification_type
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options.mobile_app.integrity_detection.google_play.verification_type
}
```

### `oidc_options.token_endpoint_authn_method` computed attribute renamed

The `oidc_options.token_endpoint_authn_method` computed attribute has been renamed to `oidc_options.token_endpoint_auth_method` to align with the API.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options[0].token_endpoint_authn_method
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.oidc_options.token_endpoint_auth_method
}
```

### `saml_options` computed attribute data type change

The `saml_options` computed attribute is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.saml_options[0].type
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.saml_options.type
}
```

### `saml_options.cors_settings` computed attribute data type change

The `saml_options.cors_settings` computed attribute is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.saml_options[0].cors_settings[0].origins
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.saml_options.cors_settings.origins
}
```

### `saml_options.idp_signing_key` computed attribute data type change

The `saml_options.idp_signing_key` computed attribute is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.saml_options[0].idp_signing_key[0].key_id
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.saml_options.idp_signing_key.key_id
}
```

### `saml_options.sp_verification` computed attribute data type change

The `saml_options.sp_verification` computed attribute is now a nested object type and no longer a list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.saml_options[0].sp_verification[0].authn_request_signed
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_application.example.saml_options.sp_verification.authn_request_signed
}
```

### `saml_options.sp_verification_certificate_ids` computed attribute removed

This parameter was previously deprecated and has been removed.  Use the `saml_options.sp_verification.certificate_ids` attribute going forward.

## Data Source: pingone_environment

### `region` parameter ENUM values changed

The `region` computed attribute's values now aligns with the API request/response payload values.  Possible values are `AP` (for Asia-Pacific `.asia` environments), `AU` (for Asia-Pacific `.com.au` environments), `CA` (for Canada `.ca` environments), `EU` (for Europe `.eu` environments) and `NA` (for North America `.com` environments).

### `service` computed attribute rename and data type change

The `service` computed attribute has been renamed to `services` and is now a nested list type and no longer a block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_environment.example.service[0].type
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_environment.example.services[0].type
}
```

### `service.bookmark` computed attribute rename and data type change

The `service.bookmark` computed attribute has been renamed to `services.bookmarks` and is now a nested list type and no longer a block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_environment.example.service[0].bookmark[0].name
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_environment.example.services[0].bookmarks[0].name
}
```

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

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_flow_policy.example.davinci_application[0].id
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_flow_policy.example.davinci_application.id
}
```

### `trigger` computed attribute data type change

The `trigger` computed attribute is now a nested object type and no longer a list block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_flow_policy.example.trigger[0].type
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_flow_policy.example.trigger.type
}
```

## Data Source: pingone_gateway

### `radius_client` computed attribute rename and data type change

The `radius_client` computed attribute has been renamed to `radius_clients` and is now a set of objects data type and no longer a block set type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_gateway.example.radius_client[0].ip
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_gateway.example.radius_clients[0].ip
}
```

### `user_type` computed attribute rename and data type change

The `user_type` computed attribute has been renamed to `user_types` and is now a map of objects data type and no longer a block set of objects type.  The map key of the new data type is the name of the user type (previously the `user_type.name` parameter).

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_gateway.example.user_type[0].id
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_gateway.example.user_types["My user types"].id
}
```

### `user_type.push_password_changes_to_ldap` computed attribute rename

The `user_type.push_password_changes_to_ldap` computed attribute has been renamed to `user_types.allow_password_changes`.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_gateway.example.user_type[0].push_password_changes_to_ldap
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_gateway.example.user_types["My user types"].allow_password_changes
}
```

### `user_type.user_migration` computed attribute rename and data type change

The `user_type.user_migration` computed attribute has been renamed to `user_types.new_user_lookup` and is now a single object data type and no longer a block set type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_gateway.example.user_type[0].user_migration[0].population_id
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_gateway.example.user_types["My user types"].new_user_lookup.population_id
}
```

### `user_type.user_migration.attribute_mapping` computed attribute rename and data type change

The `user_type.user_migration.attribute_mapping` computed attribute has been renamed to `user_types.new_user_lookup.attribute_mappings` and is now a set of objects data type and no longer a block set type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_gateway.example.user_type[0].user_migration[0].attribute_mapping[0].name
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_gateway.example.user_types["My user types"].new_user_lookup.attribute_mappings[0].name
}
```

### `user_type.user_migration.lookup_filter_pattern` computed attribute rename

The `user_type.user_migration.lookup_filter_pattern` computed attribute has been renamed to `user_types.new_user_lookup.ldap_filter_pattern`.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_gateway.example.user_type[0].user_migration[0].lookup_filter_pattern
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_gateway.example.user_types["My user types"].new_user_lookup.ldap_filter_pattern
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

## Data Source: pingone_license

### `advanced_services` computed attribute data type change

The `advanced_services` computed attribute is now a nested object type and no longer a list block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.advanced_services[0].pingid[0].included
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.advanced_services.pingid.included
}
```

### `advanced_services.pingid` computed attribute data type change

The `advanced_services.pingid` computed attribute is now a nested object type and no longer a list block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.advanced_services[0].pingid[0].included
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.advanced_services.pingid.included
}
```

### `authorize` computed attribute data type change

The `authorize` computed attribute is now a nested object type and no longer a list block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.authorize[0].allow_api_access_management
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.authorize.allow_api_access_management
}
```

### `credentials` computed attribute data type change

The `credentials` computed attribute is now a nested object type and no longer a list block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.credentials[0].allow_credentials
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.credentials.allow_credentials
}
```

### `environments` computed attribute data type change

The `environments` computed attribute is now a nested object type and no longer a list block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.environments[0].allow_add_resources
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.environments.allow_add_resources
}
```

### `fraud` computed attribute data type change

The `fraud` computed attribute is now a nested object type and no longer a list block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.fraud[0].allow_bot_malicious_device_detection
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.fraud.allow_bot_malicious_device_detection
}
```

### `gateways` computed attribute data type change

The `gateways` computed attribute is now a nested object type and no longer a list block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.gateways[0].allow_ldap_gateway
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.gateways.allow_ldap_gateway
}
```

### `intelligence` computed attribute data type change

The `intelligence` computed attribute is now a nested object type and no longer a list block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.intelligence[0].allow_advanced_predictors
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.intelligence.allow_advanced_predictors
}
```

### `mfa` computed attribute data type change

The `mfa` computed attribute is now a nested object type and no longer a list block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.mfa[0].allow_push_notification
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.mfa.allow_push_notification
}
```

### `orchestrate` computed attribute data type change

The `orchestrate` computed attribute is now a nested object type and no longer a list block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.orchestrate[0].allow_orchestration
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.orchestrate.allow_orchestration
}
```

### `users` computed attribute data type change

The `users` computed attribute is now a nested object type and no longer a list block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.users[0].allow_password_management_notifications
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.users.allow_password_management_notifications
}
```

### `verify` computed attribute data type change

The `verify` computed attribute is now a nested object type and no longer a list block type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.verify[0].allow_push_notifications
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_license.example.verify.allow_push_notifications
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

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.account_lockout[0].duration_seconds
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.lockout.duration_seconds
}
```

### `account_lockout.fail_count` computed attribute renamed

The `account_lockout.fail_count` computed attribute has been renamed to `lockout.failure_count`.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.account_lockout[0].fail_count
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.lockout.failure_count
}
```

### `bypass_policy` computed attribute removed

The `bypass_policy` computed attribute has no effect and has been removed.

### `environment_default` computed attribute renamed

The `environment_default` computed attribute has been renamed to `default`.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.environment_default
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.default
}
```

### `exclude_commonly_used_passwords` computed attribute renamed

The `exclude_commonly_used_passwords` computed attribute has been renamed to `excludes_commonly_used_passwords`.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.exclude_commonly_used_passwords
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.excludes_commonly_used_passwords
}
```

### `exclude_profile_data` computed attribute renamed

The `exclude_profile_data` computed attribute has been renamed to `excludes_profile_data`.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.exclude_profile_data
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.excludes_profile_data
}
```

### `password_age.max` computed attribute moved

The `password_age.max` computed attribute has been moved to `password_age_max`.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.password_age[0].max
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.password_age_max
}
```

### `password_age.min` computed attribute moved

The `password_age.min` computed attribute has been moved to `password_age_min`.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.password_age[0].min
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.password_age_min
}
```

### `password_history` computed attribute rename and data type change

The `password_history` computed attribute has been renamed to `history` and is now a nested object type and no longer a block list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.password_history[0].count
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.history.count
}
```

### `password_history.prior_password_count` computed attribute renamed

The `password_history.prior_password_count` computed attribute has been renamed to `history.count`.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.password_history[0].count
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.history.count
}
```

### `password_length` computed attribute rename and data type change

The `password_length` computed attribute has been renamed to `length` and is now a nested object type and no longer a block list type.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.password_length[0].max
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_password_policy.example.length.max
}
```

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

## Data Source: pingone_resource_scope

### `resource_id` parameter changed

This parameter was previously required and has now been made read only.  Use the `resource_type` and `custom_resource_id` parameters going forward.

Previous configuration example (OIDC):

```terraform
data "pingone_resource_scope" "my_awesome_resource_scope" {
  # ... other configuration parameters

  resource_id = var.my_oidc_resource_id
}
```

New configuration example (OIDC):

```terraform
data "pingone_resource_scope" "my_awesome_resource_scope" {
  # ... other configuration parameters

  resource_type = "OPENID_CONNECT"
}
```

Previous configuration example (Custom resource):

```terraform
data "pingone_resource_scope" "my_awesome_resource_scope" {
  # ... other configuration parameters

  resource_id = var.my_custom_resource_id
}
```

New configuration example (Custom resource):

```terraform
data "pingone_resource_scope" "my_awesome_resource_scope" {
  # ... other configuration parameters

  resource_type      = "CUSTOM"
  custom_resource_id = var.my_custom_resource_id
}
```

## Data Source: pingone_trusted_email_domain_dkim

### `id` computed attribute removed

The unnecessary `id` computed attribute has been removed.

### `region` computed attribute renamed

The `region` computed attribute has been renamed to `regions`.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_trusted_email_domain_dkim.example.region[0].name
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_trusted_email_domain_dkim.example.regions[0].name
}
```

### `region.token` computed attribute renamed

The `region.token` computed attribute has been renamed to `regions.tokens`.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_trusted_email_domain_dkim.example.region[0].token[0].key
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_trusted_email_domain_dkim.example.regions[0].tokens[0].key
}
```

## Data Source: pingone_trusted_email_domain_ownership

### `id` computed attribute removed

The unnecessary `id` computed attribute has been removed.

### `region` computed attribute renamed

The `region` computed attribute has been renamed to `regions`.

Previous configuration example:

```terraform
locals {
  demo_var = data.pingone_trusted_email_domain_ownership.example.region[0].name
}
```

New configuration example:

```terraform
locals {
  demo_var = data.pingone_trusted_email_domain_ownership.example.regions[0].name
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