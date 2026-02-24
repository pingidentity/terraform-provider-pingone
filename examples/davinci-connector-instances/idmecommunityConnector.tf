resource "pingone_davinci_connector_instance" "idmecommunityConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "idmecommunityConnector"
  }
  name = "My awesome idmecommunityConnector"
  properties = jsonencode({
    "openId" = var.idmecommunityconnector_property_open_id
  })
}
