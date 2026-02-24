resource "pingone_davinci_connector_instance" "connectorValidsoft" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorValidsoft"
  }
  name = "My awesome connectorValidsoft"
  properties = jsonencode({
    "customAuth" = jsonencode({})
  })
}
