resource "pingone_davinci_connector_instance" "connectorIPStack" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIPStack"
  }
  name = "My awesome connectorIPStack"
  properties = jsonencode({
    "allowInsecureIPStackConnection" = var.allow_insecure_ip_stack_connection
    "apiKey" = var.connectoripstack_property_api_key
  })
}
