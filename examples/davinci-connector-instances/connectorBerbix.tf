resource "pingone_davinci_connector_instance" "connectorBerbix" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBerbix"
  }
  name = "My awesome connectorBerbix"
  properties = jsonencode({
    "domainName" = var.connectorberbix_property_domain_name
    "path" = var.connectorberbix_property_path
    "username" = var.connectorberbix_property_username
  })
}
