---
page_title: "pingone_resource Data Source - terraform-provider-pingone"
subcategory: "SSO"
description: |-
  Datasource to read PingOne OAuth 2.0 resource data.
---

# pingone_resource (Data Source)

Datasource to read PingOne OAuth 2.0 resource data.

## Example Usage

```terraform
data "pingone_resource" "example_by_name" {
  environment_id = var.environment_id

  name = "openid"
}

data "pingone_resource" "example_by_id" {
  environment_id = var.environment_id

  resource_id = var.resource_id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment that is configured with the resource.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.

### Optional

- `name` (String) The name of the resource.  At least one of the following must be defined: `resource_id`, `name`.
- `resource_id` (String) The ID of the resource.  At least one of the following must be defined: `resource_id`, `name`.  Must be a valid PingOne resource ID.

### Read-Only

- `access_token_validity_seconds` (Number) An integer that specifies the number of seconds that the access token is valid.
- `audience` (String) A string that specifies a URL without a fragment or `@ObjectName` and must not contain `pingone` or `pingidentity` (for example, `https://api.myresource.com`). If a URL is not specified, the resource name is used.
- `client_secret` (String, Sensitive) An auto-generated resource client secret.
- `description` (String) A description of the resource.
- `id` (String) The ID of this resource.
- `introspect_endpoint_auth_method` (String) The client authentication methods supported by the token endpoint.  Options are `CLIENT_SECRET_BASIC`, `CLIENT_SECRET_POST`, `NONE`.
- `type` (String) A string that specifies the type of resource.  Options are `CUSTOM` (specifies the a resource that has been created by admin), `OPENID_CONNECT` (specifies the built-in platform resource for OpenID Connect), `PINGONE_API` (specifies the built-in platform resource for PingOne).
