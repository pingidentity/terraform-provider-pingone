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

resource "pingone_application_resource_permission" "my_custom_application_resource_permission" {
  environment_id          = pingone_environment.my_environment.id
  application_resource_id = pingone_application_resource.my_custom_application_resource.id

  action      = "Invoices:Read"
  description = "Read Invoices"
}

resource "pingone_authorize_application_role" "my_awesome_application_role" {
  environment_id = pingone_environment.my_environment.id

  name        = "CEO"
  description = "The CEO"
}

resource "pingone_authorize_application_role_permission" "my_awesome_application_role_permission" {
  environment_id = pingone_environment.my_environment.id

  application_role_permission_id = pingone_authorize_application_role.my_awesome_application_role.id
  application_resource_permission_id = pingone_application_resource_permission.my_custom_application_resource_permission.id
}
