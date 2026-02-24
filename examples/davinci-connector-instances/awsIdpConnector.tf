resource "pingone_davinci_connector_instance" "awsIdpConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "awsIdpConnector"
  }
  name = "My awesome awsIdpConnector"
  properties = jsonencode({
    "openId" = jsonencode({})
  })
}
