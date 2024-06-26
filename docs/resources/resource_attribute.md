---
page_title: "pingone_resource_attribute Resource - terraform-provider-pingone"
subcategory: "SSO"
description: |-
  Resource to create and manage resource attributes in PingOne.
---

# pingone_resource_attribute (Resource)

Resource to create and manage resource attributes in PingOne.

## Example Usage

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_resource" "my_resource" {
  environment_id = pingone_environment.my_environment.id

  name = "My resource"
}

resource "pingone_resource_attribute" "my_custom_resource_attribute" {
  environment_id = pingone_environment.my_environment.id

  resource_type      = "CUSTOM"
  custom_resource_id = pingone_resource.my_resource.id

  name  = "example_attribute"
  value = "$${user.name.family}"
}

resource "pingone_resource_attribute" "my_openid_connect_resource_attribute" {
  environment_id = pingone_environment.my_environment.id

  resource_type = "OPENID_CONNECT"

  name  = "example_attribute"
  value = "$${user.name.family}"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment to create the resource attribute in.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `name` (String) A string that specifies the name of the resource attribute to map a value for. When the resource's type property is `OPENID_CONNECT`, the following are reserved names and cannot be used: `acr`, `amr`, `aud`, `auth_time`, `client_id`, `env`, `exp`, `iat`, `iss`, `jti`, `org`, `p1.*`, `scope`, `sid`, `sub`.  The resource will also override the default configured values for a resource, rather than creating new attributes.  For resources of type `CUSTOM`, the `sub` name is overridden.  For resources of type `OPENID_CONNECT`, the following names are overridden: `address.country`, `address.formatted`, `address.locality`, `address.postal_code`, `address.region`, `address.street_address`, `birthdate`, `email`, `email_verified`, `family_name`, `gender`, `given_name`, `locale`, `middle_name`, `name`, `nickname`, `phone_number`, `phone_number_verified`, `picture`, `preferred_username`, `profile`, `updated_at`, `website`, `zoneinfo`.
- `resource_type` (String) The type of the resource to create the attribute for.  When the value is set to `CUSTOM`, `custom_resource_id` must be specified.  Options are `CUSTOM`, `OPENID_CONNECT`.
- `value` (String) A string that specifies the value of the custom resource attribute. This value can be a placeholder that references an attribute in the user schema, expressed as `${user.path.to.value}`, or it can be an expression, or a static string. Placeholders must be valid, enabled attributes in the environment’s user schema. Examples of valid values are: `${user.email}`, `${user.name.family}`, and `myClaimValueString`.  Note that definition in HCL requires escaping with the `$` character when defining attribute paths, for example `value = "$${user.email}"`.

### Optional

- `custom_resource_id` (String) A string that specifies the ID of the custom resource to create the attribute for.  Must be a valid PingOne resource ID.  Required if `resource_type` is set to `CUSTOM`, but cannot be set if `resource_type` is set to `OPENID_CONNECT`.
- `id_token_enabled` (Boolean) A boolean that specifies whether the attribute mapping should be available in the ID Token.  Only applies to resources that are of type `OPENID_CONNECT` and the `id_token_enabled` and `userinfo_enabled` properties cannot both be set to false. Defaults to `true`.
- `userinfo_enabled` (Boolean) A boolean that specifies whether the attribute mapping should be available through the /as/userinfo endpoint.  Only applies to resources that are of type `OPENID_CONNECT` and the `id_token_enabled` and `userinfo_enabled` properties cannot both be set to false. Defaults to `true`.

### Read-Only

- `id` (String) The ID of this resource.
- `resource_id` (String) The ID of the resource that the attribute is assigned to.
- `type` (String) A string that specifies the type of resource attribute. Options are: `CORE` (The claim is required and cannot not be removed), `CUSTOM` (The claim is not a CORE attribute. All created attributes are of this type), `PREDEFINED` (A designation for predefined OIDC resource attributes such as given_name. These attributes cannot be removed; however, they can be modified).

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_resource_attribute.example <environment_id>/<resource_id>/<resource_attribute_id>
```
