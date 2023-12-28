resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_group" "parent_group" {
  environment_id = pingone_environment.my_environment.id

  name = "My parent group"

  lifecycle {
    # change the `prevent_destroy` parameter value to `true` to prevent this data carrying resource from being destroyed
    prevent_destroy = false
  }
}

resource "pingone_group" "nested_group" {
  environment_id = pingone_environment.my_environment.id

  name = "My nested group"

  lifecycle {
    # change the `prevent_destroy` parameter value to `true` to prevent this data carrying resource from being destroyed
    prevent_destroy = false
  }
}

resource "pingone_group_nesting" "my_group_nesting" {
  environment_id  = pingone_environment.my_environment.id
  group_id        = pingone_group.parent_group.id
  nested_group_id = pingone_group.nested_group.id
}