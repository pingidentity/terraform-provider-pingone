resource "pingone_davinci_connector_instance" "splunkConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "splunkConnector"
  }
  name = "My awesome splunkConnector"
  properties = jsonencode({
    "apiUrl" = var.splunkconnector_property_api_url
    "port" = var.splunkconnector_property_port
    "token" = var.splunkconnector_property_token
  })
}
