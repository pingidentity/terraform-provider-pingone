resource "pingone_davinci_connector_instance" "connectorSegment" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSegment"
  }
  name = "My awesome connectorSegment"
  properties = jsonencode({
    "version" = var.connectorsegment_property_version
    "writeKey" = var.connectorsegment_property_write_key
  })
}
