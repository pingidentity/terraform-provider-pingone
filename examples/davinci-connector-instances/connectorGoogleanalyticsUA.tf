resource "pingone_davinci_connector_instance" "connectorGoogleanalyticsUA" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorGoogleanalyticsUA"
  }
  name = "My awesome connectorGoogleanalyticsUA"
  properties = jsonencode({
    "trackingID" = var.tracking_id
    "version" = var.connectorgoogleanalyticsua_property_version
  })
}
