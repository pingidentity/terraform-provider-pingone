resource "pingone_davinci_connector_instance" "splunkConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "splunkConnector"
  }
  name = "My awesome splunkConnector"
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.splunkconnector_property_api_url
  }
  property {
    name  = "event"
    type  = "string"
    value = var.splunkconnector_property_event
  }
  property {
    name  = "eventJSON"
    type  = "string"
    value = var.splunkconnector_property_event_json
  }
  property {
    name  = "host"
    type  = "string"
    value = var.splunkconnector_property_host
  }
  property {
    name  = "index"
    type  = "string"
    value = var.splunkconnector_property_index
  }
  property {
    name  = "metadataKV"
    type  = "string"
    value = var.splunkconnector_property_metadata_kv
  }
  property {
    name  = "port"
    type  = "string"
    value = var.splunkconnector_property_port
  }
  property {
    name  = "source"
    type  = "string"
    value = var.splunkconnector_property_source
  }
  property {
    name  = "sourcetype"
    type  = "string"
    value = var.splunkconnector_property_sourcetype
  }
  property {
    name  = "time"
    type  = "string"
    value = var.splunkconnector_property_time
  }
  property {
    name  = "token"
    type  = "string"
    value = var.splunkconnector_property_token
  }
}
