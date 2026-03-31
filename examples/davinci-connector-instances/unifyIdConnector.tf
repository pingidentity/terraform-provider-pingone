resource "pingone_davinci_connector_instance" "unifyIdConnector" {
  environment_id = var.pingone_environment_id

  connector = {
    id = "unifyIdConnector"
  }
  name = "My awesome unifyIdConnector"
  property {
    name  = "accountId"
    type  = "string"
    value = var.unifyidconnector_property_account_id
  }
  property {
    name  = "apiKey"
    type  = "string"
    value = var.unifyidconnector_property_api_key
  }
  property {
    name  = "connectorName"
    type  = "string"
    value = var.unifyidconnector_property_connector_name
  }
  property {
    name  = "details1"
    type  = "string"
    value = var.unifyidconnector_property_details1
  }
  property {
    name  = "details2"
    type  = "string"
    value = var.unifyidconnector_property_details2
  }
  property {
    name  = "iconUrl"
    type  = "string"
    value = var.unifyidconnector_property_icon_url
  }
  property {
    name  = "iconUrlPng"
    type  = "string"
    value = var.unifyidconnector_property_icon_url_png
  }
  property {
    name  = "initAuthSession"
    type  = "string"
    value = var.unifyidconnector_property_init_auth_session
  }
  property {
    name  = "isUserAuthorized"
    type  = "string"
    value = var.unifyidconnector_property_is_user_authorized
  }
  property {
    name  = "sdkToken"
    type  = "string"
    value = var.unifyidconnector_property_sdk_token
  }
  property {
    name  = "showCredAddedOn"
    type  = "string"
    value = var.unifyidconnector_property_show_cred_added_on
  }
  property {
    name  = "showCredAddedVia"
    type  = "string"
    value = var.unifyidconnector_property_show_cred_added_via
  }
  property {
    name  = "toolTip"
    type  = "string"
    value = var.unifyidconnector_property_tool_tip
  }
}
