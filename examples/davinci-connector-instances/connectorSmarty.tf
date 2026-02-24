resource "pingone_davinci_connector_instance" "connectorSmarty" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSmarty"
  }
  name = "My awesome connectorSmarty"
  properties = jsonencode({
    "authId" = var.connectorsmarty_property_auth_id
    "authToken" = var.connectorsmarty_property_auth_token
    "license" = var.connectorsmarty_property_license
  })
}
