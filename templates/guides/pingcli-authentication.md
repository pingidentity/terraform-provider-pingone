---
layout: ""
page_title: "Using Ping CLI authentication with the provider"
description: |-
  Authenticate the PingOne Terraform provider with a Ping CLI login session stored in the OS keychain.
---

# Using Ping CLI authentication with the provider

The provider can authenticate with the access token that [Ping CLI](https://github.com/pingidentity/pingcli) stores in the OS keychain after a Ping CLI login. This avoids managing client credentials for Terraform use cases.

## Prerequisites

- Ping CLI v1.x installed.
- A PingOne worker application with the PingOne roles Terraform requires. The application that Ping CLI logs in with must be a worker application: PingOne rejects management API calls from tokens of other application types, with a "not authorized" (`ACCESS_FAILED`) error, even for an organization admin. See the [PingOne provider setup guide](https://github.com/pingidentity/terraform-provider-pingone/blob/main/docs/index.md#authenticate-using-static-oauth-20-client-credentials-pingone-worker-application) for the application setup.

## Step 1: Log in with Ping CLI

Log in to PingOne with the default keychain storage (`secure_local`). Storage types `file_system` and `none` are not read by the provider.

```sh
pingcli pingone auth login
```

Ping CLI stores the authentication session under the profile that was active when the login was performed. Verify with `pingcli pingone auth status`.

## Step 2: Configure the provider

```hcl
provider "pingone" {
  use_pingcli_auth = true
  region_code      = "NA"
}
```

Or, when the profile holds sessions for more than one grant type, name the profile and grant type to use:

```hcl
provider "pingone" {
  use_pingcli_auth    = true
  pingcli_profile     = "myprofile"
  pingcli_grant_type  = "device_code"
  region_code         = "NA"
}
```

- `use_pingcli_auth` (environment `PINGONE_TF_USE_PINGCLI_AUTH`) enables this mode. A value in the provider configuration takes precedence over the environment variable, so `use_pingcli_auth = false` turns it off even when the variable is `true`.
- `pingcli_profile` (environment `PINGCLI_PROFILE`) selects the profile whose keychain session is used. When unset, the Ping CLI default profile (`default`) is used.
- `pingcli_grant_type` (environment `PINGCLI_PINGONE_AUTHENTICATION_TYPE`) restricts the lookup to one grant type (`client_credentials`, `authorization_code` or `device_code`). When unset, each grant type is checked and the valid session that expires last is used.
- `region_code` (environment `PINGONE_REGION_CODE`) is **required** in this mode.

If `use_pingcli_auth` is set, `client_id`, `client_secret`, `environment_id` or `api_access_token` must not be set in the provider configuration or exported as the matching `PINGONE_CLIENT_ID`, `PINGONE_CLIENT_SECRET`, `PINGONE_ENVIRONMENT_ID` or `PINGONE_API_ACCESS_TOKEN` environment variables.

## Behaviour and constraints

- The provider reads the stored session; it never refreshes it and never writes to Ping CLI storage.
- The token is static for the duration of the run. Applies that outlive the token lifetime fail with an authentication error, exactly as with `api_access_token` today. Run `pingcli pingone auth refresh` (or log in again) and re-run Terraform.
- The feature is not usable on remote runners (for example Terraform Cloud) that have no local Ping CLI session or OS keychain.
