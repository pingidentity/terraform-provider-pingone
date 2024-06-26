---
page_title: "pingone_application_resource_permission Resource - terraform-provider-pingone"
subcategory: "Authorize"
description: |-
  Resource to create and manage application resource permissions in a PingOne environment.
---

# pingone_application_resource_permission (Resource)

Resource to create and manage application resource permissions in a PingOne environment.

## Example Usage

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_resource" "my_awesome_custom_resource" {
  environment_id = pingone_environment.my_environment.id

  name = "My awesome custom resource"
}

resource "pingone_application_resource" "my_custom_application_resource" {
  environment_id = pingone_environment.my_environment.id
  resource_name  = pingone_resource.my_awesome_custom_resource.name

  name        = "Invoices"
  description = "My invoices resource application"
}

resource "pingone_application_resource_permission" "my_custom_application_resource_permission" {
  environment_id          = pingone_environment.my_environment.id
  application_resource_id = pingone_application_resource.my_custom_application_resource.id

  action      = "Invoices-Read"
  description = "Read Invoices"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `action` (String) A string that specifies the action associated with this permission.  The action must contain only Unicode letters, marks, numbers, spaces, forward slashes, dots, apostrophes, underscores, or hyphens.
- `application_resource_id` (String) The ID of the application resource to create and manage permissions for.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `environment_id` (String) The ID of the environment to create the resource attribute in.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.

### Optional

- `description` (String) A string that specifies the permission's description.

### Read-Only

- `id` (String) The ID of this resource.
- `resource` (Attributes) A single object that describes the associated application resource. (see [below for nested schema](#nestedatt--resource))

<a id="nestedatt--resource"></a>
### Nested Schema for `resource`

Read-Only:

- `id` (String) A string that specifies the ID for the associated application resource.
- `name` (String) A string that represents the name of the associated application resource.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_application_resource_permission.example <environment_id>/<application_resource_id>/<application_resource_permission_id>
```
