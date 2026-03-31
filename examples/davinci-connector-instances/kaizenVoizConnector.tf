resource "pingone_davinci_connector_instance" "kaizenVoizConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "kaizenVoizConnector"
  }
  name = "My awesome kaizenVoizConnector"
  property {
    name  = "apiUrl"
    type  = "string"
    value = var.kaizenvoizconnector_property_api_url
  }
  property {
    name  = "applicationName"
    type  = "string"
    value = var.kaizenvoizconnector_property_application_name
  }
  property {
    name  = "authDescription"
    type  = "string"
    value = var.kaizenvoizconnector_property_auth_description
  }
  property {
    name  = "connectorName"
    type  = "string"
    value = var.kaizenvoizconnector_property_connector_name
  }
  property {
    name  = "description"
    type  = "string"
    value = var.kaizenvoizconnector_property_description
  }
  property {
    name  = "details1"
    type  = "string"
    value = var.kaizenvoizconnector_property_details1
  }
  property {
    name  = "details2"
    type  = "string"
    value = var.kaizenvoizconnector_property_details2
  }
  property {
    name  = "iconUrl"
    type  = "string"
    value = var.kaizenvoizconnector_property_icon_url
  }
  property {
    name  = "iconUrlPng"
    type  = "string"
    value = var.kaizenvoizconnector_property_icon_url_png
  }
  property {
    name  = "mainHeaderText"
    type  = "string"
    value = var.kaizenvoizconnector_property_main_header_text
  }
  property {
    name  = "nextButtonText"
    type  = "string"
    value = var.kaizenvoizconnector_property_next_button_text
  }
  property {
    name  = "screen0Config"
    type  = "string"
    value = var.kaizenvoizconnector_property_screen0_config
  }
  property {
    name  = "screen1Config"
    type  = "string"
    value = var.kaizenvoizconnector_property_screen1_config
  }
  property {
    name  = "screen2Config"
    type  = "string"
    value = var.kaizenvoizconnector_property_screen2_config
  }
  property {
    name  = "screen3Config"
    type  = "string"
    value = var.kaizenvoizconnector_property_screen3_config
  }
  property {
    name  = "screen4Config"
    type  = "string"
    value = var.kaizenvoizconnector_property_screen4_config
  }
  property {
    name  = "showCredAddedOn"
    type  = "string"
    value = var.kaizenvoizconnector_property_show_cred_added_on
  }
  property {
    name  = "showCredAddedVia"
    type  = "string"
    value = var.kaizenvoizconnector_property_show_cred_added_via
  }
  property {
    name  = "title"
    type  = "string"
    value = var.kaizenvoizconnector_property_title
  }
  property {
    name  = "toolTip"
    type  = "string"
    value = var.kaizenvoizconnector_property_tool_tip
  }
  property {
    name  = "voiceAccuracy"
    type  = "string"
    value = var.kaizenvoizconnector_property_voice_accuracy
  }
}
