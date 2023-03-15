---
page_title: "pingone_environment Resource - terraform-provider-pingone"
subcategory: "Platform"
description: |-
  Resource to create and manage PingOne environments.
---

# pingone_environment (Resource)

Resource to create and manage PingOne environments.

~> PingOne environments are created with a default population and at least one service added.

~> This `pingone_environment` resource does not yet support creation of WORKFORCE enabled environments.

## Example Usage

```terraform
resource "pingone_environment" "my_environment" {
  name        = "New Environment"
  description = "My new environment"
  type        = "SANDBOX"
  license_id  = var.license_id

  default_population {
    name        = "My Population"
    description = "My new population for users"
  }

  service {
    type = "SSO"
  }

  service {
    type = "MFA"
  }

  service {
    type        = "PingFederate"
    console_url = "https://my-pingfederate-console.example.com/pingfederate"
  }
}
```

## Schema

### Required

- `default_population` (Block List, Min: 1, Max: 1) The environment's default population. (see [below for nested schema](#nestedblock--default_population))
- `license_id` (String) An ID of a valid license to apply to the environment.
- `name` (String) The name of the environment.
- `service` (Block Set, Min: 1, Max: 13) The services to enable in the environment. (see [below for nested schema](#nestedblock--service))

### Optional

- `description` (String) A description of the environment.
- `region` (String) The region to create the environment in.  Should be consistent with the PingOne organisation region.  Valid options are `AsiaPacific` `Canada` `Europe` and `NorthAmerica`.  Default can be set with the `PINGONE_REGION` environment variable.
- `solution` (String) The solution context of the environment.  Leave blank for a custom, non-workforce solution context.  Valid options are `CUSTOMER`, or no value for custom solution context.  Workforce solution environments are not yet supported in this provider resource, but can be fetched using the `pingone_environment` datasource.
- `timeouts` (Block, Optional) (see [below for nested schema](#nestedblock--timeouts))
- `type` (String) The type of the environment to create.  Options are `SANDBOX` for a development/testing environment and `PRODUCTION` for environments that require protection from deletion. Defaults to `SANDBOX`.

### Read-Only

- `default_population_id` (String) The ID of the environment's default population.  This attribute is only populated when also using the `default_population` block to define a default population.
- `id` (String) The ID of this resource.
- `organization_id` (String) The ID of the PingOne organization tenant to which the environment belongs.

<a id="nestedblock--default_population"></a>
### Nested Schema for `default_population`

Optional:

- `description` (String) A description to apply to the environment's default population.
- `name` (String) The name of the environment's default population. Defaults to `Default`.


<a id="nestedblock--service"></a>
### Nested Schema for `service`

Optional:

- `bookmark` (Block Set) Custom bookmark links for the service. (see [below for nested schema](#nestedblock--service--bookmark))
- `console_url` (String) A custom console URL to set.  Generally used with services that are deployed separately to the PingOne SaaS service, such as `PingFederate`, `PingAccess`, `PingDirectory`, `PingAuthorize` and `PingCentral`.
- `type` (String) The service type to enable in the environment.  Valid options are `APIIntelligence`, `Authorize`, `Credentials`, `DaVinci`, `MFA`, `PingAccess`, `PingAuthorize`, `PingCentral`, `PingDirectory`, `PingFederate`, `PingID`, `Risk`, `SSO`, `Verify`.  Defaults to `SSO`.

<a id="nestedblock--service--bookmark"></a>
### Nested Schema for `service.bookmark`

Required:

- `name` (String) Bookmark name.
- `url` (String) Bookmark URL.



<a id="nestedblock--timeouts"></a>
### Nested Schema for `timeouts`

Optional:

- `create` (String)

## Import

Import is supported using the following syntax examples, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
$ terraform import pingone_environment.example <environment_id>
```

~> The following import command is supported for backward capability with previous provider versions or where a default population doesn't exist in the environment.  The following example associates the defined population as the default population in the resource.  Once imported, the `environment_id` of the import command will be stored in the `id` attribute of the environment resource, and the `population_id` stored in the `default_population_id` attribute of the environment resource.

```shell
$ terraform import pingone_environment.example <environment_id>/<population_id>
```
