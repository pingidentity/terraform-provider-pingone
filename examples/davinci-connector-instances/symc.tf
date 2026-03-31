resource "pingone_davinci_connector_instance" "symc" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "symc"
  }
  name = "My awesome symc"
  property {
    name  = "authDescription"
    type  = "string"
    value = var.symc_property_auth_description
  }
  property {
    name  = "connectorName"
    type  = "string"
    value = var.symc_property_connector_name
  }
  property {
    name  = "description"
    type  = "string"
    value = var.symc_property_description
  }
  property {
    name  = "details1"
    type  = "string"
    value = var.symc_property_details1
  }
  property {
    name  = "details2"
    type  = "string"
    value = var.symc_property_details2
  }
  property {
    name  = "iconUrl"
    type  = "string"
    value = var.symc_property_icon_url
  }
  property {
    name  = "iconUrlPng"
    type  = "string"
    value = var.symc_property_icon_url_png
  }
  property {
    name  = "pfxBase64"
    type  = "string"
    value = var.symc_property_pfx_base64
  }
  property {
    name  = "pfxPassword"
    type  = "string"
    value = var.symc_property_pfx_password
  }
  property {
    name  = "pushLoginEnabled"
    type  = "string"
    value = var.symc_property_push_login_enabled
  }
  property {
    name  = "screen0Config"
    type  = "string"
    value = var.symc_property_screen0_config
  }
  property {
    name  = "screen1Config"
    type  = "string"
    value = var.symc_property_screen1_config
  }
  property {
    name  = "screen2Config"
    type  = "string"
    value = var.symc_property_screen2_config
  }
  property {
    name  = "showCredAddedOn"
    type  = "string"
    value = var.symc_property_show_cred_added_on
  }
  property {
    name  = "showCredAddedVia"
    type  = "string"
    value = var.symc_property_show_cred_added_via
  }
  property {
    name  = "title"
    type  = "string"
    value = var.symc_property_title
  }
  property {
    name  = "toolTip"
    type  = "string"
    value = var.symc_property_tool_tip
  }
}
