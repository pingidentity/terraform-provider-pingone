resource "pingone_davinci_connector_instance" "connectorBTrs" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBTrs"
  }
  name = "My awesome connectorBTrs"
  properties = jsonencode({
    "clientID" = var.connectorbtrs_property_client_i_d
    "clientSecret" = var.connectorbtrs_property_client_secret
    "rsAPIurl" = var.rs_api_url
  })
}
