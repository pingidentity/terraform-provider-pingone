---
page_title: "pingone_application Resource - terraform-provider-pingone"
subcategory: "SSO"
description: |-
  Resource to create and manage a PingOne application (SAML, OpenID Connect, External Link) in an environment.
---

# pingone_application (Resource)

Resource to create and manage a PingOne application (SAML, OpenID Connect, External Link) in an environment.

## Example Usage - Single Page Application (SPA)

```terraform
resource "pingone_application" "my_awesome_spa" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Single Page App"
  enabled        = true

  oidc_options = {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://my-website.com"]
  }
}

resource "time_rotating" "my_awesome_spa_secret_rotation" {
  rotation_days = 30
}

resource "pingone_application_secret" "my_awesome_spa" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  regenerate_trigger_values = {
    "rotation_rfc3339" : time_rotating.my_awesome_spa_secret_rotation.rotation_rfc3339,
  }
}
```

## Example Usage - Web Application

```terraform
resource "pingone_application" "my_awesome_web_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Web App"
  enabled        = true

  oidc_options = {
    type                        = "WEB_APP"
    grant_types                 = ["AUTHORIZATION_CODE", "REFRESH_TOKEN"]
    response_types              = ["CODE"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
    redirect_uris               = ["https://my-website.com"]
  }
}

resource "time_rotating" "my_awesome_web_app_secret_rotation" {
  rotation_days = 30
}

resource "pingone_application_secret" "my_awesome_web_app" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_web_app.id

  regenerate_trigger_values = {
    "rotation_rfc3339" : time_rotating.my_awesome_web_app_secret_rotation.rotation_rfc3339,
  }
}
```

## Example Usage - SAML Application

```terraform
resource "pingone_key" "my_awesome_key" {
  environment_id = pingone_environment.my_environment.id

  name                = "Example Signing Key"
  algorithm           = "RSA"
  key_length          = 4096
  signature_algorithm = "SHA512withRSA"
  subject_dn          = "CN=Example Signing Key, OU=BX Retail, O=BX Retail, L=, ST=, C=US"
  usage_type          = "SIGNING"
  validity_period     = 365
}

resource "pingone_application" "my_awesome_saml_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome SAML App"
  enabled        = true

  saml_options = {
    acs_urls           = ["https://my-saas-app.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:localhost"

    idp_signing_key = {
      key_id    = pingone_key.my_awesome_key.id
      algorithm = pingone_key.my_awesome_key.signature_algorithm
    }

    sp_verification = {
      certificate_ids      = [var.sp_verification_certificate_id]
      authn_request_signed = true
    }
  }
}
```

## Example Usage - Native Application (Mobile)

```terraform
resource "pingone_application" "my_awesome_native_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Native Mobile App"
  enabled        = true

  oidc_options = {
    type                        = "NATIVE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris = [
      "https://demo.bxretail.org/app/callback",
      "org.bxretail.app://callback"
    ]

    mobile_app = {
      bundle_id           = var.apple_bundle_id
      package_name        = var.android_package_name
      huawei_app_id       = var.huawei_app_id
      huawei_package_name = var.huawei_package_name

      universal_app_link = "https://demo.bxretail.org"

      passcode_refresh_seconds = 30

      integrity_detection = {
        enabled = true

        cache_duration = {
          amount = 24
          units  = "HOURS"
        }

        google_play = {
          verification_type = "INTERNAL"
          decryption_key    = var.google_play_integrity_api_decryption_key
          verification_key  = var.google_play_integrity_api_verification_key
        }
      }
    }
  }
}

resource "time_rotating" "my_awesome_native_app_secret_rotation" {
  rotation_days = 30
}

resource "pingone_application_secret" "my_awesome_native_app" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_native_app.id

  regenerate_trigger_values = {
    "rotation_rfc3339" : time_rotating.my_awesome_native_app_secret_rotation.rotation_rfc3339,
  }
}
```

## Example Usage - Worker Application

```terraform
resource "pingone_application" "my_awesome_worker_app" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Worker App"
  enabled        = true

  oidc_options = {
    type                        = "WORKER"
    grant_types                 = ["CLIENT_CREDENTIALS"]
    token_endpoint_authn_method = "CLIENT_SECRET_BASIC"
  }
}

resource "time_rotating" "my_awesome_worker_app_secret_rotation" {
  rotation_days = 30
}

resource "pingone_application_secret" "my_awesome_worker_app" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_worker_app.id

  regenerate_trigger_values = {
    "rotation_rfc3339" : time_rotating.my_awesome_worker_app_secret_rotation.rotation_rfc3339,
  }
}
```

## Example Usage - External Link

```terraform
resource "pingone_application" "my_awesome_external_link" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome External Link"
  enabled        = true

  external_link_options = {
    home_page_url = "https://demo.bxretail.org/"
  }
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The PingOne resource ID of the environment to create and manage the application in.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `name` (String) A string that specifies the name of the application.

### Optional

- `access_control_group_options` (Attributes) A single object that specifies group access control settings. (see [below for nested schema](#nestedatt--access_control_group_options))
- `access_control_role_type` (String) A string that specifies the user role required to access the application.  A user is an admin user if the user has one or more admin roles assigned, such as `Organization Admin`, `Environment Admin`, `Identity Data Admin`, or `Client Application Developer`.  Options are `ADMIN_USERS_ONLY`.
- `description` (String) A string that specifies the description of the application.
- `enabled` (Boolean) A boolean that specifies whether the application is enabled in the environment.  Defaults to `false`.
- `external_link_options` (Attributes) A single object that specifies External link application specific settings.  At least one of the following must be defined: `external_link_options`, `oidc_options`, `saml_options`. (see [below for nested schema](#nestedatt--external_link_options))
- `hidden_from_app_portal` (Boolean) A boolean to specify whether the application is hidden in the application portal despite the configured group access policy.  Defaults to `false`.
- `icon` (Attributes) A single object that specifies settings for the application icon. (see [below for nested schema](#nestedatt--icon))
- `login_page_url` (String) A string that specifies the custom login page URL for the application. If you set the `login_page_url` property for applications in an environment that sets a custom domain, the URL should include the top-level domain and at least one additional domain level. **Warning** To avoid issues with third-party cookies in some browsers, a custom domain must be used, giving your PingOne environment the same parent domain as your authentication application. For more information about custom domains, see Custom domains.  The provided URL is expected to use the `https://` schema.  The `http` schema is permitted where the host is `localhost` or `127.0.0.1`.
- `oidc_options` (Attributes) A single object that specifies OIDC/OAuth application specific settings.  At least one of the following must be defined: `external_link_options`, `oidc_options`, `saml_options`. (see [below for nested schema](#nestedatt--oidc_options))
- `saml_options` (Attributes) A single object that specifies SAML application specific settings.  At least one of the following must be defined: `external_link_options`, `oidc_options`, `saml_options`. (see [below for nested schema](#nestedatt--saml_options))
- `tags` (Set of String) An array of strings that specifies the list of labels associated with the application.  Options are `PING_FED_CONNECTION_INTEGRATION`.  Conflicts with `external_link_options`, `saml_options`.

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--access_control_group_options"></a>
### Nested Schema for `access_control_group_options`

Required:

- `groups` (Set of String) A set that specifies the group IDs for the groups the actor must belong to for access to the application.  Values must be valid PingOne Resource IDs.
- `type` (String) A string that specifies the group type required to access the application.  Options are `ALL_GROUPS` (the actor must belong to all groups listed in the `groups` property), `ANY_GROUP` (the actor must belong to at least one group listed in the `groups` property).


<a id="nestedatt--external_link_options"></a>
### Nested Schema for `external_link_options`

Required:

- `home_page_url` (String) A string that specifies the custom home page URL for the application.  Both `http://` and `https://` URLs are permitted.


<a id="nestedatt--icon"></a>
### Nested Schema for `icon`

Required:

- `href` (String) A string that specifies the URL for the application icon.  Both `http://` and `https://` are permitted.
- `id` (String) A string that specifies the ID for the application icon.  Must be a valid PingOne Resource ID.


<a id="nestedatt--oidc_options"></a>
### Nested Schema for `oidc_options`

Required:

- `grant_types` (Set of String) A list that specifies the grant type for the authorization request.  Options are `AUTHORIZATION_CODE`, `CLIENT_CREDENTIALS`, `IMPLICIT`, `REFRESH_TOKEN`.
- `token_endpoint_authn_method` (String) A string that specifies the client authentication methods supported by the token endpoint.  Options are `CLIENT_SECRET_BASIC`, `CLIENT_SECRET_JWT`, `CLIENT_SECRET_POST`, `NONE`, `PRIVATE_KEY_JWT`.  When `PRIVATE_KEY_JWT` is configured, either `jwks` or `jwks_url` must also be configured.
- `type` (String) A string that specifies the type associated with the application.  Options are `CUSTOM_APP`, `NATIVE_APP`, `SERVICE`, `SINGLE_PAGE_APP`, `WEB_APP`, `WORKER`.  This field is immutable and will trigger a replace plan if changed.

Optional:

- `additional_refresh_token_replay_protection_enabled` (Boolean) A boolean that, when set to `true` (the default), if you attempt to reuse the refresh token, the authorization server immediately revokes the reused refresh token, as well as all descendant tokens. Setting this to null equates to a `false` setting.  Defaults to `true`.
- `allow_wildcards_in_redirect_uris` (Boolean) A boolean to specify whether wildcards are allowed in redirect URIs. For more information, see [Wildcards in Redirect URIs](https://docs.pingidentity.com/csh?context=p1_c_wildcard_redirect_uri).  Defaults to `false`.
- `certificate_based_authentication` (Attributes) A single object that specifies Certificate based authentication settings. This parameter block can only be set where the application's `type` parameter is set to `NATIVE_APP`. (see [below for nested schema](#nestedatt--oidc_options--certificate_based_authentication))
- `cors_settings` (Attributes) A single object that allows customization of how the Authorization and Authentication APIs interact with CORS requests that reference the application. If omitted, the application allows CORS requests from any origin except for operations that expose sensitive information (e.g. `/as/authorize` and `/as/token`).  This is legacy behavior, and it is recommended that applications migrate to include specific CORS settings. (see [below for nested schema](#nestedatt--oidc_options--cors_settings))
- `home_page_url` (String) A string that specifies the custom home page URL for the application.  The provided URL is expected to use the `https://` schema.  The `http` schema is permitted where the host is `localhost` or `127.0.0.1`.
- `initiate_login_uri` (String) A string that specifies the URI to use for third-parties to begin the sign-on process for the application. If specified, PingOne redirects users to this URI to initiate SSO to PingOne. The application is responsible for implementing the relevant OIDC flow when the initiate login URI is requested. This property is required if you want the application to appear in the PingOne Application Portal. See the OIDC specification section of [Initiating Login from a Third Party](https://openid.net/specs/openid-connect-core-1_0.html#ThirdPartyInitiatedLogin) for more information.  The provided URL is expected to use the `https://` schema.  The `http` schema is permitted where the host is `localhost` or `127.0.0.1`.
- `jwks` (String) A string that specifies a JWKS string that validates the signature of signed JWTs for applications that use the `PRIVATE_KEY_JWT` option for the `token_endpoint_authn_method`. This property is required when `token_endpoint_authn_method` is `PRIVATE_KEY_JWT` and the `jwks_url` property is empty. For more information, see [Create a private_key_jwt JWKS string](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-private_key_jwt-jwks-string). This property is also required if the optional `request` property JWT on the authorize endpoint is signed using the RS256 (or RS384, RS512) signing algorithm and the `jwks_url` property is empty. For more infornmation about signing the `request` property JWT, see [Create a request property JWT](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-request-property-jwt).  Conflicts with `jwks_url`.
- `jwks_url` (String) A string that specifies a URL (supports `https://` only) that provides access to a JWKS string that validates the signature of signed JWTs for applications that use the `PRIVATE_KEY_JWT` option for the `token_endpoint_authn_method`. This property is required when `token_endpoint_authn_method` is `PRIVATE_KEY_JWT` and the `jwks` property is empty. For more information, see [Create a private_key_jwt JWKS string](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-private_key_jwt-jwks-string). This property is also required if the optional `request` property JWT on the authorize endpoint is signed using the RS256 (or RS384, RS512) signing algorithm and the `jwks` property is empty. For more infornmation about signing the `request` property JWT, see [Create a request property JWT](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-request-property-jwt).  Conflicts with `jwks`.
- `mobile_app` (Attributes) A single object that specifies Mobile application integration settings for `NATIVE_APP` type applications. (see [below for nested schema](#nestedatt--oidc_options--mobile_app))
- `par_requirement` (String) A string that specifies whether pushed authorization requests (PAR) are required.  Options are `OPTIONAL`, `REQUIRED`.  Defaults to `OPTIONAL`.
- `par_timeout` (Number) An integer that specifies the pushed authorization request (PAR) timeout in seconds.  Valid values are between `1` and `600`.  Defaults to `60`.
- `pkce_enforcement` (String) A string that specifies how `PKCE` request parameters are handled on the authorize request.  Options are `OPTIONAL`, `REQUIRED`, `S256_REQUIRED`.  Defaults to `OPTIONAL`.
- `post_logout_redirect_uris` (Set of String) A list of strings that specifies the URLs that the browser can be redirected to after logout.  The provided URLs are expected to use the `https://`, `http://` schema, or a custom mobile native schema (e.g., `org.bxretail.app://logout`).
- `redirect_uris` (Set of String) A list of strings that specifies the allowed callback URIs for the authentication response.    The provided URLs are expected to use the `https://` schema, or a custom mobile native schema (e.g., `org.bxretail.app://callback`).  The `http` schema is only permitted where the host is `localhost` or `127.0.0.1`.
- `refresh_token_duration` (Number) An integer that specifies the lifetime in seconds of the refresh token. Valid values are between `60` and `2147483647`. If the `refresh_token_rolling_duration` property is specified for the application, then this property value must be less than or equal to the value of `refresh_token_rolling_duration`. After this property is set, the value cannot be nullified - this will force recreation of the resource. This value is used to generate the value for the exp claim when minting a new refresh token.  Defaults to `2592000`.
- `refresh_token_rolling_duration` (Number) An integer that specifies the number of seconds a refresh token can be exchanged before re-authentication is required. Valid values are between `60` and `2147483647`. After this property is set, the value cannot be nullified - this will force recreation of the resource. This value is used to generate the value for the exp claim when minting a new refresh token.  Defaults to `15552000`.
- `refresh_token_rolling_grace_period_duration` (Number) The number of seconds that a refresh token may be reused after having been exchanged for a new set of tokens. This is useful in the case of network errors on the client. Valid values are between `0` and `86400` seconds. `Null` is treated the same as `0`.
- `require_signed_request_object` (Boolean) A boolean that indicates that the Java Web Token (JWT) for the [request query](https://openid.net/specs/openid-connect-core-1_0.html#RequestObject) parameter is required to be signed. If `false` or null, a signed request object is not required. Both `support_unsigned_request_object` and this property cannot be set to `true`.  Defaults to `false`.
- `response_types` (Set of String) A list that specifies the code or token type returned by an authorization request.  Options are `CODE`, `ID_TOKEN`, `TOKEN`.  Note that `CODE` cannot be used in an authorization request with `TOKEN` or `ID_TOKEN` because PingOne does not currently support OIDC hybrid flows.
- `support_unsigned_request_object` (Boolean) A boolean that specifies whether the request query parameter JWT is allowed to be unsigned. If `false` or null, an unsigned request object is not allowed.  Defaults to `false`.
- `target_link_uri` (String) The URI for the application. If specified, PingOne will redirect application users to this URI after a user is authenticated. In the PingOne admin console, this becomes the value of the `target_link_uri` parameter used for the Initiate Single Sign-On URL field.  Both `http://` and `https://` URLs are permitted as well as custom mobile native schema (e.g., `org.bxretail.app://target`).

<a id="nestedatt--oidc_options--certificate_based_authentication"></a>
### Nested Schema for `oidc_options.certificate_based_authentication`

Required:

- `key_id` (String) A string that represents a PingOne ID for the issuance certificate key.  The key must be of type `ISSUANCE`.  Must be a valid PingOne Resource ID.


<a id="nestedatt--oidc_options--cors_settings"></a>
### Nested Schema for `oidc_options.cors_settings`

Required:

- `behavior` (String) A string that specifies the behavior of how Authorization and Authentication APIs interact with CORS requests that reference the application.  Options are `ALLOW_NO_ORIGINS` (rejects all CORS requests), `ALLOW_SPECIFIC_ORIGINS` (rejects all CORS requests except those listed in `origins`).

Optional:

- `origins` (Set of String) A set of strings that represent the origins from which CORS requests to the Authorization and Authentication APIs are allowed.  Each value must be a `http` or `https` URL without a path.  The host may be a domain name (including `localhost`), or an IPv4 address.  Subdomains may use the wildcard (`*`) to match any string.  Must be non-empty when `behavior` is `ALLOW_SPECIFIC_ORIGINS` and must be omitted or empty when `behavior` is `ALLOW_NO_ORIGINS`.  Limited to 20 values.


<a id="nestedatt--oidc_options--mobile_app"></a>
### Nested Schema for `oidc_options.mobile_app`

Optional:

- `bundle_id` (String) A string that specifies the bundle associated with the application, for push notifications in native apps. The value of the `bundle_id` property is unique per environment, and once defined, is immutable.  This field is immutable and will trigger a replace plan if changed.
- `huawei_app_id` (String) The unique identifier for the app on the device and in the Huawei Mobile Service AppGallery. The value of this property is unique per environment, and once defined, is immutable.  Required with `huawei_package_name`.  This field is immutable and will trigger a replace plan if changed.
- `huawei_package_name` (String) The package name associated with the application, for push notifications in native apps. The value of this property is unique per environment, and once defined, is immutable.  Required with `huawei_app_id`.  This field is immutable and will trigger a replace plan if changed.
- `integrity_detection` (Attributes) A single object that specifies mobile application integrity detection settings. (see [below for nested schema](#nestedatt--oidc_options--mobile_app--integrity_detection))
- `package_name` (String) A string that specifies the package name associated with the application, for push notifications in native apps. The value of the `package_name` property is unique per environment, and once defined, is immutable.  This field is immutable and will trigger a replace plan if changed.
- `passcode_refresh_seconds` (Number) The amount of time a passcode should be displayed before being replaced with a new passcode - must be between `30` and `60` seconds.  Defaults to `30`.
- `universal_app_link` (String) A string that specifies a URI prefix that enables direct triggering of the mobile application when scanning a QR code. The URI prefix can be set to a universal link with a valid value (which can be a URL address that starts with `HTTP://` or `HTTPS://`, such as `https://www.bxretail.org`), or an app schema, which is just a string and requires no special validation.

<a id="nestedatt--oidc_options--mobile_app--integrity_detection"></a>
### Nested Schema for `oidc_options.mobile_app.integrity_detection`

Optional:

- `cache_duration` (Attributes) A single object that specifies settings for the caching duration of successful integrity detection calls.  Every attestation request entails a certain time tradeoff. You can choose to cache successful integrity detection calls for a predefined duration, between a minimum of 1 minute and a maximum of 48 hours. If integrity detection is ENABLED, the cache duration must be set. (see [below for nested schema](#nestedatt--oidc_options--mobile_app--integrity_detection--cache_duration))
- `enabled` (Boolean) A boolean that specifies whether device integrity detection takes place on mobile devices.  Defaults to `false`.
- `excluded_platforms` (Set of String) You can enable device integrity checking separately for Android and iOS by setting `enabled` to `true` and then using `excluded_platforms` to specify the OS where you do not want to use device integrity checking. The values to use are `GOOGLE` and `IOS` (all upper case). Note that this is implemented as an array even though currently you can only include a single value.  If `GOOGLE` is not included in this list, the `google_play` attribute block must be configured.
- `google_play` (Attributes) A single object that describes Google Play Integrity API credential settings for Android device integrity detection.  Required when `excluded_platforms` is unset or does not include `GOOGLE`. (see [below for nested schema](#nestedatt--oidc_options--mobile_app--integrity_detection--google_play))

<a id="nestedatt--oidc_options--mobile_app--integrity_detection--cache_duration"></a>
### Nested Schema for `oidc_options.mobile_app.integrity_detection.google_play`

Required:

- `amount` (Number) An integer that specifies the number of minutes or hours that specify the duration between successful integrity detection calls.

Optional:

- `units` (String) A string that specifies the time units of the cache `amount` parameter.  Options are `HOURS`, `MINUTES`.  Defaults to `MINUTES`.


<a id="nestedatt--oidc_options--mobile_app--integrity_detection--google_play"></a>
### Nested Schema for `oidc_options.mobile_app.integrity_detection.google_play`

Required:

- `verification_type` (String) The type of verification that should be used.  Options are `GOOGLE`, `INTERNAL`.  Using internal verification will not count against your Google API call quota. The value you select for this attribute determines what other parameters you must provide. When set to `GOOGLE`, you must provide `service_account_credentials_json`. When set to `INTERNAL`, you must provide both `decryption_key` and `verification_key`.

Optional:

- `decryption_key` (String, Sensitive) Play Integrity verdict decryption key from your Google Play Services account. This parameter must be provided if you have set `verification_type` to `INTERNAL`.  Conflicts with `service_account_credentials_json`.
- `service_account_credentials_json` (String, Sensitive) Contents of the JSON file that represents your Service Account Credentials. This parameter must be provided if you have set `verification_type` to `GOOGLE`.  Conflicts with `decryption_key`, `verification_key`.
- `verification_key` (String, Sensitive) Play Integrity verdict signature verification key from your Google Play Services account. This parameter must be provided if you have set `verification_type` to `INTERNAL`.  Conflicts with `service_account_credentials_json`.





<a id="nestedatt--saml_options"></a>
### Nested Schema for `saml_options`

Required:

- `acs_urls` (Set of String) A list of string that specifies the Assertion Consumer Service URLs. The first URL in the list is used as default (there must be at least one URL).
- `assertion_duration` (Number) An integer that specifies the assertion validity duration in seconds.
- `idp_signing_key` (Attributes) SAML application assertion/response signing key settings.  Use with `assertion_signed_enabled` to enable assertion signing and/or `response_is_signed` to enable response signing.  It's highly recommended, and best practice, to define signing key settings for the configured SAML application.  However if this property is omitted, the default signing certificate for the environment is used.  This parameter will become a required field in the next major release of the provider. (see [below for nested schema](#nestedatt--saml_options--idp_signing_key))
- `sp_entity_id` (String) A string that specifies the service provider entity ID used to lookup the application. This is a required property and is unique within the environment.

Optional:

- `assertion_signed_enabled` (Boolean) A boolean that specifies whether the SAML assertion itself should be signed.  Defaults to `true`.
- `cors_settings` (Attributes) A single object that allows customization of how the Authorization and Authentication APIs interact with CORS requests that reference the application. If omitted, the application allows CORS requests from any origin except for operations that expose sensitive information (e.g. `/as/authorize` and `/as/token`).  This is legacy behavior, and it is recommended that applications migrate to include specific CORS settings. (see [below for nested schema](#nestedatt--saml_options--cors_settings))
- `default_target_url` (String) A string that specfies a default URL used as the `RelayState` parameter by the IdP to deep link into the application after authentication. This value can be overridden by the `applicationUrl` query parameter for [GET Identity Provider Initiated SSO](https://apidocs.pingidentity.com/pingone/platform/v1/api/#get-identity-provider-initiated-sso). Although both of these parameters are generally URLs, because they are used as deep links, this is not enforced. If neither `defaultTargetUrl` nor `applicationUrl` is specified during a SAML authentication flow, no `RelayState` value is supplied to the application. The `defaultTargetUrl` (or the `applicationUrl`) value is passed to the SAML application’s ACS URL as a separate `RelayState` key value (not within the SAMLResponse key value).
- `enable_requested_authn_context` (Boolean) A boolean that specifies whether `requestedAuthnContext` is taken into account in policy decision-making.
- `home_page_url` (String) A string that specifies the custom home page URL for the application.
- `nameid_format` (String) A string that specifies the format of the Subject NameID attibute in the SAML assertion.
- `response_is_signed` (Boolean) A boolean that specifies whether the SAML assertion response itself should be signed.  Defaults to `false`.
- `slo_binding` (String) A string that specifies the binding protocol to be used for the logout response.  Options are `HTTP_POST`, `HTTP_REDIRECT`.  Existing configurations with no data default to `HTTP_POST`.  Defaults to `HTTP_POST`.
- `slo_endpoint` (String) A string that specifies the logout endpoint URL. This is an optional property. However, if a logout endpoint URL is not defined, logout actions result in an error.
- `slo_response_endpoint` (String) A string that specifies the endpoint URL to submit the logout response. If a value is not provided, the `slo_endpoint` property value is used to submit SLO response.
- `slo_window` (Number) An integer that defines how long (hours) PingOne can exchange logout messages with the application, specifically a logout request from the application, since the initial request.  The minimum value is `0` hour and the maximum is `24` hours.
- `sp_verification` (Attributes) A single object item that specifies SP signature verification settings. (see [below for nested schema](#nestedatt--saml_options--sp_verification))
- `type` (String) A string that specifies the type associated with the application.  Options are `CUSTOM_APP`, `WEB_APP`.  Defaults to `WEB_APP`.  This field is immutable and will trigger a replace plan if changed.

<a id="nestedatt--saml_options--idp_signing_key"></a>
### Nested Schema for `saml_options.idp_signing_key`

Required:

- `algorithm` (String) Specifies the signature algorithm of the key. For RSA keys, options are `SHA256withRSA`, `SHA384withRSA` and `SHA512withRSA`. For elliptical curve (EC) keys, options are `SHA256withECDSA`, `SHA384withECDSA` and `SHA512withECDSA`.
- `key_id` (String) An ID for the certificate key pair to be used by the identity provider to sign assertions and responses.  Must be a valid PingOne resource ID.


<a id="nestedatt--saml_options--cors_settings"></a>
### Nested Schema for `saml_options.cors_settings`

Required:

- `behavior` (String) A string that specifies the behavior of how Authorization and Authentication APIs interact with CORS requests that reference the application.  Options are `ALLOW_NO_ORIGINS` (rejects all CORS requests), `ALLOW_SPECIFIC_ORIGINS` (rejects all CORS requests except those listed in `origins`).

Optional:

- `origins` (Set of String) A set of strings that represent the origins from which CORS requests to the Authorization and Authentication APIs are allowed.  Each value must be a `http` or `https` URL without a path.  The host may be a domain name (including `localhost`), or an IPv4 address.  Subdomains may use the wildcard (`*`) to match any string.  Must be non-empty when `behavior` is `ALLOW_SPECIFIC_ORIGINS` and must be omitted or empty when `behavior` is `ALLOW_NO_ORIGINS`.  Limited to 20 values.


<a id="nestedatt--saml_options--sp_verification"></a>
### Nested Schema for `saml_options.sp_verification`

Required:

- `certificate_ids` (Set of String) A list that specifies the certificate IDs used to verify the service provider signature.  Values must be valid PingOne resource IDs.

Optional:

- `authn_request_signed` (Boolean) A boolean that specifies whether the Authn Request signing should be enforced.  Defaults to `false`.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_application.example <environment_id>/<application_id>
```
