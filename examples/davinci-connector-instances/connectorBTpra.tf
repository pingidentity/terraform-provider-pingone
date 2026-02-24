resource "pingone_davinci_connector_instance" "connectorBTpra" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBTpra"
  }
  name = "My awesome connectorBTpra"
  properties = jsonencode({
    "clientID" = var.connectorbtpra_property_client_i_d
    "clientSecret" = var.connectorbtpra_property_client_secret
    "praAPIurl" = var.pra_api_url
  })
}
