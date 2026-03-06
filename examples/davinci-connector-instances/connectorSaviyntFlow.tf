resource "pingone_davinci_connector_instance" "connectorSaviyntFlow" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSaviyntFlow"
  }
  name = "My awesome connectorSaviyntFlow"
  properties = jsonencode({
    "domainName" = var.connectorsaviyntflow_property_domain_name
    "path" = var.connectorsaviyntflow_property_path
    "saviyntPassword" = var.connectorsaviyntflow_property_saviynt_password
    "saviyntUserName" = var.connectorsaviyntflow_property_saviynt_user_name
  })
}
