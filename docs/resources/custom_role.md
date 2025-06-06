---
page_title: "pingone_custom_role Resource - terraform-provider-pingone"
subcategory: "Platform"
description: |-
  Resource to create and manage a custom administrator role in an environment.
---

# pingone_custom_role (Resource)

Resource to create and manage a custom administrator role in an environment.

## Example Usage

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

data "pingone_role" "environment_admin" {
  name = "Environment Admin"
}

resource "pingone_custom_role" "my_custom_role" {
  environment_id = pingone_environment.my_environment.id

  name        = "My custom role"
  description = "My custom role for reading role assignments"

  applicable_to = [
    "ENVIRONMENT",
    "POPULATION"
  ]

  can_be_assigned_by = [
    {
      id = pingone_role.environment_admin.id
    }
  ]

  permissions = [
    {
      id = "permissions:read:userRoleAssignments"
    },
    {
      id = "permissions:read:groupRoleAssignments"
    },
  ]

}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `applicable_to` (Set of String) The scope types to which the role can be applied. Options are `ORGANIZATION`, `ENVIRONMENT`, `POPULATION`, `APPLICATION`. At least one value must be set.
- `can_be_assigned_by` (Attributes Set) A relationship that determines whether a user assigned to one of this set of roles for a jurisdiction can assign the current custom role to another user for the same jurisdiction or sub-jurisdiction. (see [below for nested schema](#nestedatt--can_be_assigned_by))
- `environment_id` (String) The ID of the environment to create and manage the custom role in.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `name` (String) The role name.
- `permissions` (Attributes Set) The set of permissions assigned to the role. For possible values, see the [list of available permissions](https://apidocs.pingidentity.com/pingone/platform/v1/api/#pingone-permissions-by-identifier). At least one permission must be set. (see [below for nested schema](#nestedatt--permissions))

### Optional

- `description` (String) The description of the role.

### Read-Only

- `can_assign` (Attributes Set) A relationship that specifies if an actor is assigned the current custom role for a jurisdiction, then the actor can assign any of this set of roles to another actor for the same jurisdiction or sub-jurisdiction. This capability is derived from the `can_be_assigned_by` property. (see [below for nested schema](#nestedatt--can_assign))
- `id` (String) The ID of this resource.
- `type` (String) A value that indicates whether the role is a built-in role or a custom role. Options are `PLATFORM` and `CUSTOM`. This will always be `CUSTOM` for custom roles.

<a id="nestedatt--can_be_assigned_by"></a>
### Nested Schema for `can_be_assigned_by`

Required:

- `id` (String) The ID of the role that can assign the current custom role.


<a id="nestedatt--permissions"></a>
### Nested Schema for `permissions`

Required:

- `id` (String) The ID of the permission assigned to this role.


<a id="nestedatt--can_assign"></a>
### Nested Schema for `can_assign`

Read-Only:

- `id` (String) The ID of a role that can be assigned by an actor assigned the current custom role.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_custom_role.example <environment_id>/<role_id>
```
