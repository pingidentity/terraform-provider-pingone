resource "pingone_davinci_connector_instance" "connectorDaonidv" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorDaonidv"
  }
  name = "My awesome connectorDaonidv"
  properties = jsonencode({
    "openId" = jsonencode({})
  })
}
