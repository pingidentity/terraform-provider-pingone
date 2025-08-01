---
page_title: "pingone_population_default Resource - terraform-provider-pingone"
subcategory: "SSO"
description: |-
  Resource to overwrite the default PingOne population, or create it if it doesn't already exist.
---

# pingone_population_default (Resource)

Resource to overwrite the default PingOne population, or create it if it doesn't already exist.

~> When destroying the resource, the default population will be reset to it's original configuration, and then removed from Terraform's state.  The population itself (and any user data contained in the population) will not be removed from the PingOne service.

## Example Usage

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_population_default" "my_default_population" {
  environment_id = pingone_environment.my_environment.id

  name        = "My default population"
  description = "A resource that overwrites the default population"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment to manage the default population in.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `name` (String) The name to apply to the default population.

### Optional

- `alternative_identifiers` (Set of String) Alternative identifiers that can be used to search for populations besides `name`.
- `description` (String) A description to apply to the default population.
- `password_policy` (Attributes) The object reference to the password policy resource. This is an optional property. Conflicts with `password_policy_id`. (see [below for nested schema](#nestedatt--password_policy))
- `password_policy_id` (String, Deprecated) A string that specifies the ID of a password policy to assign to the population.  Must be a valid PingOne resource ID. The "password_policy.id" attribute should be used instead of this attribute.  Conflicts with "password_policy".
- `preferred_language` (String) The language locale for the population. If absent, the environment default is used.
- `theme` (Attributes) The object reference to the theme resource. (see [below for nested schema](#nestedatt--theme))

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedatt--password_policy"></a>
### Nested Schema for `password_policy`

Required:

- `id` (String) The ID of the password policy that is used for this population. If absent, the environment's default is used. Must be a valid PingOne resource ID.


<a id="nestedatt--theme"></a>
### Nested Schema for `theme`

Required:

- `id` (String) The ID of the theme to use for the population. If absent, the environment's default is used. Must be a valid PingOne resource ID.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_population_default.example <environment_id>
```
