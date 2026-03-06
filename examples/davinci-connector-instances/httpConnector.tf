resource "pingone_davinci_connector_instance" "httpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "httpConnector"
  }
  name = "My awesome httpConnector"
  properties = jsonencode({
    "connectionId" = var.httpconnector_property_connection_id
    "recaptchaSecretKey" = var.httpconnector_property_recaptcha_secret_key
    "recaptchaSiteKey" = var.httpconnector_property_recaptcha_site_key
    "whiteList" = var.httpconnector_property_white_list
  })
}
