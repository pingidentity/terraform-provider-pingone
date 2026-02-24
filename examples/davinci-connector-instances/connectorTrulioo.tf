resource "pingone_davinci_connector_instance" "connectorTrulioo" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorTrulioo"
  }
  name = "My awesome connectorTrulioo"
  properties = jsonencode({
    "clientID" = var.connectortrulioo_property_client_i_d
    "clientSecret" = var.connectortrulioo_property_client_secret
  })
}
