resource "pingone_davinci_connector_instance" "connectorHyprAdapt" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorHyprAdapt"
  }
  name = "My awesome connectorHyprAdapt"
  properties = jsonencode({
    "accessToken" = var.connectorhypradapt_property_access_token
  })
}
