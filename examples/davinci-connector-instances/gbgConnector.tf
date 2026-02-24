resource "pingone_davinci_connector_instance" "gbgConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "gbgConnector"
  }
  name = "My awesome gbgConnector"
  properties = jsonencode({
    "password" = var.gbgconnector_property_password
    "requestUrl" = var.gbgconnector_property_request_url
    "soapAction" = var.gbgconnector_property_soap_action
    "username" = var.gbgconnector_property_username
  })
}
