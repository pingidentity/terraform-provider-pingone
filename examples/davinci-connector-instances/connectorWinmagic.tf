resource "pingone_davinci_connector_instance" "connectorWinmagic" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorWinmagic"
  }
  name = "My awesome connectorWinmagic"
  properties = jsonencode({
    "openId" = jsonencode({})
  })
}
