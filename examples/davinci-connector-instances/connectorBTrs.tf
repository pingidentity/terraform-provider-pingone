resource "pingone_davinci_connector_instance" "connectorBTrs" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBTrs"
  }
  name = "My awesome connectorBTrs"
  property {
    name  = "clientID"
    type  = "string"
    value = var.connectorbtrs_property_client_id
  }
  property {
    name  = "clientSecret"
    type  = "string"
    value = var.connectorbtrs_property_client_secret
  }
  property {
    name  = "hostName"
    type  = "string"
    value = var.connectorbtrs_property_host_name
  }
  property {
    name  = "rsAPIurl"
    type  = "string"
    value = var.rs_api_url
  }
  property {
    name  = "userName"
    type  = "string"
    value = var.connectorbtrs_property_user_name
  }
}
