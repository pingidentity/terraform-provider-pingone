---
page_title: "pingone_identity_propagation_plan Resource - terraform-provider-pingone"
subcategory: "Platform"
description: |-
  Resource to create and manage PingOne identity propagation (provisioning) plans for an environment.
---

# pingone_identity_propagation_plan (Resource)

Resource to create and manage PingOne identity propagation (provisioning) plans for an environment.

~> Only one `pingone_identity_propagation_plan` resource can be configured for an environment.  If multiple `pingone_identity_propagation_plan` resource definitions exist in HCL code for a single environment, the platform will return an error.

## Example Usage

```terraform
resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_identity_propagation_plan" "my_awesome_propagation_plan" {
  environment_id = pingone_environment.my_environment.id
  name           = "My Awesome Identity Provisioning Plan"
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment to manage the identity propagation plan in.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `name` (String) A string that specifies the unique name of the identity propagation plan.

### Read-Only

- `id` (String) The ID of this resource.

## Import

Import is supported using the following syntax, where attributes in `<>` brackets are replaced with the relevant ID.  For example, `<environment_id>` should be replaced with the ID of the environment to import from.

```shell
terraform import pingone_identity_propagation_plan.example <environment_id>/<identity_propagation_plan_id>
```
