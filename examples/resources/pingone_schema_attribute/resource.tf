resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_schema_attribute" "my_attribute" {
  environment_id = pingone_environment.my_environment.id

  name         = "myAttribute"
  display_name = "My Awesome Attribute"
  description  = "My new awesome attribute"

  type        = "STRING"
  unique      = false
  multivalued = false

  lifecycle {
    # change the `prevent_destroy` parameter value to `true` to prevent this data carrying resource from being destroyed
    prevent_destroy = false
  }
}
