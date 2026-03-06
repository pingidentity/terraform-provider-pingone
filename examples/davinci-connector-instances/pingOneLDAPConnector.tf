resource "pingone_davinci_connector_instance" "pingOneLDAPConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "pingOneLDAPConnector"
  }
  name = "My awesome pingOneLDAPConnector"
  properties = jsonencode({
    "clientId" = var.pingoneldapconnector_property_client_id
    "clientSecret" = var.pingoneldapconnector_property_client_secret
    "envId" = var.pingoneldapconnector_property_env_id
    "gatewayId" = var.pingoneldapconnector_property_gateway_id
    "region" = var.pingoneldapconnector_property_region
  })
}
