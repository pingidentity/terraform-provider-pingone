---
layout: ""
page_title: "Using Ping CLI authentication with the provider"
description: |-
  Authenticate the PingOne Terraform provider with a Ping CLI login session stored in the OS keychain.
---

# Using Ping CLI authentication with the provider

The provider can authenticate with the access token that [Ping CLI](https://github.com/pingidentity/pingcli) stores in the OS keychain after a Ping CLI login. This avoids managing client credentials for Terraform use cases.

-> Ping CLI authentication is an additional option. Authenticating with `client_id`, `client_secret` and `environment_id` (or `api_access_token`) is unchanged and remains fully supported.

## Prerequisites

- Ping CLI v1.x installed.
- A PingOne worker application with the PingOne roles Terraform requires. The application that Ping CLI logs in with must be a worker application: PingOne rejects management API calls from tokens of other application types, even for an organization admin. See the [PingOne provider setup guide](https://github.com/pingidentity/terraform-provider-pingone/blob/main/docs/index.md#authenticate-using-static-oauth-20-client-credentials-pingone-worker-application) for the application setup.

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

- `pingcli_profile` (environment `PINGCLI_PROFILE`) selects the profile whose keychain session is used. When unset, the Ping CLI default profile (`default`) is used.
- `pingcli_grant_type` (environment `PINGCLI_PINGONE_AUTHENTICATION_TYPE`) restricts the lookup to one grant type (`client_credentials`, `authorization_code` or `device_code`). When unset, each grant type is checked and the valid session that expires last is used.
- `region_code` (environment `PINGONE_REGION_CODE`) is **required** in this mode.

If `use_pingcli_auth` is set, `client_id`, `client_secret`, `environment_id` or `api_access_token` must not be set in the provider configuration, and credentials exported in the environment are ignored.

## Behaviour and constraints

- The provider reads the stored session; it never refreshes it and never writes to Ping CLI storage.
- The token is static for the duration of the run. Applies that outlive the token lifetime fail with an authentication error, exactly as with `api_access_token` today. Run `pingcli pingone auth refresh` (or log in again) and re-run Terraform.
- The feature is not usable on remote runners (for example Terraform Cloud) that have no local Ping CLI session or OS keychain.

## Error reference

| Error | Cause | Fix |
|---|---|---|
| `Conflicting provider authentication configuration` | `use_pingcli_auth` is combined with credentials set in the provider configuration | Remove `use_pingcli_auth` to use credentials, or remove the credential attributes to use Ping CLI |
| `Invalid Ping CLI grant type` | `pingcli_grant_type` is not `client_credentials`, `authorization_code` or `device_code` | Use one of the supported grant types |
| `Missing region code` | No `region_code` attribute or `PINGONE_REGION_CODE` environment variable | Set the region for the PingOne tenant |
| `no Ping CLI authentication session was found in the OS keychain` | No session for the profile, or the session is under a different profile than expected | Run `pingcli pingone auth login`; note the session is stored under the profile active at login time |
| `has a format this provider does not recognise` | The keychain entry has neither an access token nor Ping CLI's split-session marker, for example because a newer Ping CLI changed the layout | Update the provider, or authenticate with client credentials or an access token instead |
| `has expired` | The stored access token has passed its expiry (with a small safety margin) | Run `pingcli pingone auth refresh` or `pingcli pingone auth login` |
| `The actor attempting to perform the request is not authorized` (`ACCESS_FAILED`) | The Ping CLI login's application is not a worker application, or it has no PingOne admin role for the environment being managed | Use a worker application for the Ping CLI login and assign it the needed admin roles, or log in with a different profile or grant type |
| `failed to read the Ping CLI authentication session from the OS keychain` | The OS keychain is unavailable (for example headless Linux without a Secret Service) | Run where a keychain is available, or authenticate with client credentials or an access token instead |
