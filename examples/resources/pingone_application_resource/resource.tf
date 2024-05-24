resource "pingone_environment" "my_environment" {
  # ...
}

resource "pingone_resource" "my_awesome_custom_resource" {
  environment_id = pingone_environment.my_environment.id

  name = "My awesome custom resource"
}

resource "pingone_application_resource" "my_custom_application_resource" {
  environment_id = pingone_environment.my_environment.id
  resource_name  = pingone_resource.my_resource.name

  name        = "Invoices"
  description = "My invoices resource application"
}
