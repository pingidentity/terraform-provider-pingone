---
layout: ""
page_title: "Using PingCLI with the Provider"
description: |-
  Configure the PingOne Terraform provider to authenticate using PingCLI profiles and stored tokens.
---

# Using PingCLI with the Provider

The PingOne Terraform provider can authenticate using credentials and tokens managed by [PingCLI](https://github.com/pingidentity/pingcli). This allows you to log in once with PingCLI and have Terraform automatically discover and use the stored token or client credentials from your PingCLI profile.

## Overview

- Run `pingcli login` with your preferred grant type.
- Ensure your PingCLI config’s `authentication.type` matches the grant.
- Point the Terraform provider to the PingCLI config file and profile via `config_path` and `config_profile` (or the env vars `PINGCLI_CONFIG` and `PINGCLI_PROFILE`).

## Step 1: Log in with PingCLI

Choose a grant type and profile name (e.g., `dev`). This writes a token to secure storage (macOS Keychain by default) and optionally to `~/.pingcli/credentials` if file storage is enabled.

```sh
# Authorization Code flow (interactive)
pingcli login --type authorization_code --profile dev

# Device Code flow (interactive, alternative to auth code)
pingcli login --type device_code --profile dev

# Client Credentials flow (non-interactive)
pingcli login --type client_credentials --profile dev
```

## Step 2: Verify your PingCLI configuration

The provider reads the PingCLI config to resolve region, environment ID, auth type, and (for some flows) client credentials. A minimal example for the `authorization_code` grant type:

```yaml
activeProfile: dev

dev:
  service:
    pingOne:
      regionCode: NA
      authentication:
        type: authorization_code
        environmentID: 11111111-1111-1111-1111-111111111111
        authorizationCode:
          clientID: 22222222-2222-2222-2222-222222222222
          # Optional fields if you prefer custom redirect params:
          # redirectURIPath: /callback
          # redirectURIPort: "3000"
```

Other grant type snippets:

Device Code:
```yaml
dev:
  service:
    pingOne:
      regionCode: NA
      authentication:
        type: device_code
        environmentID: 11111111-1111-1111-1111-111111111111
        deviceCode:
          clientID: 22222222-2222-2222-2222-222222222222
```

Client Credentials:
```yaml
cc:
  service:
    pingOne:
      regionCode: NA
      authentication:
        type: client_credentials
        environmentID: 11111111-1111-1111-1111-111111111111
        clientCredentials:
          clientID: 22222222-2222-2222-2222-222222222222
          clientSecret: ${CLIENT_SECRET}
```

Worker (alias of client credentials):
```yaml
worker:
  service:
    pingOne:
      regionCode: NA
      authentication:
        type: worker
        worker:
          environmentID: 11111111-1111-1111-1111-111111111111
          clientID: 22222222-2222-2222-2222-222222222222
          clientSecret: ${CLIENT_SECRET}
```

## Step 3: Point the provider at your PingCLI profile

You can configure this in Terraform:

```hcl
provider "pingone" {
  config_path    = "~/.pingcli/config.yaml"
  config_profile = "dev"
}
```

Or via environment variables:

```sh
export PINGCLI_CONFIG="~/.pingcli/config.yaml"
export PINGCLI_PROFILE="dev"
```

## Behavior and Fallbacks

When `config_path` is set, the provider resolves credentials in this order:

1. Use a valid stored token from PingCLI (Keychain or `~/.pingcli/credentials`).
2. If no stored token is found and `api_access_token` is set in the provider, use that token.
3. If the PingCLI profile contains client credentials (`client_credentials` or `worker`) and an environment ID, use those to obtain a token.
4. If none of the above apply, the provider returns an error and prompts you to run `pingcli login`.

Region resolution:
- If using a stored token, the region is taken from the PingCLI profile (`regionCode`).
- If using `api_access_token` or provider-supplied credentials, `region_code` in the provider (or `PINGONE_REGION_CODE`) can be used to select endpoints. When both are present, the provider value takes precedence.

## FAQ

- Do I need to run `pingcli login` for client credentials? For non-interactive flows, you can either run `pingcli login --type client_credentials` or specify the client credentials directly in the PingCLI config. The provider will use them if no stored token is found.
- Where does PingCLI store tokens? On macOS, tokens are persisted in the Keychain. If file storage is enabled, it also writes to `~/.pingcli/credentials`.
- Can I still use `client_id`, `client_secret`, and `environment_id` in the provider? Yes, but not together with `config_path`. For pingcli-based auth, prefer `config_path` and `config_profile`.
