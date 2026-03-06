resource "pingone_davinci_connector_instance" "microsoftIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "microsoftIdpConnector"
  }
  name = "My awesome microsoftIdpConnector"
  properties = jsonencode({
    "openId" = jsonencode({})
  })
}
