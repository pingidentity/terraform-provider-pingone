---
page_title: "pingone_application_sign_on_policy_assignments Data Source - terraform-provider-pingone"
subcategory: "SSO"
description: |-
  Datasource to retrieve the IDs, as a collection, of PingOne Sign On Policy assignments for an application in an environment.
---

# pingone_application_sign_on_policy_assignments (Data Source)

Datasource to retrieve the IDs, as a collection, of PingOne Sign On Policy assignments for an application in an environment.

## Example Usage

```terraform
data "pingone_application_sign_on_policy_assignments" "all_sop_assignments_by_app" {
  environment_id = var.environment_id

  application_id = pingone_application.my_awesome_application.id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `application_id` (String) The ID of the application to filter application sign on policy assignments from.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.
- `environment_id` (String) The ID of the environment to filter application sign on policy assignments from.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.

### Read-Only

- `ids` (List of String) The list of resulting IDs of application sign on policy assignments that have been successfully retrieved for an application.
