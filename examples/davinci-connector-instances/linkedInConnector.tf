resource "pingone_davinci_connector_instance" "linkedInConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "linkedInConnector"
  }
  name = "My awesome linkedInConnector"
  properties = jsonencode({
    "oauth2" = jsonencode({})
  })
}
