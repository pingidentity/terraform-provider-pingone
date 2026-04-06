resource "pingone_davinci_connector_instance" "connectorIPStack" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIPStack"
  }
  name = "My awesome connectorIPStack"
  property {
    name  = "allowInsecureIPStackConnection"
    type  = "string"
    value = var.allow_insecure_ip_stack_connection
  }
  property {
    name  = "apiKey"
    type  = "string"
    value = var.connectoripstack_property_api_key
  }
  property {
    name  = "ip"
    type  = "string"
    value = var.connectoripstack_property_ip
  }
}
