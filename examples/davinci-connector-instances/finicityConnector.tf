resource "pingone_davinci_connector_instance" "finicityConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "finicityConnector"
  }
  name = "My awesome finicityConnector"
  properties = jsonencode({
    "appKey" = var.finicityconnector_property_app_key
    "baseUrl" = var.finicityconnector_property_base_url
    "partnerId" = var.finicityconnector_property_partner_id
    "partnerSecret" = var.finicityconnector_property_partner_secret
  })
}
