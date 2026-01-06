---
layout: ""
page_title: "Provider: PingOne"
description: |-
  This PingOne provider provides resources and data sources to manage the PingOne platform as infrastructure-as-code, through the PingOne management API.
---

# PingOne Provider

The PingOne provider interacts with the configuration of the PingOne platform via the management API. The provider requires credentials from worker application client before it can be used.

## Getting Started

To get started using the PingOne Terraform provider, first you'll need an active PingOne cloud subscription.  Get instant access with a [PingOne trial account](https://www.pingidentity.com/en/try-ping.html), or read more about Ping Identity at [pingidentity.com](https://www.pingidentity.com)

### Configure PingOne for Terraform access

For detailed instructions on how to prepare PingOne for Terraform access, see the [PingOne getting started guide](https://terraform.pingidentity.com/getting-started/pingone/#configure-pingone-for-terraform-access) at [terraform.pingidentity.com](https://terraform.pingidentity.com).

### Ping Identity Developer Experiences

For code examples showing how to configure the PingOne service using Terraform, the following resources are available:

- [Tutorials and use-case examples](https://terraform.pingidentity.com/examples/) at [terraform.pingidentity.com](https://terraform.pingidentity.com).
- End-to-end solution deployment examples, complete with sample applications at the [Ping Identity Developers Experience Github](https://github.com/pingidentity-developers-experience?tab=repositories).

## Provider Authentication

### Authenticate using PingCLI profile (stored token)

You can use a saved token from [PingCLI](https://github.com/pingidentity/pingcli) instead of managing client credentials directly in Terraform. This is convenient for interactive users and CI environments that already run `pingcli login`.

1) Use PingCLI to log in with your preferred grant type and profile (e.g., `authorization_code`, `device_code`, or `client_credentials`). Ensure the profile’s `authentication.type` in the PingCLI config matches the grant you used.

```shell
pingcli login --type authorization_code --profile dev
# or
pingcli login --type device_code --profile dev
# or
pingcli login --type client_credentials --profile dev
```

2) In Terraform, point the provider to your PingCLI config file and profile:

```terraform
terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = ">= 1.13, < 1.14"
    }
  }
}

provider "pingone" {
  # Path to your PingCLI config file and profile name
  config_path    = "~/.pingcli/config.yaml"
  config_profile = "dev"
}
```

Alternatively, you can set environment variables:

```shell
export PINGCLI_CONFIG="~/.pingcli/config.yaml"
export PINGCLI_PROFILE="dev"
terraform plan
```

Notes:
- The provider will first try to use a valid stored token from PingCLI (Keychain or `~/.pingcli/credentials`).
- If no stored token is found and `api_access_token` is set on the provider, it uses that token.
- If the PingCLI profile contains `client_credentials` (client ID/secret) and environment ID, the provider can use them to obtain a token.
- Otherwise, the provider returns an error prompting you to run `pingcli login`.

For full details and examples of the PingCLI config structure, see the guide: [Using PingCLI with the Provider](./guides/pingcli-authentication.md).

### Authenticate using static OAuth 2.0 Client Credentials (PingOne Worker Application)

```terraform
terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = ">= 1.14, < 1.15"
    }
  }
}

provider "pingone" {
  client_id      = var.client_id
  client_secret  = var.client_secret
  environment_id = var.environment_id
  region_code    = var.region_code
}

resource "pingone_environment" "my_environment" {
  # ...
}
```

### Authenticate using environment variable OAuth 2.0 Client Credentials (PingOne Worker Application)

```terraform
terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = ">= 1.14, < 1.15"
    }
  }
}

provider "pingone" {
}

resource "pingone_environment" "my_environment" {
  # ...
}
```

```shell
export PINGONE_CLIENT_ID="admin-client-id-value"
export PINGONE_CLIENT_SECRET="admin-client-secret-value"
export PINGONE_ENVIRONMENT_ID="admin-environment-id-value"
export PINGONE_REGION_CODE="AP | AU | CA | EU | NA | SG"
terraform plan
```

### Authenticate using an environment variable access token

```terraform
terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = ">= 1.14, < 1.15"
    }
  }
}

provider "pingone" {
}

resource "pingone_environment" "my_environment" {
  # ...
}
```

```shell
export PINGONE_API_ACCESS_TOKEN="worker-access-token-value"
export PINGONE_REGION_CODE="AP | AU | CA | EU | NA | SG"
terraform plan
```

## Custom User Agent information

The PingOne provider allows custom information to be appended to the default user agent string (that includes Terraform provider version information) by setting the `PINGONE_TF_APPEND_USER_AGENT` environment variable, or the `append_user_agent` provider parameter.  This can be useful when troubleshooting issues with Ping Identity Support, or adding context to HTTP requests.

```shell
export PINGONE_TF_APPEND_USER_AGENT="Jenkins/2.426.2"
```

## Global Options

The PingOne provider provides global options to override API behaviours in PingOne, for example to override data protection features for development, testing and demo use cases.

The following example shows how to configure the provider with global options to force-delete populations that contain users not managed by Terraform.  Note this only applies to environments that are of type `SANDBOX`.

```terraform
provider "pingone" {
  client_id      = var.client_id
  client_secret  = var.client_secret
  environment_id = var.environment_id
  region_code    = var.region_code

  global_options {

    population {
      // This option should not be used in environments that contain production data.  Data loss may occur.
      contains_users_force_delete = true
    }

  }
}
```

## Provider Schema Reference

- `api_access_token` (String) The access token used for provider resource management against the PingOne management API.  Default value can be set with the `PINGONE_API_ACCESS_TOKEN` environment variable.  Must provide only one of `api_access_token` (when obtaining the worker token outside of the provider) and `client_id` (when the provider should fetch the worker token during operations).
- `client_id` (String) Client ID for the worker app client.  Default value can be set with the `PINGONE_CLIENT_ID` environment variable.  Must provide only one of `api_access_token` (when obtaining the worker token outside of the provider) and `client_id` (when the provider should fetch the worker token during operations).  Must be configured with `client_secret` and `environment_id`.
- `client_secret` (String) Client secret for the worker app client.  Default value can be set with the `PINGONE_CLIENT_SECRET` environment variable.  Must be configured with `client_id` and `environment_id`.
- `environment_id` (String) Environment ID for the worker app client.  Default value can be set with the `PINGONE_ENVIRONMENT_ID` environment variable.  Must be configured with `client_id` and `client_secret`.
- `global_options` (Block List) A single block containing configuration items to override API behaviours in PingOne. (see [below for nested schema](#nestedblock--global_options))
- `http_proxy` (String) Full URL for the http/https proxy service, for example `http://127.0.0.1:8090`.  Default value can be set with the `HTTP_PROXY` or `HTTPS_PROXY` environment variables.
- `region_code` (String) The PingOne region to use, which selects the appropriate service endpoints.  Options are `AP` (for Asia-Pacific `.asia` tenants), `AU` (for Asia-Pacific `.com.au` tenants), `CA` (for Canada `.ca` tenants), `EU` (for Europe `.eu` tenants), `NA` (for North America `.com` tenants) and `SG` (for Singapore `.sg` tenants).  Default value can be set with the `PINGONE_REGION_CODE` environment variable.
- `service_endpoints` (Block List) A single block containing configuration items to override the service API endpoints of PingOne. (see [below for nested schema](#nestedblock--service_endpoints))
- `append_user_agent` (String) A custom string value to append to the end of the `User-Agent` header when making API requests to the PingOne service. Default value can be set with the `PINGONE_TF_APPEND_USER_AGENT` environment variable.
- `config_path` (String) Path to a PingCLI configuration file containing authentication credentials. Default value can be set with the `PINGCLI_CONFIG` environment variable. Cannot be used together with `client_id`, `client_secret`, `environment_id`, or `api_access_token`. If set, `config_profile` can optionally specify which profile to use (defaults to the active profile in the config file).
- `config_profile` (String) Name of the profile to use from the PingCLI configuration file. If not specified, uses the active profile defined in the config file. Default value can be set with the `PINGCLI_PROFILE` environment variable. Requires `config_path` to be set.

<a id="nestedblock--global_options"></a>
### Nested Schema for `global_options`

Optional:

- `population` (Block List) A single block containing configuration items to override population resource settings in PingOne. (see [below for nested schema](#nestedblock--global_options))

<a id="nestedblock--service_endpoints"></a>
### Nested Schema for `service_endpoints`

Required:

- `api_hostname` (String) Hostname for the PingOne management service API.  Default value can be set with the `PINGONE_API_SERVICE_HOSTNAME` environment variable.
- `auth_hostname` (String) Hostname for the PingOne authentication service API.  Default value can be set with the `PINGONE_AUTH_SERVICE_HOSTNAME` environment variable.

<a id="nestedblock--global_options-population"></a>
### Nested Schema for `global_options.population`

Optional:

- `contains_users_force_delete` (Boolean) Choose whether to force-delete populations that contain users not managed by Terraform.  Useful for development and testing use cases, and only applies if the environment that contains the population is of type `SANDBOX`, or the `global_options.environment.production_type_force_delete` parameter is set to `true`.  The platform default is that populations cannot be removed if they contain user data.  By default this parameter is set to `false`. This option should not be set to `true` when the environment contains production data. Data loss may occur.
