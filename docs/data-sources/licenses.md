---
page_title: "pingone_licenses Data Source - terraform-provider-pingone"
subcategory: "Platform"
description: |-
  Datasource to retrieve multiple PingOne license IDs selected by a SCIM filter or a name/value list combination.
---

# pingone_licenses (Data Source)

Datasource to retrieve multiple PingOne license IDs selected by a SCIM filter or a name/value list combination.

## Example Usage

```terraform
data "pingone_licenses" "my_licenses_by_scim_filter" {
  organization_id = var.organization_id
  scim_filter     = "(status eq \"active\") and (beginsAt lt \"2009-11-10T23:00:00Z\")"
}

data "pingone_licenses" "my_licenses_by_data_filter" {
  organization_id = var.organization_id

  data_filters = [
    {
      name   = "name"
      values = ["My License"]
    },
    {
      name   = "status"
      values = ["ACTIVE"]
    }
  ]
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `organization_id` (String) The ID of the organization to retrieve licenses for.  Must be a valid PingOne resource ID.  This field is immutable and will trigger a replace plan if changed.

### Optional

- `data_filters` (Attributes List) Individual data filters to apply to the license selection.  If the attribute filter is `status`, available values are `ACTIVE`, `EXPIRED`, `FUTURE` and `TERMINATED`.  Allowed attributes to filter: `name`, `package`, `status`.  Exactly one of the following must be defined: `scim_filter`, `data_filters`. (see [below for nested schema](#nestedatt--data_filters))
- `scim_filter` (String) A SCIM filter to apply to the license selection.  A SCIM filter offers the greatest flexibility in filtering licenses.  If the attribute filter is `status`, available values are `ACTIVE`, `EXPIRED`, `FUTURE` and `TERMINATED`.  The SCIM filter can use the following attributes: `name`, `package`, `status`.  Exactly one of the following must be defined: `scim_filter`, `data_filters`.

### Read-Only

- `id` (String) The ID of this resource.
- `ids` (List of String) The list of resulting IDs of licenses that have been successfully retrieved and filtered.

<a id="nestedatt--data_filters"></a>
### Nested Schema for `data_filters`

Required:

- `name` (String) The attribute name to filter on.  Must be one of the following values: `name`, `package`, `status`.
- `values` (List of String) The possible values (case sensitive) of the attribute defined in the `name` parameter to filter.
