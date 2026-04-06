resource "pingone_davinci_connector_instance" "connectorAbuseipdb" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorAbuseipdb"
  }
  name = "My awesome connectorAbuseipdb"
  property {
    name  = "apiKey"
    type  = "string"
    value = var.connectorabuseipdb_property_api_key
  }
  property {
    name  = "ipAddress"
    type  = "string"
    value = var.connectorabuseipdb_property_ip_address
  }
  property {
    name  = "maxDays"
    type  = "string"
    value = var.connectorabuseipdb_property_max_days
  }
}
