resource "pingone_davinci_connector_instance" "connectorBeyondIdentity" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorBeyondIdentity"
  }
  name = "My awesome connectorBeyondIdentity"
  properties = jsonencode({
    "openId" = jsonencode({})
  })
}
