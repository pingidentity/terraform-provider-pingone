resource "pingone_davinci_connector_instance" "connectorOpswat" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorOpswat"
  }
  name = "My awesome connectorOpswat"
  properties = jsonencode({
    "clientID" = var.connectoropswat_property_client_i_d
    "clientSecret" = var.connectoropswat_property_client_secret
    "crossDomainApiPort" = var.connectoropswat_property_cross_domain_api_port
    "maDomain" = var.connectoropswat_property_ma_domain
  })
}
