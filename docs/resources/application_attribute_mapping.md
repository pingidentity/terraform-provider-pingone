---
page_title: "pingone_application_attribute_mapping Resource - terraform-provider-pingone"
subcategory: "SSO"
description: |-
  Resource to create and manage custom attribute mappings for administrator defined applications configured in PingOne.
---

# pingone_application_attribute_mapping (Resource)

Resource to create and manage custom attribute mappings for administrator defined applications configured in PingOne.

~> Attributes can only be mapped to administrator defined applications that are managed through the `pingone_application` resource.

## Example Usage - OIDC Application

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_application" "my_awesome_spa" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Single Page App"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://my-website.com"]
  }
}

resource "pingone_application_attribute_mapping" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  name  = "email"
  value = "$${user.email}"
}

resource "pingone_application_attribute_mapping" "bar" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  name  = "full_name"
  value = "$${user.name.given + ', ' + user.name.family}"
}
```

## Example Usage - OIDC Application (OIDC Resource Scope)

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

data "pingone_resource" "openid_resource" {
  environment_id = var.environment_id

  name = "openid"
}

data "pingone_resource_scope" "openid_profile" {
  environment_id = var.environment_id
  resource_id    = data.pingone_resource.openid_resource.id

  name = "profile"
}

resource "pingone_application" "my_awesome_spa" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Single Page App"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://my-website.com"]
  }
}

resource "pingone_application_resource_grant" "oidc_grant" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  resource_id = data.pingone_resource.openid_resource.id

  scopes = [
    data.pingone_resource_scope.openid_profile.id
  ]
}

resource "pingone_application_attribute_mapping" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  name  = "customAttribute"
  value = "$${user.email}"

  oidc_scopes = [
    data.pingone_resource_scope.openid_profile.id
  ]

  oidc_id_token_enabled = true
  oidc_userinfo_enabled = false

  depends_on = [
    pingone_application_resource_grant.oidc_grant
  ]
}
```

## Example Usage - OIDC Application (Custom Resource Scope)

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_resource" "my_resource" {
  environment_id = pingone_environment.my_environment.id

  name = "My resource"
}

resource "pingone_resource_scope" "my_resource_scope" {
  environment_id = pingone_environment.my_environment.id
  resource_id    = pingone_resource.my_resource.id

  name = "example_scope"
}

resource "pingone_application" "my_awesome_spa" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Single Page App"
  enabled        = true

  oidc_options {
    type                        = "SINGLE_PAGE_APP"
    grant_types                 = ["AUTHORIZATION_CODE"]
    response_types              = ["CODE"]
    pkce_enforcement            = "S256_REQUIRED"
    token_endpoint_authn_method = "NONE"
    redirect_uris               = ["https://my-website.com"]
  }
}

resource "pingone_application_resource_grant" "custom_grant" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  resource_id = pingone_resource.my_resource.id

  scopes = [
    pingone_resource_scope.my_resource_scope.id
  ]
}

resource "pingone_application_attribute_mapping" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_spa.id

  name  = "customAttribute"
  value = "$${user.email}"

  oidc_scopes = [
    pingone_resource_scope.my_resource_scope.id
  ]

  oidc_id_token_enabled = true
  oidc_userinfo_enabled = false

  depends_on = [
    pingone_application_resource_grant.custom_grant
  ]
}
```

## Example Usage - SAML Application

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

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

  saml_options {
    acs_urls           = ["https://my-saas-app.com"]
    assertion_duration = 3600
    sp_entity_id       = "sp:entity:localhost"

    idp_signing_key {
      key_id    = pingone_key.my_awesome_key.id
      algorithm = pingone_key.my_awesome_key.signature_algorithm
    }

    sp_verification_certificate_ids = [var.sp_verification_certificate_id]
  }
}

resource "pingone_application_attribute_mapping" "foo" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_saml_app.id

  name  = "email"
  value = "$${user.email}"
}

resource "pingone_application_attribute_mapping" "bar" {
  environment_id = pingone_environment.my_environment.id
  application_id = pingone_application.my_awesome_saml_app.id

  name  = "full_name"
  value = "$${user.name.given + ', ' + user.name.family}"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `application_id` (String) The ID of the application to create the attribute mapping for.
- `environment_id` (String) The ID of the environment to create the application attribute mapping in.
- `name` (String) A string that specifies the name of attribute and must be unique within an application. For SAML applications, the `saml_subject` name is a case-insensitive name which indicates the mapping to be used for the subject in an assertion and can be overridden. For OpenID Connect applications, the `sub` name indicates the mapping to be used for the subject in the token and can be overridden.  The following OpenID Connect names are reserved and cannot be used: `acr`, `amr`, `at_hash`, `aud`, `auth_time`, `azp`, `client_id`, `exp`, `iat`, `iss`, `jti`, `nbf`, `nonce`, `org`, `scope`, `sid`.
- `value` (String) A string that specifies the string constants or expression for mapping the attribute path against a specific source. The expression format is `${<source>.<attribute_path>}`. The only supported source is user (for example, `${user.id}`).  When defining attribute mapping values in Terraform, the expression must be escaped (for example `value = "$${user.id}}"`)

### Optional

- `oidc_id_token_enabled` (Boolean) Whether the attribute mapping should be available in the ID Token. This property is applicable only when the application's `protocol` property is `OPENID_CONNECT`. If omitted, the default is `true`. Note that the `id_token_enabled` and `userinfo_enabled` properties cannot both be set to `false`. At least one of these properties must have a value of `true`.
- `oidc_scopes` (Set of String) OIDC resource scope IDs that this attribute mapping is available for exclusively. This setting overrides any global OIDC resource scopes that contain an attribute mapping with the same name. The list can contain only scope IDs that have been granted for the application through the `/grants` endpoint. At least one scope ID is expected.
- `oidc_userinfo_enabled` (Boolean) Whether the attribute mapping should be available through the `/as/userinfo` endpoint. This property is applicable only when the application's protocol property is `OPENID_CONNECT`. If omitted, the default is `true`. Note that the `id_token_enabled` and `userinfo_enabled` properties cannot both be set to `false`. At least one of these properties must have a value of `true`.
- `required` (Boolean) A boolean to specify whether a mapping value is required for this attribute. If `true`, a value must be set and a non-empty value must be available in the SAML assertion or ID token. If overriding a core attribute mapping (`saml_subject` for SAML applications and `sub` for OpenID Connect applications), then this value must be set to `true`.  Defaults to `false`.
- `saml_subject_nameformat` (String) A URI reference representing the classification of the attribute, which helps the service provider interpret the attribute format.  This property is applicable only when the application's protocol property is `SAML` and the name is the `saml_subject` core attribute.  Examples include `urn:oasis:names:tc:SAML:2.0:attrname-format:unspecified`, `urn:oasis:names:tc:SAML:2.0:attrname-format:uri`, `urn:oasis:names:tc:SAML:2.0:attrname-format:basic`.

### Read-Only

- `id` (String) The ID of this resource.
- `mapping_type` (String) A string that specifies the mapping type of the attribute. Options are `CORE`, `SCOPE`, and `CUSTOM`.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
$ terraform import pingone_application_attribute_mapping.example <environment_id>/<application_id>/<attribute_mapping_id>
```
