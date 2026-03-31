resource "pingone_davinci_connector_instance" "connectorBTpra" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBTpra"
  }
  name = "My awesome connectorBTpra"
  property {
    name  = "clientID"
    type  = "string"
    value = var.connectorbtpra_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.connectorbtpra_property_client_secret
  }
  property {
    name  = "hostName"
    type  = "string"
    value = var.connectorbtpra_property_host_name
  }
  property {
    name  = "praAPIurl"
    type  = "string"
    value = var.pra_api_url
  }
  property {
    name  = "userName"
    type  = "string"
    value = var.connectorbtpra_property_user_name
  }
}
