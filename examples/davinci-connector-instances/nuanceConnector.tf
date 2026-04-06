resource "pingone_davinci_connector_instance" "nuanceConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "nuanceConnector"
  }
  name = "My awesome nuanceConnector"
  property {
    name  = "authDescription"
    type  = "string"
    value = var.nuanceconnector_property_auth_description
  }
  property {
    name  = "configSetName"
    type  = "string"
    value = var.nuanceconnector_property_config_set_name
  }
  property {
    name  = "connectorName"
    type  = "string"
    value = var.nuanceconnector_property_connector_name
  }
  property {
    name  = "credId"
    type  = "string"
    value = var.nuanceconnector_property_cred_id
  }
  property {
    name  = "customCSS"
    type  = "string"
    value = var.nuanceconnector_property_custom_css
  }
  property {
    name  = "customHTML"
    type  = "string"
    value = var.nuanceconnector_property_custom_html
  }
  property {
    name  = "customScript"
    type  = "string"
    value = var.nuanceconnector_property_custom_script
  }
  property {
    name  = "description"
    type  = "string"
    value = var.nuanceconnector_property_description
  }
  property {
    name  = "details1"
    type  = "string"
    value = var.nuanceconnector_property_details1
  }
  property {
    name  = "details2"
    type  = "string"
    value = var.nuanceconnector_property_details2
  }
  property {
    name  = "htmlConfig"
    type  = "string"
    value = var.nuanceconnector_property_html_config
  }
  property {
    name  = "iconUrl"
    type  = "string"
    value = var.nuanceconnector_property_icon_url
  }
  property {
    name  = "iconUrlPng"
    type  = "string"
    value = var.nuanceconnector_property_icon_url_png
  }
  property {
    name  = "mainHeaderText"
    type  = "string"
    value = var.nuanceconnector_property_main_header_text
  }
  property {
    name  = "nextButtonText"
    type  = "string"
    value = var.nuanceconnector_property_next_button_text
  }
  property {
    name  = "passphrase"
    type  = "string"
    value = var.nuanceconnector_property_passphrase
  }
  property {
    name  = "passphrase1"
    type  = "string"
    value = var.nuanceconnector_property_passphrase1
  }
  property {
    name  = "passphrase2"
    type  = "string"
    value = var.nuanceconnector_property_passphrase2
  }
  property {
    name  = "passphrase3"
    type  = "string"
    value = var.nuanceconnector_property_passphrase3
  }
  property {
    name  = "passphrase4"
    type  = "string"
    value = var.nuanceconnector_property_passphrase4
  }
  property {
    name  = "passphrase5"
    type  = "string"
    value = var.nuanceconnector_property_passphrase5
  }
  property {
    name  = "screen0Config"
    type  = "string"
    value = var.nuanceconnector_property_screen0_config
  }
  property {
    name  = "screen1Config"
    type  = "string"
    value = var.nuanceconnector_property_screen1_config
  }
  property {
    name  = "screen2Config"
    type  = "string"
    value = var.nuanceconnector_property_screen2_config
  }
  property {
    name  = "screen3Config"
    type  = "string"
    value = var.nuanceconnector_property_screen3_config
  }
  property {
    name  = "screen4Config"
    type  = "string"
    value = var.nuanceconnector_property_screen4_config
  }
  property {
    name  = "screen5Config"
    type  = "string"
    value = var.nuanceconnector_property_screen5_config
  }
  property {
    name  = "sessionId"
    type  = "string"
    value = var.nuanceconnector_property_session_id
  }
  property {
    name  = "showCredAddedOn"
    type  = "string"
    value = var.nuanceconnector_property_show_cred_added_on
  }
  property {
    name  = "showCredAddedVia"
    type  = "string"
    value = var.nuanceconnector_property_show_cred_added_via
  }
  property {
    name  = "title"
    type  = "string"
    value = var.nuanceconnector_property_title
  }
  property {
    name  = "toolTip"
    type  = "string"
    value = var.nuanceconnector_property_tool_tip
  }
  property {
    name  = "useCustomScreens"
    type  = "string"
    value = var.nuanceconnector_property_use_custom_screens
  }
  property {
    name  = "voiceUrl"
    type  = "string"
    value = var.nuanceconnector_property_voice_url
  }
  property {
    name  = "voiceprintTag"
    type  = "string"
    value = var.nuanceconnector_property_voiceprint_tag
  }
}
