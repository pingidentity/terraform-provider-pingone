resource "pingone_davinci_connector_instance" "connectorGoogleChromeEnterprise" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorGoogleChromeEnterprise"
  }
  name = "My awesome connectorGoogleChromeEnterprise"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
