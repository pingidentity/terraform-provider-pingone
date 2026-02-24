resource "pingone_davinci_connector_instance" "googleConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "googleConnector"
  }
  name = "My awesome googleConnector"
  properties = jsonencode({
    "openId" = jsonencode({})
  })
}
