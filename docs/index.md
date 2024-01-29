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

### Authenticate using static OAuth 2.0 Client Credentials (PingOne Worker Application)

```terraform
terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = "~> 0.26"
    }
  }
}

provider "pingone" {
  client_id      = var.client_id
  client_secret  = var.client_secret
  environment_id = var.environment_id
  region         = var.region

  force_delete_production_type = false
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
      version = "~> 0.26"
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
$ export PINGONE_CLIENT_ID="admin-client-id-value"
$ export PINGONE_CLIENT_SECRET="admin-client-secret-value"
$ export PINGONE_ENVIRONMENT_ID="admin-environment-id-value"
$ export PINGONE_REGION="admin-environment-region-code"
$ terraform plan
```

### Authenticate using an environment variable access token

```terraform
terraform {
  required_providers {
    pingone = {
      source  = "pingidentity/pingone"
      version = "~> 0.26"
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
$ export PINGONE_API_ACCESS_TOKEN="worker-access-token-value"
$ export PINGONE_REGION="admin-environment-region-code"
$ terraform plan
```

## Custom User Agent information

The PingOne provider allows custom information to be appended to the default user agent string (that includes Terraform provider version information) by setting the `PINGONE_TF_APPEND_USER_AGENT` environment variable.  This can be useful when troubleshooting issues with Ping Identity Support, or adding context to HTTP requests.

```shell
$ export PINGONE_TF_APPEND_USER_AGENT="Jenkins/2.426.2"
```

## Provider Schema Reference

- `api_access_token` (String) The access token used for provider resource management against the PingOne management API.  Default value can be set with the `PINGONE_API_ACCESS_TOKEN` environment variable.  Must provide only one of `api_access_token` (when obtaining the worker token outside of the provider) and `client_id` (when the provider should fetch the worker token during operations).
- `client_id` (String) Client ID for the worker app client.  Default value can be set with the `PINGONE_CLIENT_ID` environment variable.  Must provide only one of `api_access_token` (when obtaining the worker token outside of the provider) and `client_id` (when the provider should fetch the worker token during operations).  Must be configured with `client_secret` and `environment_id`.
- `client_secret` (String) Client secret for the worker app client.  Default value can be set with the `PINGONE_CLIENT_SECRET` environment variable.  Must be configured with `client_id` and `environment_id`.
- `environment_id` (String) Environment ID for the worker app client.  Default value can be set with the `PINGONE_ENVIRONMENT_ID` environment variable.  Must be configured with `client_id` and `client_secret`.
- `force_delete_production_type` (Boolean) Choose whether to force-delete any configuration that has a `PRODUCTION` type parameter.  The platform default is that `PRODUCTION` type configuration will not destroy without intervention to protect stored data.  By default this parameter is set to `false` and can be overridden with the `PINGONE_FORCE_DELETE_PRODUCTION_TYPE` environment variable.
- `http_proxy` (String) Full URL for the http/https proxy service, for example `http://127.0.0.1:8090`.  Default value can be set with the `HTTP_PROXY` or `HTTPS_PROXY` environment variables.
- `region` (String) The PingOne region to use.  Options are `AsiaPacific` `Canada` `Europe` and `NorthAmerica`.  Default value can be set with the `PINGONE_REGION` environment variable.
- `service_endpoints` (Block List) A single block containing configuration items to override the service API endpoints of PingOne. (see [below for nested schema](#nestedblock--service_endpoints))

<a id="nestedblock--service_endpoints"></a>
### Nested Schema for `service_endpoints`

Required:

- `api_hostname` (String) Hostname for the PingOne management service API.  Default value can be set with the `PINGONE_API_SERVICE_HOSTNAME` environment variable.
- `auth_hostname` (String) Hostname for the PingOne authentication service API.  Default value can be set with the `PINGONE_AUTH_SERVICE_HOSTNAME` environment variable.