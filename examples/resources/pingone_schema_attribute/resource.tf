resource "pingone_environment" "my_environment" {
  # ...
}

data "pingone_schema" "users" {
  environment_id = pingone_environment.my_environment.id

  name = "User"

}

resource "pingone_schema_attribute" "my_attribute" {
  environment_id = pingone_environment.my_environment.id
  schema_id      = data.pingone_schema.users.id

  name         = "myAttribute"
  display_name = "My Awesome Attribute"
  description  = "My new awesome attribute"

  type        = "STRING"
  unique      = false
  multivalued = false
}
