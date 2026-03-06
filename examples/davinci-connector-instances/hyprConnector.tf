resource "pingone_davinci_connector_instance" "hyprConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "hyprConnector"
  }
  name = "My awesome hyprConnector"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
