resource "pingone_davinci_connector_instance" "connector1Kosmos" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connector1Kosmos"
  }
  name = "My awesome connector1Kosmos"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
