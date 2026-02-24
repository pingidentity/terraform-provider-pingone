resource "pingone_davinci_connector_instance" "sinchConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "sinchConnector"
  }
  name = "My awesome sinchConnector"
  properties = jsonencode({
    "acceptLanguage" = var.sinchconnector_property_accept_language
    "applicationKey" = var.sinchconnector_property_application_key
    "secretKey" = var.sinchconnector_property_secret_key
  })
}
