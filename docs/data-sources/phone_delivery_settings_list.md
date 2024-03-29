---
page_title: "pingone_phone_delivery_settings_list Data Source - terraform-provider-pingone"
subcategory: "Platform"
description: |-
  Datasource to retrieve multiple phone delivery settings in a PingOne environments.
---

# pingone_phone_delivery_settings_list (Data Source)

Datasource to retrieve multiple phone delivery settings in a PingOne environments.

## Example Usage

```terraform
data "pingone_phone_delivery_settings_list" "all" {
  environment_id = var.environment_id
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `environment_id` (String) The ID of the environment to filter phone delivery settings senders from.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (List of String) The list of resulting IDs of phone delivery settings senders that have been successfully retrieved.
