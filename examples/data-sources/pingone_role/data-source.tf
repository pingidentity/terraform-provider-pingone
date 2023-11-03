data "pingone_role" "organisation_admin" {
  name = "Organization Admin"
}

data "pingone_role" "environment_admin" {
  name = "Environment Admin"
}

data "pingone_role" "identity_data_admin" {
  name = "Identity Data Admin"
}

data "pingone_role" "davinci_admin" {
  name = "DaVinci Admin"
}

data "pingone_role" "client_application_developer" {
  name = "Client Application Developer"
}

data "pingone_role" "identity_data_admin_ro" {
  name = "Identity Data Read Only"
}

data "pingone_role" "davinci_admin_ro" {
  name = "DaVinci Admin Read Only"
}

data "pingone_role" "role_by_id" {
  role_id = var.role_id
}
