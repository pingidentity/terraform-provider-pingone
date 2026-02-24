resource "pingone_davinci_connector_instance" "cookieConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "cookieConnector"
  }
  name = "My awesome cookieConnector"
  properties = jsonencode({
    "hmacSigningKey" = var.cookieconnector_property_hmac_signing_key
  })
}
