resource "pingone_davinci_connector_instance" "connectorClearbit" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorClearbit"
  }
  name = "My awesome connectorClearbit"
  properties = jsonencode({
    "apiKey" = var.connectorclearbit_property_api_key
    "riskApiVersion" = var.connectorclearbit_property_risk_api_version
    "version" = var.connectorclearbit_property_version
  })
}
