resource "pingone_davinci_connector_instance" "connectorMailgun" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorMailgun"
  }
  name = "My awesome connectorMailgun"
  properties = jsonencode({
    "apiKey" = var.connectormailgun_property_api_key
    "apiVersion" = var.connectormailgun_property_api_version
    "mailgunDomain" = var.connectormailgun_property_mailgun_domain
  })
}
