---
layout: ""
page_title: "Using Ping CLI authentication with the provider"
description: |-
  Authenticate the PingOne Terraform provider with a Ping CLI login session stored in the OS keychain.
---

# Using Ping CLI authentication with the provider

The provider can authenticate with the access token that [Ping CLI](https://github.com/pingidentity/pingcli) stores in the OS keychain after a Ping CLI login. This avoids creating a worker application and managing client credentials for Terraform use cases.

-> Ping CLI authentication is an additional option. Authenticating with `client_id`, `client_secret` and `environment_id` (or `api_access_token`) is unchanged and remains fully supported; the provider only uses a Ping CLI session when you opt in as described below.

-> This guide applies to Ping CLI v1.x. Ping CLI v0.8 used a different storage layout that the provider does not read.

## Prerequisites

- Ping CLI v1.x installed and logged in with its default keychain storage.
- For `authorization_code` and `device_code`, the PingOne application Ping CLI logs in with must be a **Worker** application with the admin roles Terraform needs (for example Environment Admin) assigned to the application. PingOne treats tokens from other application types as self-service tokens and rejects management API calls with `ACCESS_FAILED`, even when the signed-in user is an organization admin.
- A PingOne region configured for the provider (see below).

## Step 1: Log in with Ping CLI

Log in to PingOne with the default keychain storage (`secure_local`). Storage types `file_system` and `none` are not read by the provider. Use `pingcli pingone auth login` rather than `pingcli auth login`, which authenticates every enabled Ping service instead of only PingOne.

```sh
pingcli pingone auth login
```

Note that Ping CLI stores the authentication session under the profile that was **active** when the login was performed. A login invoked with `--profile dev` while another profile is active stores the session under the active profile, not `dev`. Verify with:

```sh
pingcli pingone auth status
```

## Step 2: Configure the provider

Point the provider at the Ping CLI profile:

```hcl
provider "pingone" {
  config_profile = "myprofile"
  region_code    = "NA"
}
```

If the profile holds sessions for more than one grant type, name the one to use:

```hcl
provider "pingone" {
  config_profile    = "myprofile"
  config_grant_type = "device_code"
  region_code       = "NA"
}
```

Or use environment variables:

```sh
export PINGCLI_PROFILE="myprofile"
export PINGCLI_PINGONE_AUTHENTICATION_TYPE="device_code" # optional
export PINGONE_REGION_CODE="NA"
```

- `config_profile` (environment `PINGCLI_PROFILE`) selects the profile whose keychain session is used. The provider derives the keychain entry from the profile name and does not read the Ping CLI configuration file, so the profile must be named explicitly.
- `config_grant_type` (environment `PINGCLI_PINGONE_AUTHENTICATION_TYPE`) optionally restricts the lookup to one grant type (`client_credentials`, `authorization_code` or `device_code`). Set it when a profile holds sessions for more than one grant type, for example after switching the profile's authentication method.
- `region_code` (environment `PINGONE_REGION_CODE`) is **required** in this mode; the provider does not derive the region from the Ping CLI configuration file.

If `config_profile` (or `config_grant_type`) is set in HCL, or `PINGCLI_PROFILE` is set in the environment and no client credentials are configured, the provider uses Ping CLI authentication automatically. `PINGCLI_PINGONE_AUTHENTICATION_TYPE` on its own does not enable it.

## Behaviour and constraints

- The keychain entry is derived from the profile name and the grant type. Unless `config_grant_type` is set, the provider checks `client_credentials`, `authorization_code` and `device_code` and uses the valid session that expires last. Client IDs and secrets are not part of the key and are not needed.
- The provider **reads** the stored session; it never refreshes it and never writes to Ping CLI storage.
- The token is static for the duration of the run. Applies that outlive the token lifetime fail with an authentication error, exactly as with `api_access_token` today. Run `pingcli pingone auth refresh` (or log in again) and re-run Terraform.
- Ping CLI authentication cannot be combined with `client_id`, `client_secret`, `environment_id` or `api_access_token` set in the provider configuration; that is an error.
- When `config_profile` or `config_grant_type` is set in the provider configuration, it wins over any `PINGONE_CLIENT_ID`, `PINGONE_CLIENT_SECRET`, `PINGONE_ENVIRONMENT_ID` or `PINGONE_API_ACCESS_TOKEN` exported in your shell: the provider ignores them. When Ping CLI authentication would be enabled only by `PINGCLI_PROFILE` in the environment, any credential (in the configuration or the environment) takes precedence instead, so an exported Ping CLI variable never breaks an existing credentials setup.
- The feature is not usable on remote runners (for example Terraform Cloud) that have no local Ping CLI session or OS keychain.

## Error reference

| Error | Cause | Fix |
|---|---|---|
| `Conflicting provider authentication configuration` | `config_profile` or `config_grant_type` is combined with credentials set in the provider configuration | Remove `config_profile`/`config_grant_type` to use credentials, or remove `client_id`, `client_secret`, `environment_id` and `api_access_token` from the provider configuration to use Ping CLI |
| `Missing Ping CLI profile` | `config_grant_type` was set without a profile | Set `config_profile` or `PINGCLI_PROFILE` |
| `Invalid Ping CLI grant type` | `config_grant_type` is not `client_credentials`, `authorization_code` or `device_code` | Use one of the supported grant types |
| `Unknown Ping CLI configuration value` | `config_profile` or `config_grant_type` depends on a value that is only known after apply | Use a literal value, an input variable or the matching environment variable |
| `Missing region code` | No `region_code` attribute or `PINGONE_REGION_CODE` environment variable | Set the region for the PingOne tenant |
| `no Ping CLI authentication session was found in the OS keychain` | No session for the profile (and grant type, when `config_grant_type` is set), or the session is under a different profile than expected | Run `pingcli pingone auth login`; note the session is stored under the profile active at login time. If `config_grant_type` is set, check it matches the grant type used to log in |
| `The actor attempting to perform the request is not authorized` (`ACCESS_FAILED`) | The Ping CLI login's application is not a Worker application, or it has no PingOne admin role for the environment being managed. With Ping CLI authentication the provider retries it for up to 30 seconds, in case a resource was created during the run, and then returns it | Use a Worker application for the Ping CLI login and assign it the needed admin roles (user roles alone are not enough for other application types), or log in with a different profile or grant type (`config_grant_type`) |
| `has a format this provider does not recognise` | The keychain entry has neither an access token nor Ping CLI's split-session marker, for example because a newer Ping CLI changed the layout | Update the provider, or authenticate with client credentials or an access token instead |
| `has expired` | The stored access token has passed its expiry (with a small safety margin) | Run `pingcli pingone auth refresh` or `pingcli pingone auth login` |
| `failed to read the Ping CLI authentication session from the OS keychain` | The OS keychain is unavailable (for example headless Linux without a Secret Service) | Run where an OS keychain is available; storage types `file_system` and `none` are not read by the provider |
