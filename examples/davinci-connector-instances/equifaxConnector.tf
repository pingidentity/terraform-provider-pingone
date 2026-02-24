resource "pingone_davinci_connector_instance" "equifaxConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "equifaxConnector"
  }
  name = "My awesome equifaxConnector"
  properties = jsonencode({
    "baseUrl" = var.equifaxconnector_property_base_url
    "clientId" = var.equifaxconnector_property_client_id
    "clientSecret" = var.equifaxconnector_property_client_secret
    "equifaxSoapApiEnvironment" = var.equifaxconnector_property_equifax_soap_api_environment
    "memberNumber" = var.equifaxconnector_property_member_number
    "password" = var.equifaxconnector_property_password
    "username" = var.equifaxconnector_property_username
  })
}
