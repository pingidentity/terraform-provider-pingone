resource "pingone_davinci_connector_instance" "connectorIdMeIdentity" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorIdMeIdentity"
  }
  name = "My awesome connectorIdMeIdentity"
  properties = jsonencode({
    "openId" = jsonencode({})
  })
}
