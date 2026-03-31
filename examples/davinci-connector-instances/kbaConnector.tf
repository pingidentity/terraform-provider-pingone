resource "pingone_davinci_connector_instance" "kbaConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "kbaConnector"
  }
  name = "My awesome kbaConnector"
  property {
    name  = "authDescription"
    type  = "string"
    value = var.kbaconnector_property_auth_description
  }
  property {
    name  = "bodyHeaderText"
    type  = "string"
    value = var.kbaconnector_property_body_header_text
  }
  property {
    name  = "connectorName"
    type  = "string"
    value = var.kbaconnector_property_connector_name
  }
  property {
    name  = "description"
    type  = "string"
    value = var.kbaconnector_property_description
  }
  property {
    name  = "details1"
    type  = "string"
    value = var.kbaconnector_property_details1
  }
  property {
    name  = "details2"
    type  = "string"
    value = var.kbaconnector_property_details2
  }
  property {
    name  = "formFieldsList"
    type  = "string"
    value = var.kbaconnector_property_form_fields_list
  }
  property {
    name  = "iconUrl"
    type  = "string"
    value = var.kbaconnector_property_icon_url
  }
  property {
    name  = "iconUrlPng"
    type  = "string"
    value = var.kbaconnector_property_icon_url_png
  }
  property {
    name  = "nextButtonText"
    type  = "string"
    value = var.kbaconnector_property_next_button_text
  }
  property {
    name  = "nextEvent"
    type  = "string"
    value = var.kbaconnector_property_next_event
  }
  property {
    name  = "parameters"
    type  = "string"
    value = var.kbaconnector_property_parameters
  }
  property {
    name  = "screen0Config"
    type  = "string"
    value = var.kbaconnector_property_screen0_config
  }
  property {
    name  = "screen1Config"
    type  = "string"
    value = var.kbaconnector_property_screen1_config
  }
  property {
    name  = "screen2Config"
    type  = "string"
    value = var.kbaconnector_property_screen2_config
  }
  property {
    name  = "showCredAddedOn"
    type  = "string"
    value = var.kbaconnector_property_show_cred_added_on
  }
  property {
    name  = "showCredAddedVia"
    type  = "string"
    value = var.kbaconnector_property_show_cred_added_via
  }
  property {
    name  = "title"
    type  = "string"
    value = var.kbaconnector_property_title
  }
  property {
    name  = "toolTip"
    type  = "string"
    value = var.kbaconnector_property_tool_tip
  }
}
