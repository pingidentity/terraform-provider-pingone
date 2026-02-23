resource "pingone_environment" "my_environment" {
  # ...
}

data "pingone_schema" "user" {
  # ...
}

data "pingone_schema_attribute" "email" {
  environment_id = pingone_environment.my_environment.id
  schema_id      = data.pingone_schema.user.id
  name           = "email"
}

import {
  to = pingone_schema_attribute.email
  id = "${pingone_environment.my_environment.id}/${data.pingone_schema.user.id}/${data.pingone_schema_attribute.email.id}"
}

resource "pingone_schema_attribute" "email" {
  environment_id = pingone_environment.my_environment.id
  name           = data.pingone_schema.attribute.email.name

  enabled = true
}
