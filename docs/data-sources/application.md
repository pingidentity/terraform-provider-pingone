---
page_title: "pingone_application Data Source - terraform-provider-pingone"
subcategory: "SSO"
description: |-
  Data source to retrieve a PingOne application in an environment by ID or by name.
---

# pingone_application (Data Source)

Data source to retrieve a PingOne application in an environment by ID or by name.

## Example Usage

```terraform
data "pingone_application" "example_by_name" {
  environment_id = var.environment_id
  name           = "foo"
}

data "pingone_application" "example_by_id" {
  environment_id = var.environment_id
  application_id = var.application_id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) PingOne environment identifier (UUID) in which the application exists.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.

### Optional

- `application_id` (String) The identifier (UUID) of the application.  Exactly one of the following must be defined: `application_id`, `name`.  Must be a valid PingOne resource ID.
- `name` (String) The name of the application.  Exactly one of the following must be defined: `application_id`, `name`.

### Read-Only

- `access_control_group_options` (Attributes) Group access control settings. (see [below for nested schema](#nestedatt--access_control_group_options))
- `access_control_role_type` (String) The user role required to access the application.  A user is an admin user if the user has one or more of the following roles: `Organization Admin`, `Environment Admin`, `Identity Data Admin`, or `Client Application Developer`.
- `description` (String) A string that specifies the description of the application.
- `enabled` (Boolean) A boolean that specifies whether the application is enabled in the environment.
- `external_link_options` (Attributes) External link application specific settings. (see [below for nested schema](#nestedatt--external_link_options))
- `hidden_from_app_portal` (Boolean) A boolean to specify whether the application is hidden in the application portal despite the configured group access policy.
- `icon` (Attributes) The HREF and the ID for the application icon. (see [below for nested schema](#nestedatt--icon))
- `id` (String) The ID of this resource.
- `login_page_url` (String) A string that specifies the custom login page URL for the application.
- `oidc_options` (Attributes) OIDC/OAuth application specific settings. (see [below for nested schema](#nestedatt--oidc_options))
- `saml_options` (Attributes) SAML application specific settings. (see [below for nested schema](#nestedatt--saml_options))
- `tags` (Set of String) An array that specifies the list of labels associated with the application.
- `wsfed_options` (Attributes) A single object that specifies WS-Fed application specific settings. (see [below for nested schema](#nestedatt--wsfed_options))

<a id="nestedatt--access_control_group_options"></a>
### Nested Schema for `access_control_group_options`

Read-Only:

- `groups` (Set of String) A set that specifies the group IDs for the groups the actor must belong to for access to the application.
- `type` (String) A string that specifies the group type required to access the application.


<a id="nestedatt--external_link_options"></a>
### Nested Schema for `external_link_options`

Read-Only:

- `home_page_url` (String) The custom home page URL for the application.  Both `http://` and `https://` URLs are permitted.


<a id="nestedatt--icon"></a>
### Nested Schema for `icon`

Read-Only:

- `href` (String) The HREF for the application icon.
- `id` (String) The ID for the application icon.


<a id="nestedatt--oidc_options"></a>
### Nested Schema for `oidc_options`

Read-Only:

- `additional_refresh_token_replay_protection_enabled` (Boolean) A boolean that, when set to `true`, if you attempt to reuse the refresh token, the authorization server immediately revokes the reused refresh token, as well as all descendant tokens.  Defaults to `true`.
- `allow_wildcard_in_redirect_uris` (Boolean) A boolean to specify whether wildcards are allowed in redirect URIs. For more information, see [Wildcards in Redirect URIs](https://docs.pingidentity.com/csh?context=p1_c_wildcard_redirect_uri).
- `certificate_based_authentication` (Attributes) Certificate based authentication settings. (see [below for nested schema](#nestedatt--oidc_options--certificate_based_authentication))
- `client_id` (String) A string that specifies the application ID used to authenticate to the authorization server.
- `cors_settings` (Attributes) A single object that allows customization of how the Authorization and Authentication APIs interact with CORS requests that reference the application. If omitted, the application allows CORS requests from any origin except for operations that expose sensitive information (e.g. `/as/authorize` and `/as/token`).  This is legacy behavior, and it is recommended that applications migrate to include specific CORS settings. (see [below for nested schema](#nestedatt--oidc_options--cors_settings))
- `device_custom_verification_uri` (String) A string that specifies an optional custom verification URI that is returned for the `/device_authorization` endpoint.
- `device_path_id` (String) A string that describes a unique identifier within an environment for a device authorization grant flow to provide a short identifier to the application. This property is ignored when the `device_custom_verification_uri` property is configured.
- `device_polling_interval` (Number) An integer that specifies the frequency (in seconds) for the client to poll the `/as/token` endpoint.
- `device_timeout` (Number) An integer that specifies the length of time (in seconds) that the `userCode` and `deviceCode` returned by the `/device_authorization` endpoint are valid.
- `grant_types` (Set of String) A list that specifies the grant type for the authorization request.
- `home_page_url` (String) The custom home page URL for the application.  The provided URL is expected to use the `https://` schema.  The `http` schema is permitted where the host is `localhost` or `127.0.0.1`.
- `idp_signoff` (Boolean) A boolean flag to allow signoff without access to the session token cookie.
- `initiate_login_uri` (String) A string that specifies the URI to use for third-parties to begin the sign-on process for the application.
- `jwks` (String) A string that specifies a JWKS string that validates the signature of signed JWTs for applications that use the `PRIVATE_KEY_JWT` option for the `token_endpoint_auth_method`. This property is required when `token_endpoint_auth_method` is `PRIVATE_KEY_JWT` and the `jwks_url` property is empty. For more information, see [Create a private_key_jwt JWKS string](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-private_key_jwt-jwks-string). This property is also required if the optional `request` property JWT on the authorize endpoint is signed using the RS256 (or RS384, RS512) signing algorithm and the `jwks_url` property is empty. For more infornmation about signing the `request` property JWT, see [Create a request property JWT](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-request-property-jwt).
- `jwks_url` (String) A string that specifies a URL (supports `https://` only) that provides access to a JWKS string that validates the signature of signed JWTs for applications that use the `PRIVATE_KEY_JWT` option for the `token_endpoint_auth_method`. This property is required when `token_endpoint_auth_method` is `PRIVATE_KEY_JWT` and the `jwks` property is empty. For more information, see [Create a private_key_jwt JWKS string](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-private_key_jwt-jwks-string). This property is also required if the optional `request` property JWT on the authorize endpoint is signed using the RS256 (or RS384, RS512) signing algorithm and the `jwks` property is empty. For more infornmation about signing the `request` property JWT, see [Create a request property JWT](https://apidocs.pingidentity.com/pingone/platform/v1/api/#create-a-request-property-jwt).
- `mobile_app` (Attributes) Mobile application integration settings. (see [below for nested schema](#nestedatt--oidc_options--mobile_app))
- `par_requirement` (String) A string that specifies whether pushed authorization requests (PAR) are required.
- `par_timeout` (Number) An integer that specifies the pushed authorization request (PAR) timeout in seconds.
- `pkce_enforcement` (String) A string that specifies how `PKCE` request parameters are handled on the authorize request.
- `post_logout_redirect_uris` (Set of String) A list of strings that specifies the URLs that the browser can be redirected to after logout.  The provided URLs are expected to use the `https://`, `http://` schema, or a custom mobile native schema (e.g., `org.bxretail.app://logout`).
- `redirect_uris` (Set of String) A list of strings that specifies the allowed callback URIs for the authentication response.
- `refresh_token_duration` (Number) An integer that specifies the lifetime in seconds of the refresh token.
- `refresh_token_rolling_duration` (Number) An integer that specifies the number of seconds a refresh token can be exchanged before re-authentication is required.
- `refresh_token_rolling_grace_period_duration` (Number) The number of seconds that a refresh token may be reused after having been exchanged for a new set of tokens.
- `require_signed_request_object` (Boolean) A boolean that indicates that the Java Web Token (JWT) for the [request query](https://openid.net/specs/openid-connect-core-1_0.html#RequestObject) parameter is required to be signed. If `false` or null, a signed request object is not required. Both `support_unsigned_request_object` and this property cannot be set to `true`.  Defaults to `false`.
- `response_types` (Set of String) A list that specifies the code or token type returned by an authorization request.
- `support_unsigned_request_object` (Boolean) A boolean that specifies whether the request query parameter JWT is allowed to be unsigned.
- `target_link_uri` (String) The URI for the application.
- `token_endpoint_auth_method` (String) A string that specifies the client authentication methods supported by the token endpoint.
- `type` (String) A string that specifies the type associated with the application.

<a id="nestedatt--oidc_options--certificate_based_authentication"></a>
### Nested Schema for `oidc_options.certificate_based_authentication`

Read-Only:

- `key_id` (String) A string that represents a PingOne ID for the issuance certificate key.


<a id="nestedatt--oidc_options--cors_settings"></a>
### Nested Schema for `oidc_options.cors_settings`

Read-Only:

- `behavior` (String) A string that represents the behavior of how Authorization and Authentication APIs interact with CORS requests that reference the application.  Options are `ALLOW_NO_ORIGINS` (rejects all CORS requests), `ALLOW_SPECIFIC_ORIGINS` (rejects all CORS requests except those listed in `origins`).
- `origins` (Set of String) A set of strings that represent the origins from which CORS requests to the Authorization and Authentication APIs are allowed.  Each value will be a `http` or `https` URL without a path.  The host may be a domain name (including `localhost`), or an IPv4 address.  Subdomains may use the wildcard (`*`) to match any string.  Is expected to be non-empty when `behavior` is `ALLOW_SPECIFIC_ORIGINS` and is expected to be omitted or empty when `behavior` is `ALLOW_NO_ORIGINS`.  Limited to 20 values.


<a id="nestedatt--oidc_options--mobile_app"></a>
### Nested Schema for `oidc_options.mobile_app`

Read-Only:

- `bundle_id` (String) A string that specifies the bundle associated with the application, for push notifications in native apps.
- `huawei_app_id` (String) The unique identifier for the app on the device and in the Huawei Mobile Service AppGallery.
- `huawei_package_name` (String) The package name associated with the application, for push notifications in native apps.
- `integrity_detection` (Attributes) Mobile application integrity detection settings. (see [below for nested schema](#nestedatt--oidc_options--mobile_app--integrity_detection))
- `package_name` (String) A string that specifies the package name associated with the application, for push notifications in native apps.
- `passcode_refresh_seconds` (Number) The amount of time a passcode should be displayed before being replaced with a new passcode.
- `universal_app_link` (String) A string that specifies a URI prefix that enables direct triggering of the mobile application when scanning a QR code.

<a id="nestedatt--oidc_options--mobile_app--integrity_detection"></a>
### Nested Schema for `oidc_options.mobile_app.integrity_detection`

Read-Only:

- `cache_duration` (Attributes) Indicates the caching duration of successful integrity detection calls. (see [below for nested schema](#nestedatt--oidc_options--mobile_app--integrity_detection--cache_duration))
- `enabled` (Boolean) A boolean that specifies whether device integrity detection takes place on mobile devices.
- `excluded_platforms` (Set of String) Indicates OS excluded from device integrity checking.
- `google_play` (Attributes) A single object that describes Google Play Integrity API credential settings for Android device integrity detection. (see [below for nested schema](#nestedatt--oidc_options--mobile_app--integrity_detection--google_play))

<a id="nestedatt--oidc_options--mobile_app--integrity_detection--cache_duration"></a>
### Nested Schema for `oidc_options.mobile_app.integrity_detection.cache_duration`

Read-Only:

- `amount` (Number) An integer that specifies the number of minutes or hours that specify the duration between successful integrity detection calls.
- `units` (String) A string that specifies the cache duration time units.


<a id="nestedatt--oidc_options--mobile_app--integrity_detection--google_play"></a>
### Nested Schema for `oidc_options.mobile_app.integrity_detection.google_play`

Read-Only:

- `decryption_key` (String, Sensitive) Play Integrity verdict decryption key from your Google Play Services account. This parameter must be provided if you have set `verification_type` to `INTERNAL`.  Cannot be set with `service_account_credentials_json`.
- `service_account_credentials_json` (String, Sensitive) Contents of the JSON file that represents your Service Account Credentials.
- `verification_key` (String, Sensitive) Play Integrity verdict signature verification key from your Google Play Services account.
- `verification_type` (String) The type of verification.





<a id="nestedatt--saml_options"></a>
### Nested Schema for `saml_options`

Read-Only:

- `acs_urls` (Set of String) A list of string that specifies the Assertion Consumer Service URLs. The first URL in the list is used as default (there must be at least one URL).
- `assertion_duration` (Number) An integer that specifies the assertion validity duration in seconds.
- `assertion_signed_enabled` (Boolean) A boolean that specifies whether the SAML assertion itself should be signed.
- `cors_settings` (Attributes) A single object that allows customization of how the Authorization and Authentication APIs interact with CORS requests that reference the application. If omitted, the application allows CORS requests from any origin except for operations that expose sensitive information (e.g. `/as/authorize` and `/as/token`).  This is legacy behavior, and it is recommended that applications migrate to include specific CORS settings. (see [below for nested schema](#nestedatt--saml_options--cors_settings))
- `default_target_url` (String) A string that specfies a default URL used as the `RelayState` parameter by the IdP to deep link into the application after authentication. This value can be overridden by the `applicationUrl` query parameter for [GET Identity Provider Initiated SSO](https://apidocs.pingidentity.com/pingone/platform/v1/api/#get-identity-provider-initiated-sso). Although both of these parameters are generally URLs, because they are used as deep links, this is not enforced. If neither `defaultTargetUrl` nor `applicationUrl` is specified during a SAML authentication flow, no `RelayState` value is supplied to the application. The `defaultTargetUrl` (or the `applicationUrl`) value is passed to the SAML application’s ACS URL as a separate `RelayState` key value (not within the SAMLResponse key value).
- `enable_requested_authn_context` (Boolean) A boolean that specifies whether `requestedAuthnContext` is taken into account in policy decision-making.
- `home_page_url` (String) A string that specifies the custom home page URL for the application.
- `idp_signing_key` (Attributes) SAML application assertion/response signing key settings. (see [below for nested schema](#nestedatt--saml_options--idp_signing_key))
- `nameid_format` (String) A string that specifies the format of the Subject NameID attibute in the SAML assertion.
- `response_is_signed` (Boolean) A boolean that specifies whether the SAML assertion response itself should be signed.
- `session_not_on_or_after_duration` (Number) An integer that specifies a value for if the SAML application requires a different `SessionNotOnOrAfter` attribute value within the `AuthnStatement` element than the `NotOnOrAfter` value set by the `assertion_duration` property.
- `slo_binding` (String) A string that specifies the binding protocol to be used for the logout response.
- `slo_endpoint` (String) A string that specifies the logout endpoint URL.
- `slo_response_endpoint` (String) A string that specifies the endpoint URL to submit the logout response.
- `slo_window` (Number) An integer that defines how long (hours) PingOne can exchange logout messages with the application, specifically a logout request from the application, since the initial request.
- `sp_encryption` (Attributes) A single object that specifies settings for PingOne to encrypt SAML assertions to be sent to the application. Assertions are not encrypted by default. (see [below for nested schema](#nestedatt--saml_options--sp_encryption))
- `sp_entity_id` (String) A string that specifies the service provider entity ID used to lookup the application. This is a required property and is unique within the environment.
- `sp_verification` (Attributes) A single object that specifies SP signature verification settings. (see [below for nested schema](#nestedatt--saml_options--sp_verification))
- `type` (String) A string that specifies the type associated with the application.
- `virtual_server_id_settings` (Attributes) A single object that specifies settings for SAML virtual server IDs. (see [below for nested schema](#nestedatt--saml_options--virtual_server_id_settings))

<a id="nestedatt--saml_options--cors_settings"></a>
### Nested Schema for `saml_options.cors_settings`

Read-Only:

- `behavior` (String) A string that represents the behavior of how Authorization and Authentication APIs interact with CORS requests that reference the application.  Options are `ALLOW_NO_ORIGINS` (rejects all CORS requests), `ALLOW_SPECIFIC_ORIGINS` (rejects all CORS requests except those listed in `origins`).
- `origins` (Set of String) A set of strings that represent the origins from which CORS requests to the Authorization and Authentication APIs are allowed.  Each value will be a `http` or `https` URL without a path.  The host may be a domain name (including `localhost`), or an IPv4 address.  Subdomains may use the wildcard (`*`) to match any string.  Is expected to be non-empty when `behavior` is `ALLOW_SPECIFIC_ORIGINS` and is expected to be omitted or empty when `behavior` is `ALLOW_NO_ORIGINS`.  Limited to 20 values.


<a id="nestedatt--saml_options--idp_signing_key"></a>
### Nested Schema for `saml_options.idp_signing_key`

Read-Only:

- `algorithm` (String) A string that specifies the signature algorithm of the key.
- `key_id` (String) An ID for the certificate key pair to be used by the identity provider to sign assertions and responses.


<a id="nestedatt--saml_options--sp_encryption"></a>
### Nested Schema for `saml_options.sp_encryption`

Read-Only:

- `algorithm` (String) The algorithm to use when encrypting assertions.  Options are `AES_128`, `AES_256`, `TRIPLEDES`.
- `certificate` (Attributes) A single object that specifies the certificate settings used to encrypt SAML assertions. (see [below for nested schema](#nestedatt--saml_options--sp_encryption--certificate))

<a id="nestedatt--saml_options--sp_encryption--certificate"></a>
### Nested Schema for `saml_options.sp_encryption.certificate`

Read-Only:

- `id` (String) A string that specifies the unique identifier of the encryption public certificate that has been uploaded to PingOne.



<a id="nestedatt--saml_options--sp_verification"></a>
### Nested Schema for `saml_options.sp_verification`

Read-Only:

- `authn_request_signed` (Boolean) A boolean that specifies whether the Authn Request signing should be enforced.
- `certificate_ids` (Set of String) A list that specifies the certificate IDs used to verify the service provider signature.


<a id="nestedatt--saml_options--virtual_server_id_settings"></a>
### Nested Schema for `saml_options.virtual_server_id_settings`

Read-Only:

- `enabled` (Boolean) A boolean that specifies whether virtual server IDs are enabled for this SAML application.
- `virtual_server_ids` (Attributes List) A list of virtual server ID objects. Each object contains a virtual server ID and a flag indicating if it is the default. (see [below for nested schema](#nestedatt--saml_options--virtual_server_id_settings--virtual_server_ids))

<a id="nestedatt--saml_options--virtual_server_id_settings--virtual_server_ids"></a>
### Nested Schema for `saml_options.virtual_server_id_settings.virtual_server_ids`

Read-Only:

- `default` (Boolean) Whether this virtual server ID is the default.
- `vs_id` (String) The virtual server ID.




<a id="nestedatt--wsfed_options"></a>
### Nested Schema for `wsfed_options`

Read-Only:

- `audience_restriction` (String) The service provider ID. The default value is `urn:federation:MicrosoftOnline`.
- `cors_settings` (Attributes) A single object that allows customization of how the Authorization and Authentication APIs interact with CORS requests that reference the application. If omitted, the application allows CORS requests from any origin except for operations that expose sensitive information (e.g. `/as/authorize` and `/as/token`).  This is legacy behavior, and it is recommended that applications migrate to include specific CORS settings. (see [below for nested schema](#nestedatt--wsfed_options--cors_settings))
- `domain_name` (String) The federated domain name (for example, the Azure custom domain).
- `idp_signing_key` (Attributes) Contains the information about the signing of requests by the identity provider (IdP). (see [below for nested schema](#nestedatt--wsfed_options--idp_signing_key))
- `kerberos` (Attributes) The Kerberos authentication settings. Leave this out of the configuration to disable Kerberos authentication. (see [below for nested schema](#nestedatt--wsfed_options--kerberos))
- `reply_url` (String) The URL that the replying party (such as, Office365) uses to accept submissions of RequestSecurityTokenResponse messages that are a result of SSO requests.
- `slo_endpoint` (String) The single logout endpoint URL.
- `subject_name_identifier_format` (String) The format to use for the SubjectNameIdentifier element. Options are `urn:oasis:names:tc:SAML:1.1:nameid-format:unspecified`, `urn:oasis:names:tc:SAML:1.1:nameid-format:emailAddress`.
- `type` (String) A string that specifies the type associated with the application. This is a required property. Options are `WEB_APP`, `NATIVE_APP`, `SINGLE_PAGE_APP`, `WORKER`, `SERVICE`, `CUSTOM_APP`, `PORTAL_LINK_APP`.

<a id="nestedatt--wsfed_options--cors_settings"></a>
### Nested Schema for `wsfed_options.cors_settings`

Read-Only:

- `behavior` (String) A string that represents the behavior of how Authorization and Authentication APIs interact with CORS requests that reference the application.  Options are `ALLOW_NO_ORIGINS` (rejects all CORS requests), `ALLOW_SPECIFIC_ORIGINS` (rejects all CORS requests except those listed in `origins`).
- `origins` (Set of String) A set of strings that represent the origins from which CORS requests to the Authorization and Authentication APIs are allowed.  Each value will be a `http` or `https` URL without a path.  The host may be a domain name (including `localhost`), or an IPv4 address.  Subdomains may use the wildcard (`*`) to match any string.  Is expected to be non-empty when `behavior` is `ALLOW_SPECIFIC_ORIGINS` and is expected to be omitted or empty when `behavior` is `ALLOW_NO_ORIGINS`.  Limited to 20 values.


<a id="nestedatt--wsfed_options--idp_signing_key"></a>
### Nested Schema for `wsfed_options.idp_signing_key`

Read-Only:

- `algorithm` (String) A string that specifies the signature algorithm of the key.
- `key_id` (String) An ID for the certificate key pair to be used by the identity provider to sign assertions and responses.


<a id="nestedatt--wsfed_options--kerberos"></a>
### Nested Schema for `wsfed_options.kerberos`

Read-Only:

- `gateways` (Attributes Set) The LDAP gateway properties. (see [below for nested schema](#nestedatt--wsfed_options--kerberos--gateways))

<a id="nestedatt--wsfed_options--kerberos--gateways"></a>
### Nested Schema for `wsfed_options.kerberos.gateways`

Read-Only:

- `id` (String) The UUID of the LDAP gateway. Must be a valid PingOne resource ID.
- `type` (String) The gateway type. This must be "LDAP".
- `user_type` (Attributes) The object reference to the user type in the list of "userTypes" for the LDAP gateway. (see [below for nested schema](#nestedatt--wsfed_options--kerberos--gateways--user_type))

<a id="nestedatt--wsfed_options--kerberos--gateways--user_type"></a>
### Nested Schema for `wsfed_options.kerberos.gateways.user_type`

Read-Only:

- `id` (String) The UUID of a user type in the list of `userTypes` for the LDAP gateway. Must be a valid PingOne resource ID.
