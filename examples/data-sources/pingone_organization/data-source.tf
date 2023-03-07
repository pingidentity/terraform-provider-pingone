data "pingone_organization" "example_by_id" {
  organization_id = var.organization_id
}

data "pingone_organization" "example_by_name" {
  name = var.organization_name
}