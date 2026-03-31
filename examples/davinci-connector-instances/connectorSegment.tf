resource "pingone_davinci_connector_instance" "connectorSegment" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "connectorSegment"
  }
  name = "My awesome connectorSegment"
  property {
    name  = "body"
    type  = "string"
    value = var.connectorsegment_property_body
  }
  property {
    name  = "endpoint"
    type  = "string"
    value = var.connectorsegment_property_endpoint
  }
  property {
    name  = "genericProperties"
    type  = "string"
    value = var.connectorsegment_property_generic_properties
  }
  property {
    name  = "genericTraits"
    type  = "string"
    value = var.connectorsegment_property_generic_traits
  }
  property {
    name  = "groupId"
    type  = "string"
    value = var.connectorsegment_property_group_id
  }
  property {
    name  = "headers"
    type  = "string"
    value = var.connectorsegment_property_headers
  }
  property {
    name  = "method"
    type  = "string"
    value = var.connectorsegment_property_method
  }
  property {
    name  = "pageName"
    type  = "string"
    value = var.connectorsegment_property_page_name
  }
  property {
    name  = "previousId"
    type  = "string"
    value = var.connectorsegment_property_previous_id
  }
  property {
    name  = "queryParameters"
    type  = "string"
    value = var.connectorsegment_property_query_parameters
  }
  property {
    name  = "screenName"
    type  = "string"
    value = var.connectorsegment_property_screen_name
  }
  property {
    name  = "userEvent"
    type  = "string"
    value = var.connectorsegment_property_user_event
  }
  property {
    name  = "userId"
    type  = "string"
    value = var.connectorsegment_property_user_id
  }
  property {
    name  = "version"
    type  = "string"
    value = var.connectorsegment_property_version
  }
  property {
    name  = "writeKey"
    type  = "string"
    value = var.connectorsegment_property_write_key
  }
}
